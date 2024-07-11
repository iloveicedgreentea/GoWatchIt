package plex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

// Stuff to interface directly with Plex
// of course their api is undocumented and worst of all, in xml. I had to piece it together reading various unofficial API implementations or trial and error

type PlexClient struct {
	ServerURL  string
	Port       string
	HTTPClient http.Client
	ImdbClient *http.Client
	MachineID  string
	ClientIP   string
	MediaType  string
}

// return a new instance of a plex client
func NewClient(url, port string) *PlexClient {
	// remove scheme
	url = strings.Replace(url, "http://", "", -1)
	return &PlexClient{
		ServerURL: url,
		Port:      port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// only used for get all movies
func parseAllMediaContainer(payload []byte) (models.AllMediaContainer, error) {
	var data models.AllMediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// unmarshal xml into a struct
func parseMediaContainer(payload []byte) (models.MediaContainer, error) {
	var data models.MediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func parseSessionMediaContainer(payload []byte) (models.SessionMediaContainer, error) {
	var data models.SessionMediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *PlexClient) getRunningSession() (models.SessionMediaContainer, error) {
	// Get session object
	var data models.SessionMediaContainer
	var err error

	// loop until not empty for 30s
	for i := 0; i < 30; i++ {
		res, err := c.makePlexReq("/status/sessions")
		if err != nil {
			return models.SessionMediaContainer{}, err
		}
		// if no response, keep trying
		if len(res) == 0 {
			log.Debugf("Plex session empty, waiting for %v", 30-i)
			time.Sleep(1 * time.Second)
			continue
		}
		data, err = parseSessionMediaContainer(res)
		if err != nil {
			return models.SessionMediaContainer{}, err
		}

		break
	}

	return data, err
}

// TODO: make this an interface
// GetCodecFromSession gets the codec from a running session
func (c *PlexClient) GetCodecFromSession(uuid string) (string, error) {
	sess, err := c.getRunningSession()
	if err != nil {
		return "", err
	}
	// log.Debugf("Session data: %#v", sess.Video)
	// filter by uuid
	// try up to 15 times until session is active. webhook sends before session is ready
	for i := 0; i < 15; i++ {
		for _, video := range sess.Video {
			log.Debugf("Machine identifier: %s", video.Player.MachineIdentifier)
			if video.Player.MachineIdentifier == uuid {
				log.Debug("Found session matching uuid")
				for _, stream := range video.Media.Part.Stream {
					log.Debugf("Stream: %#v", stream)
					if stream.StreamType == "2" {
						return MapPlexToBeqAudioCodec(stream.DisplayTitle, stream.ExtendedDisplayTitle), nil
					}
				}
			}
		}
		log.Debug("Session not found, waiting 2 seconds")
		time.Sleep(time.Second * 2)
	}
	return "", fmt.Errorf("error getting codec. no session found with uuid %s", uuid)
}

// send a request to Plex to get data about something
func (c *PlexClient) GetMediaData(libraryKey string) (models.MediaContainer, error) {
	res, err := c.makePlexReq(libraryKey)
	if err != nil {
		return models.MediaContainer{}, err
	}

	data, err := parseMediaContainer(res)
	if err != nil {
		return models.MediaContainer{}, err
	}
	return data, nil
}

func insensitiveContains(s string, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// check if its DD+ codec
func containsDDP(s string) bool {
	//English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	ddPlusNames := []string{"ddp", "eac3", "e-ac3", "dd+"}
	for _, name := range ddPlusNames {
		if insensitiveContains(strings.ToLower(s), name) {
			return true
		}
	}

	return false
}

// mapPlexToBeqAudioCodec maps a plex codec metadata to a beq catalog codec name
func MapPlexToBeqAudioCodec(codecTitle, codecExtendTitle string) string {
	log.Debugf("Codecs from plex received: %v, %v", codecTitle, codecExtendTitle)

	// Titles are more likely to have atmos so check it first

	// Atmos logic
	atmosFlag := insensitiveContains(codecExtendTitle, "Atmos") || insensitiveContains(codecTitle, "Atmos")

	// check if contains DDP
	ddpFlag := containsDDP(codecTitle) || containsDDP(codecExtendTitle)

	log.Debugf("Atmos: %v - DD+: %v", atmosFlag, ddpFlag)
	// if true and false, then Atmos
	if atmosFlag && !ddpFlag {
		return "Atmos"
	}

	// if true and true, DD+ Atmos
	if atmosFlag && ddpFlag {
		return "DD+ Atmos"
	}

	// if not atmos and DD+, check later for DD+ Atmos, DD+ 7.1/5.1
	if !atmosFlag && ddpFlag {
		if insensitiveContains(codecTitle, "5.1") {
			return "DD+Atmos5.1Maybe"
		}
		if insensitiveContains(codecTitle, "7.1") {
			return "DD+Atmos7.1Maybe"
		}
	}

	// if False and false, then check others
	// TODO: simplify this like with jellyfin
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	case insensitiveContains(codecTitle, "TRUEHD 7.1") && insensitiveContains(codecExtendTitle, "TrueHD 7.1"):
		return "AtmosMaybe"
	case insensitiveContains(codecTitle, "TRUEHD 7.1") && insensitiveContains(codecExtendTitle, "Surround 7.1"):
		return "AtmosMaybe"
	// DTS:X
	case insensitiveContains(codecExtendTitle, "DTS:X") || insensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-X"
	// DTS MA 7.1 containers but not DTS:X codecs
	case insensitiveContains(codecTitle, "DTS-HD MA 7.1") && !insensitiveContains(codecExtendTitle, "DTS:X") && !insensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case insensitiveContains(codecTitle, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case insensitiveContains(codecTitle, "DTS 5.1"):
		return "DTS 5.1"
	// TrueHD 5.1
	case insensitiveContains(codecTitle, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case insensitiveContains(codecTitle, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// DTS HRA
	case insensitiveContains(codecTitle, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case insensitiveContains(codecTitle, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case insensitiveContains(codecTitle, "LPCM 5.1"):
		return "LPCM 5.1"
	case insensitiveContains(codecTitle, "LPCM 7.1"):
		return "LPCM 7.1"
	case insensitiveContains(codecTitle, "LPCM 2.0"):
		return "LPCM 2.0"
	case insensitiveContains(codecTitle, "AAC Stereo"):
		return "AAC 2.0"
	case insensitiveContains(codecTitle, "AC3 5.1") || insensitiveContains(codecTitle, "EAC3 5.1"):
		return "AC3 5.1"
	case insensitiveContains(codecTitle, "EAC3") || insensitiveContains(codecExtendTitle, "EAC3"):
		return "DD+"
	default:
		return "Empty"
	}

}

// get the type of audio codec for BEQ purpose like atmos, dts-x, etc
func (c *PlexClient) GetAudioCodec(data interface{}) (string, error) {
	var plexAudioCodec string
	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
	// loop instead of index because of edge case with two or more video streams
	log.Debugf("Data type: %T", data)
	if mc, ok := data.(models.MediaContainer); ok {
		// try to get Atmos from file because metadata with Truehd is usually misleading
		f := mc.Video.Media.Part.File
		if strings.Contains(strings.ToLower(f), "atmos") {
			log.Debug("Got atmos codec from filename")
			return MapPlexToBeqAudioCodec(f, f), nil
		}
		for _, val := range mc.Video.Media.Part.Stream {
			if val.StreamType == "2" {
				log.Debugf("Found codecs: %s, %s", val.DisplayTitle, val.ExtendedDisplayTitle)
				return MapPlexToBeqAudioCodec(val.DisplayTitle, val.ExtendedDisplayTitle), nil
			}
		}

		if plexAudioCodec == "" {
			log.Error("did not find codec from plex metadata")
			log.Error("Dumping stream data")
			log.Error(mc.Video.Media.Part.Stream)
			return "", errors.New("no codec found")
		}
	} else {
		return "", errors.New("invalid data type")
	}
	return plexAudioCodec, nil
}

// TODO: rename
// GetPlexMovieDb is used because of the Client interface
func (c *PlexClient) GetPlexMovieDb(payload interface{}) string {
	return ""
}

func (c *PlexClient) makePlexReq(path string) ([]byte, error) {
	// Construct the URL with url.URL
	var u *url.URL

	// Add query parameters if needed
	if strings.Contains(path, "playback") {
		// this MUST use the CLIENT IP and 32500 port not server
		// god forbid plex makes any documentation for their APIs they dont want you using
		u = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%s", config.GetString("signal.playerip"), "32500"),
			Path:   path,
		}
		params := url.Values{}
		// only X-Plex-Target-Client-Identifier MUST be sent and it MUST match the client machine id found in clientIP:32500/resources
		params.Add("X-Plex-Target-Client-Identifier", config.GetString("signal.playermachineidentifier"))
		// API docs says these must be sent, but thats not true at all
		// params.Add("commandID", "0")
		// params.Add("type", "video")
		u.RawQuery = params.Encode()

		log.Debugf("using params for playback query: %s", u.RawQuery)
	} else {
		u = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%s", c.ServerURL, c.Port),
			Path:   path,
		}
	}
	// Create the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("X-Plex-Target-Client-Identifier", c.MachineID)
	log.Debugf("Plex: sending request to %s", u.String())
	// Execute the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when calling plex API: %v", err)
	}
	defer res.Body.Close()

	// Read the response
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// only do this for playback
	if strings.Contains(path, "playback") {
		sData := string(data)

		if strings.Contains(sData, "Bad Request") {
			return nil, errors.New("bad request when calling plex API")
		}
		if strings.Contains(sData, "404") {
			return nil, errors.New("machine ID not found in Plex API - triple check your machine ID and client IP, then check it twice more")
		}
	}

	return data, err
}

// DoPlaybackAction generic func to do playback - play, pause, stop
func (c *PlexClient) DoPlaybackAction(action string) error {
	s := fmt.Sprintf("/player/playback/%s", action)
	log.Debugf("Plex: sending %s request with %s", action, s)
	_, err := c.makePlexReq(s)

	return err
}
