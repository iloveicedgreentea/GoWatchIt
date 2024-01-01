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

func  (c *JellyfinClient) GetMetadata(userID, itemID string) (metadata models.JellyfinMetadata, err error) {
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
func (c *JellyfinClient) GetCodec(payload models.JellyfinMetadata) (codec, displayTitle string, err error) {
	// get the audio stream
	for _, stream := range payload.MediaStreams {
		if stream.Type == "Audio" {
			log.Debugf("Codec: %#v", stream)
			return stream.Codec, stream.DisplayTitle, nil
		}
	}

	return "", "", errors.New("no audio stream found")
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
