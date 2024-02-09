package jellyfin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
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

// GetAudioCodec is a wrapper for common.Client - returns the audio codec of a given payload
func (c *JellyfinClient) GetAudioCodec(payload interface{}) (string, error) {
	codec, title, profile, layout, err := c.GetCodec(payload.(models.JellyfinMetadata))
	if err != nil {
		return "", err
	}
	// parse the response and map this to the beq codec standards
	return MapJFToBeqAudioCodec(codec, title, profile, layout), nil
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

// GetMetadata returns the metadata for a given itemID
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
		log.Errorf("GetMetadata Response failed: %#v", string(b))
		return metadata, err
	}

	return payload, nil
}

// get the codec of a media file returns the codec and the display title e.g eac3, Dolby Digital+ and profile becuase they are different
func (c *JellyfinClient) GetCodec(payload models.JellyfinMetadata) (codec, displayTitle, codecProfile, layout string, err error) {
	// get the audio stream
	for _, stream := range payload.MediaStreams {
		if stream.Type == "Audio" {
			log.Debugf("Audio stream: codec: %s // display: %s // profile: %s // layout: %s", stream.Codec, stream.DisplayTitle, stream.Profile, stream.ChannelLayout)
			return stream.Codec, stream.DisplayTitle, stream.Profile, stream.ChannelLayout, nil
		}
	}

	return codec, displayTitle, codecProfile, layout, errors.New("no audio stream found")
}

// GetEdition extracts the edition of a media file from a metadata payload
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

// GetJfTMDB extracts the tmdb id of a given itemID because its not returned directly in the metadata for some reason
func (c *JellyfinClient) GetJfTMDB(payload models.JellyfinMetadata) (string, error) {
	urls := payload.ExternalUrls
	log.Debugf("External urls: %#v", urls)
	for _, u := range urls {
		if u.Name == "TheMovieDb" {
			s := strings.Replace(u.URL, "https://www.themoviedb.org/", "", -1)
			// extract the numbers
			re, err := regexp.Compile(`\d+$`)
			if err != nil {
				return "", err
			}
			return re.FindString(s), nil
		}
	}

	return "", errors.New("no tmdb id found")
}


// containsDDP looks for typical DD+ audio codec names
func containsDDP(s string) bool {
	//English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	ddPlusNames := []string{"ddp", "eac3", "e-ac3", "dd+", "dolby digital+"}
	for _, name := range ddPlusNames {
		if common.InsensitiveContains(strings.ToLower(s), name) {
			return true
		}
	}

	return false
}

func containsDtsx(codec, displayTitle, profile, layout string) bool {
	// display title must contain dts:x or dts-x or dtsx
	if common.InsensitiveContains(displayTitle, "DTS:X") || common.InsensitiveContains(displayTitle, "DTS-X") || common.InsensitiveContains(displayTitle, "DTSX") {
		return true
	}
	return false
}

func isDtsMA71(codec, displayTitle, profile, layout string) bool {
	// must have dts and 7.1 layout
	if common.InsensitiveContains(layout, "7.1") && common.InsensitiveContains(profile, "DTS-HD MA") {
		// dts-ha ma 7.1 is a container for dts:x so we have to discount it
		if !common.InsensitiveContains(displayTitle, "DTS:X") && !common.InsensitiveContains(displayTitle, "DTS-X") && !common.InsensitiveContains(displayTitle, "DTSX") {
			return true
		}
	}

	return false
}

func MapJFToBeqAudioCodec(codec, displayTitle, profile, layout string) string {
	log.Debugf("Codecs from jellyfin received: title: %v, display: %v, profile: %v, layout: %s", codec, displayTitle, profile, layout)

	// Titles are more likely to have atmos so check it first

	// Atmos logic
	// if display title contains atmos or codec contains atmos, then its very likely atmos
	atmosFlag := common.InsensitiveContains(displayTitle, "Atmos") || common.InsensitiveContains(codec, "Atmos")

	// check if contains DDP
	ddpFlag := containsDDP(codec) || containsDDP(displayTitle)

	log.Debugf("Atmos: %v - DD+: %v", atmosFlag, ddpFlag)

	// Check most common cases

	// if Atmos not ddp, then Atmos
	if atmosFlag && !ddpFlag {
		return "Atmos"
	}

	// if atmos and ddp, DD+ Atmos
	if atmosFlag && ddpFlag {
		return "DD+ Atmos"
	}

	// Assume eac-3 5.1 or 7.1 is dd+ atmos since it usually is e.x the old guard is "English - Dolby Digital+ - 5.1 - Default" except its actually atmos over dd+5.1
	// without AVR check this is just not granular enough
	// this will attempt DD+ Atmos, then DD+ 5.1, and then DD+
	if !atmosFlag && ddpFlag {
		if common.InsensitiveContains(layout, "5.1") {
			return "DD+Atmos5.1Maybe"
		}
		if common.InsensitiveContains(layout, "7.1") {
			return "DD+Atmos7.1Maybe"
		}
	}
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	// most non-atmos 7.1 titles are actually dts-hd 7.1
	// if codec is truehd and display title contains 7.1, then maybe atmos (will be confirmed when trying to search and it will fallback to THD7.1)
	case common.InsensitiveContains(codec, "truehd") && common.InsensitiveContains(displayTitle, "7.1"):
		return "AtmosMaybe"
	// All DTS based codecs
	case common.InsensitiveContains(codec, "DTS"):
		// DTS:X
		if containsDtsx(codec, displayTitle, profile, layout) {
			return "DTS-X"
		}
		// DTS MA 7.1 containers but not DTS:X codecs
		if isDtsMA71(codec, displayTitle, profile, layout) {
			return "DTS-HD MA 7.1"
		}
		// DTS HA MA 5.1
		if common.InsensitiveContains(displayTitle, "DTS-HD MA") && (common.InsensitiveContains(displayTitle, "5.1") || common.InsensitiveContains(layout, "5.1")) {
			return "DTS-HD MA 5.1"
		}
		// DTS 5.1
		if common.InsensitiveContains(layout, "5.1") {
			return "DTS 5.1"
		}
		// DTS HRA
		if common.InsensitiveContains(displayTitle, "DTS-HD HRA") && common.InsensitiveContains(layout, "7.1") {
			return "DTS-HD HR 7.1"
		}
		if common.InsensitiveContains(displayTitle, "DTS-HD HRA") && common.InsensitiveContains(layout, "5.1") {
			return "DTS-HD HR 5.1"
		}

	
	// TrueHD 5.1
	case common.InsensitiveContains(codec, "truehd") && common.InsensitiveContains(layout, "5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case common.InsensitiveContains(codec, "truehd") && common.InsensitiveContains(layout, "6.1"):
		return "TrueHD 6.1"
	// LPCM
	case common.InsensitiveContains(codec, "lpcm") && common.InsensitiveContains(layout, "5.1"):
		return "LPCM 5.1"
	case common.InsensitiveContains(codec, "lpcm") && common.InsensitiveContains(layout, "7.1"):
		return "LPCM 7.1"
	case common.InsensitiveContains(codec, "lpcm"):
		return "LPCM 2.0"
	case common.InsensitiveContains(codec, "aac"):
		return "AAC 2.0"
	case common.InsensitiveContains(codec, "AC3") && (common.InsensitiveContains(layout, "5.1") || common.InsensitiveContains(displayTitle, "5.1")):
		return "AC3 5.1"
	default:
		return "Empty"
	}

	return "Empty"

}
