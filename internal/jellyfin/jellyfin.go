package jellyfin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

// Stuff to interface directly with Plex
// of course their api is undocumented and worst of all, in xml. I had to piece it together reading various unofficial API implementations

type JellyfinClient struct {
	ServerURL  string
	Port       string
	HTTPClient http.Client
	ImdbClient *http.Client
	MachineID  string
	ClientIP   string
	MediaType  string
}

// return a new instance of a plex client
func NewClient(url, port string, machineID string, clientIP string) *JellyfinClient {
	return &JellyfinClient{
		ServerURL: url,
		Port:      port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
		MachineID: machineID,
		ClientIP:  clientIP,
	}
}
func (c *JellyfinClient) DoPlaybackAction(action string) error {
    // Implement the action logic specific to Jellyfin
    return nil
}
func (c *JellyfinClient) GetPlexMovieDb(payload interface{}) string {
    // Implement the action logic specific to Jellyfin
    return ""
}

// TODO: finish
func (c *JellyfinClient) GetAudioCodec(payload interface{}) (string, error) {
	return "", nil
}
// generic function to make a request
func (c *JellyfinClient) makeRequest(endpoint string, method string) (io.ReadCloser, error) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%v:%v", c.ServerURL, c.Port),
		Path:   endpoint,
	}
	log.Debugf("Making request to %v", u.String())
	// create request with auth
	r := http.Request{
		Method: strings.ToUpper(method),
		URL:    &u,
		Header: http.Header{},
	}
	// add auth
	// url encoded header value
	r.Header.Add("Authorization", fmt.Sprintf("MediaBrowser Token=\"%v\"", config.GetString("jellyfin.apitoken")))
	// support emby also
	r.Header.Add("X-Emby-Token", fmt.Sprintf("%v", config.GetString("jellyfin.apitoken")))
	// make request
	resp, err := c.HTTPClient.Do(&r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error making request to %#v: %v", u, resp.Status)
	}

	return resp.Body, err
}

// TODO: is paused

func (c *JellyfinClient) GetMetadata(userID, itemID string) (metadata models.JellyfinMetadata, err error) {
	// take the itemID and get the codec
	endpoint := fmt.Sprintf("/Users/%s/Items/%s", userID, itemID)
	r, err := c.makeRequest(endpoint, "get")
	if err != nil {
		return metadata, err
	}
	defer r.Close()

	// unmarshal the response
	var payload models.JellyfinMetadata

	b, err := io.ReadAll(r)
	if err != nil {
		return metadata, err
	}
	err = json.Unmarshal(b, &payload)

	if err != nil {
		log.Debugf("GetCodec Response: %#v", string(b))
		return metadata, err
	}

	return payload, nil
}

// get the codec of a media file returns the codec and the display title e.g eac3, Dolby Digital+
func (c *JellyfinClient) GetCodec(payload models.JellyfinMetadata) (codec, displayTitle, codecProfile string, err error) {
	// get the audio stream
	for _, stream := range payload.MediaStreams {
		if stream.Type == "Audio" {
			// TODO: get profile in additoin to codec? diplsay title too
			log.Debugf("Audio stream: %#v", stream)
			return stream.Codec, stream.DisplayTitle, stream.Profile, nil
		}
	}

	return "", "", "", errors.New("no audio stream found")
}

func (c *JellyfinClient) GetEdition(payload models.JellyfinMetadata) (edition string) {
	var path string
	// extract file name from sources
	for _, source := range payload.MediaSources {
		if source.Type == "Default" {
			path = source.Path
		}
	}
	if path == "" {
		log.Errorf("No path found for edition in jellyfin - %#v", payload)
		return ""
	}

	f := strings.ToLower(path)

	// otherwise try to extract from file name
	switch {
	case strings.Contains(f, "extended"):
		return "Extended"
	case strings.Contains(f, "unrated"):
		return "Unrated"
	case strings.Contains(f, "theatrical"):
		return "Theatrical"
	case strings.Contains(f, "ultimate"):
		return "Ultimate"
	case strings.Contains(f, "director"):
		return "Director"
	case strings.Contains(f, "criterion"):
		return "Criterion"
	default:
		return ""
	}
}



func containsDDP(s string) bool {
	//English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	ddPlusNames := []string{"ddp", "eac3", "e-ac3", "dd+"}
	for _, name := range ddPlusNames {
		if common.InsensitiveContains(strings.ToLower(s), name) {
			return true
		}
	}

	return false
}

func MapJFToBeqAudioCodec(codecTitle, codecExtendTitle string) string {
	log.Debugf("Codecs from jellyfin received: %v, %v", codecTitle, codecExtendTitle)

	// Titles are more likely to have atmos so check it first

	// Atmos logic
	atmosFlag := common.InsensitiveContains(codecExtendTitle, "Atmos") || common.InsensitiveContains(codecTitle, "Atmos")

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

	// Assume eac-3 5.1 is dd+ atmos since almost all metadata says so
	if strings.Contains(codecExtendTitle, "5.1") && ddpFlag {
		return "DD+ Atmos"
	}

	// if not atmos and DD+, return DD+
	if !atmosFlag && ddpFlag {
		return "DD+"
	}

	// if False and false, then check others
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	case common.InsensitiveContains(codecTitle, "TRUEHD 7.1") && common.InsensitiveContains(codecExtendTitle, "TrueHD 7.1"):
		return "AtmosMaybe"
	case common.InsensitiveContains(codecTitle, "TRUEHD 7.1") && common.InsensitiveContains(codecExtendTitle, "Surround 7.1"):
		return "AtmosMaybe"
	// DTS:X
	case common.InsensitiveContains(codecExtendTitle, "DTS:X") || common.InsensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-X"
	// DTS MA 7.1 containers but not DTS:X codecs
	case common.InsensitiveContains(codecTitle, "DTS-HD MA 7.1") && !common.InsensitiveContains(codecExtendTitle, "DTS:X") && !common.InsensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case common.InsensitiveContains(codecExtendTitle, "DTS-HD MA 5.1") || common.InsensitiveContains(codecTitle, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	case common.InsensitiveContains(codecTitle, "DTS") && common.InsensitiveContains(codecExtendTitle, "DTS-HD MA") && common.InsensitiveContains(codecExtendTitle, "5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case common.InsensitiveContains(codecTitle, "DTS 5.1"):
		return "DTS 5.1"
	// TrueHD 5.1
	case common.InsensitiveContains(codecTitle, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case common.InsensitiveContains(codecTitle, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// DTS HRA
	case common.InsensitiveContains(codecTitle, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case common.InsensitiveContains(codecTitle, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case common.InsensitiveContains(codecTitle, "LPCM 5.1"):
		return "LPCM 5.1"
	case common.InsensitiveContains(codecTitle, "LPCM 7.1"):
		return "LPCM 7.1"
	case common.InsensitiveContains(codecTitle, "LPCM 2.0"):
		return "LPCM 2.0"
	case common.InsensitiveContains(codecTitle, "AAC Stereo"):
		return "AAC 2.0"
	case common.InsensitiveContains(codecTitle, "AC3") && common.InsensitiveContains(codecExtendTitle, "5.1"):
		return "AC3 5.1"
	// case common.InsensitiveContains(codecTitle, "AC3 5.1") || common.InsensitiveContains(codecTitle, "EAC3 5.1"):
	// 	return "AC3 5.1"
	default:
		return "Empty"
	}

}