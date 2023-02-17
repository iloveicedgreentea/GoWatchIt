package ezbeq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

type BeqClient struct {
	ServerURL           string
	Port                string
	DeviceName          string
	CurrentProfile      string
	CurrentMasterVolume float64
	CurrentMediaType    string
	MuteStatus          bool
	MasterVolume        float64
	HTTPClient          http.Client
}

// return a new instance of a plex client
func NewClient(url, port string) (*BeqClient, error) {
	c := &BeqClient{
		ServerURL: url,
		Port:      port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}

	// update client with latest metadata
	err := c.GetStatus()
	if err != nil {
		return c, errors.New("rrror initializing beq client")
	}
	return c, nil
}

// GetStatus will get metadata from eqbeq and load into client
func (c *BeqClient) GetStatus() error {
	var beqPayload models.BeqDevices

	res, err := c.makeReq("/api/1/devices", nil, "get")
	if err != nil {
		return err
	}

	err = json.Unmarshal(res, &beqPayload)
	if err != nil {
		return err
	}
	if beqPayload.Name == "" {
		return errors.New("could not get device name")
	}

	c.DeviceName = beqPayload.Name
	c.MuteStatus = beqPayload.Mute
	c.MasterVolume = beqPayload.MasterVolume

	return nil
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

// MuteCommand sends a mute on/off true = muted, false = not muted
func (c *BeqClient) MuteCommand(status bool) error {
	endpoint := fmt.Sprintf("/api/1/devices/%s/mute", c.DeviceName)
	log.Debugf("ezbeq: Using endpoint %s", endpoint)
	var method string
	switch status {
	case true:
		method = "put"
	case false:
		method = "delete"
	}
	log.Debugf("Running request with method: %s", method)
	resp, err := c.makeReq(endpoint, nil, method)
	if err != nil {
		return err
	}

	// ensure we changed the status
	var out models.BeqDevices
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return err
	}
	log.Debugf("Current mute status is %v", out.Mute)
	if out.Mute != status {
		return fmt.Errorf("mute value %v requested but mute status is now %v", status, out.Mute)
	}

	return nil
}

func (c *BeqClient) MakeCommand(payload []byte) error {
	endpoint := fmt.Sprintf("/api/1/devices/%s", c.DeviceName)
	log.Debugf("ezbeq: Using endpoint %s", endpoint)
	_, err := c.makeReq(endpoint, payload, "patch")

	return err
}

// generic func for beq requests. Payload should be nil
func (c *BeqClient) makeReq(endpoint string, payload []byte, methodType string) ([]byte, error) {
	var method string
	var setHeader bool
	var req *http.Request
	var err error

	switch methodType {
	case "put":
		method = http.MethodPut
		setHeader = true
	case "patch":
		method = http.MethodPatch
		setHeader = true
	case "delete":
		method = http.MethodDelete
	case "get":
		method = http.MethodGet
	}

	url := fmt.Sprintf("%s:%s%s", c.ServerURL, c.Port, endpoint)

	// stupid - https://github.com/golang/go/issues/32897 can't pass a typed nil without panic, because its not an untyped nil
	// extra check in case you pass in []byte{}
	if len(payload) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	}
	if err != nil {
		return []byte{}, err
	}

	if setHeader {
		req.Header.Set("Content-Type", "application/json")
	}

	// simple retry
	res, err := c.makeCallWithRetry(req, 5, endpoint)

	return res, err
}

// makeCallWithRetry returns response body and err
func (c *BeqClient) makeCallWithRetry(req *http.Request, maxRetries int, endpoint string) ([]byte, error) {
	// declaring here so we can return err outside of loop just by exiting it
	var status int
	var res *http.Response
	var resp []byte
	var err error

	for i := 0; i < maxRetries; i++ {
		res, err = c.HTTPClient.Do(req)
		if err != nil {
			log.Debugf("Error with request - Retrying %v", err)
			continue
		}
		defer res.Body.Close()

		resp, err = io.ReadAll(res.Body)
		if err != nil {
			log.Debugf("Reading body failed - Retrying %v", err)
			continue
		}

		status = res.StatusCode

		if status != http.StatusOK {
			return nil, fmt.Errorf("got status: %d", status)
		}

		// don't retry for 404
		if status == 404 {
			return resp, fmt.Errorf("404 for %s", endpoint)
		}

		if status >= 204 && status != 404 {
			log.Debug(string(resp), status)
			log.Debug("Retrying request...")
			err = fmt.Errorf("error in response: %v", res.Status)
			continue
		}
	}

	return resp, err
}

// searchCatalog will use ezbeq to search the catalog and then find the right match. tmdb data comes from plex, matched to ezbeq catalog
func (c *BeqClient) searchCatalog(tmdb string, year int, codec string, preferredAuthor string, edition string) (models.BeqCatalog, error) {
	// url encode because of spaces and stuff
	code := urlEncode(codec)
	var endpoint string
	// done this way to make it easier to add future authors
	switch preferredAuthor {
	case "None", "none", "":
		endpoint = fmt.Sprintf("/api/1/search?audiotypes=%s&years=%d", code, year)
	default:
		endpoint = fmt.Sprintf("/api/1/search?audiotypes=%s&years=%d&authors=%s", code, year, urlEncode(preferredAuthor))
	}
	log.Debugf("sending ezbeq search request to %s", endpoint)

	var payload []models.BeqCatalog
	res, err := c.makeReq(endpoint, nil, "get")
	if err != nil {
		return models.BeqCatalog{}, err
	}

	err = json.Unmarshal(res, &payload)
	if err != nil {
		return models.BeqCatalog{}, fmt.Errorf("error: %v // response: %v", err, string(res))
	}

	// search through results and find match
	for _, val := range payload {
		log.Debugf("Beq results: %v", val)
		// if we find a match, return it. Much easier to match on tmdb since plex provides it also
		if val.MovieDbID == tmdb && val.Year == year && val.AudioTypes[0] == codec {
			// if it matches, check edition
			if checkEdition(val, edition) {
				return val, nil
			} else {
				log.Error("Found a match but editions did not match entry. Not loading")
			}
		}
	}

	return models.BeqCatalog{}, errors.New("beq profile was not found in catalog")
}

// map to Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func checkEdition(val models.BeqCatalog, edition string) bool {
	// if edition from beq is empty, any match will do
	if val.Edition == "" {
		return true
	}

	// if the beq edition contains the string like Extended for "Extended Cut", its ok
	if strings.Contains(val.Edition, edition) {
		return true
	}

	return false
}

// Edition support doesn't seem important ATM, might revisit later
// LoadBeqProfile will load a profile into slot 1. If skipSearch true, rest of the params will be used (good for quick reload)
func (c *BeqClient) LoadBeqProfile(tmdb string, year int, codec string, skipSearch bool, entryID string, mvAdjust float64, dryrunMode bool, preferredAuthor string, edition string, mediaType string) error {
	if tmdb == "" {
		return errors.New("tmdb is empty. Can't find a match")
	}
	
	var err error
	var catalog models.BeqCatalog

	// if provided stuff is blank, we cant skip search
	if entryID == "" || mvAdjust == 0 {
		skipSearch = false
	}

	// skip searching when resuming for speed
	if !skipSearch {
		// if AtmosMaybe, check if its really truehd 7.1. If fails, its atmos
		if codec == "AtmosMaybe" {
			catalog, err = c.searchCatalog(tmdb, year, "TrueHD 7.1", preferredAuthor, edition)
			if err != nil {
				catalog, err = c.searchCatalog(tmdb, year, "Atmos", preferredAuthor, edition)
				if err != nil {
					return err
				}
			}
		} else {
			catalog, err = c.searchCatalog(tmdb, year, codec, preferredAuthor, edition)
			if err != nil {
				return err
			}
		}
		// get the values from catalog search
		entryID = catalog.ID
		mvAdjust = catalog.MvAdjust
	} else {
		log.Debug("Skipping search for extra speed")
	}

	// save the current stuff for later, used in media.resume
	c.CurrentMasterVolume = mvAdjust
	c.CurrentProfile = entryID
	c.CurrentMediaType = mediaType

	if entryID == "" {
		return errors.New("could not find catalog entry for ezbeq")
	}

	if dryrunMode {
		return fmt.Errorf("BEQ Dry run msg - Would load title %s -- codec %s -- edition: %s, ezbeq entry ID %s", catalog.SortTitle, codec, catalog.Edition, entryID)
	}

	// build payload
	var payload models.BeqPatchV2
	payload.Slots = append(payload.Slots, models.SlotsV2{
		ID:     "1",
		Gains:  []float64{mvAdjust, mvAdjust},
		Active: true,
		Mutes:  []bool{false, false},
		Entry:  entryID,
	})
	log.Debugf("sending BEQ payload: %#v", payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/api/2/devices/%s", c.DeviceName)
	log.Debugf("json payload %v", string(jsonPayload))
	log.Debugf("using endpoint %s", endpoint)
	_, err = c.makeReq(endpoint, jsonPayload, "patch")
	if err != nil {
		log.Debugf("json payload %v", string(jsonPayload))
		log.Debugf("using endpoint %s", endpoint)
		return err
	}

	return nil
}

func (c *BeqClient) UnloadBeqProfile(dryrunMode bool) error {
	if dryrunMode {
		return nil
	}
	log.Debug("Unloading ezBEQ profile")

	// add our last entry id and stuff before deleting
	endpoint := fmt.Sprintf("/api/1/devices/%s/filter/1", c.DeviceName)
	_, err := c.makeReq(endpoint, nil, "delete")
	if err != nil {
		return err
	}

	return nil
}
