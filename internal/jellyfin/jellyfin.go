package jellyfin

import (
	"encoding/json"
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
	log.Debugf("Request: %#v", r)
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


// get the codec of a media file
func (c *JellyfinClient) GetCodec(userID, itemID string) (string, error) {
	// take the itemID and get the codec
	endpoint := fmt.Sprintf("/Users/%s/Items/%s", userID, itemID)
	r, err := c.makeRequest(endpoint, "get")
	if err != nil {
		return "", err
	}
	defer r.Close()

	// unmarshal the response
	var payload models.JellyfinMetadata

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &payload)
	
	if err != nil {
		log.Debugf("GetCodec Response: %v", string(b))
		return "", err
	}
	// TODO: get which stream is audio
	log.Debug(payload.MediaStreams)
	
	return "ok", nil
	
}

func GetEdition() string {
	return "jellyfin"
}

