package ezbeq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
)

var log = logger.GetLogger()

type BeqClient struct {
	ServerURL           string
	Port                string
	CurrentProfile      string
	CurrentMasterVolume float64
	CurrentMediaType    string
	MuteStatus          bool
	MasterVolume        float64
	HTTPClient          http.Client
	DeviceInfo          []models.BeqDevices
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

	// update client with latest metadata from minidsp
	err := c.GetStatus()
	if err != nil {
		return c, errors.New("error initializing beq client")
	}
	return c, nil
}

// GetStatus will get metadata from ezbeq and load into client
func (c *BeqClient) GetStatus() error {
	// get all devices
	res, err := c.makeReq("/api/2/devices", nil, http.MethodGet)
	if err != nil {
		return err
	}
	payload, err := mapToBeqDevice(res)
	if err != nil {
		return err
	}
	log.Debugf("BEQ payload: %#v", payload)

	log.Debugf("Len of payload is: %v", len(payload))
	// add devices to client, it returns as a map not list
	for _, v := range payload {
		log.Debugf("BEQ device: %#v", v.Name)
		c.DeviceInfo = append(c.DeviceInfo, v)
	}

	if len(c.DeviceInfo) == 0 || c.DeviceInfo == nil {
		return errors.New("no devices found")
	}
	log.Debug("c.DeviceInfo is not 0")

	return nil
}

func mapToBeqDevice(jsonData []byte) (beqPayload map[string]models.BeqDevices, err error) {
	err = json.Unmarshal(jsonData, &beqPayload)

	return beqPayload, err
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func publishWrapper(topic string, msg string) error {
	// trigger automation
	return mqtt.Publish([]byte(msg), config.GetString(fmt.Sprintf("mqtt.%s", topic)))
}

// MuteCommand sends a mute on/off true = muted, false = not muted
func (c *BeqClient) MuteCommand(status bool) error {
	log.Debug("Running mute command")
	for _, v := range c.DeviceInfo {
		endpoint := fmt.Sprintf("/api/1/devices/%s/mute", v.Name)
		log.Debugf("ezbeq: Using endpoint %s", endpoint)
		var method string
		switch status {
		case true:
			method = http.MethodPut
		case false:
			method = http.MethodDelete
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
		log.Infof("Mute status set to %v", out.Mute)

		if out.Mute != status {
			return fmt.Errorf("mute value %v requested but mute status is now %v", status, out.Mute)
		}

		
	}

	return publishWrapper("topicMinidspMuteStatus", fmt.Sprintf("%v", status))
}

// MakeCommand sends the command of payload
func (c *BeqClient) MakeCommand(payload []byte) error {
	for _, v := range c.DeviceInfo {
		endpoint := fmt.Sprintf("/api/1/devices/%s", v.Name)
		log.Debugf("ezbeq: Using endpoint %s", endpoint)
		_, err := c.makeReq(endpoint, payload, http.MethodPatch)
		if err != nil {
			return err
		}
	}

	return nil
}

// generic func for beq requests. Payload should be nil
func (c *BeqClient) makeReq(endpoint string, payload []byte, methodType string) ([]byte, error) {
	var setHeader bool
	var req *http.Request
	var err error

	log.Debugf("Using method %s", methodType)
	switch methodType {
	case http.MethodPut:
		setHeader = true
	case http.MethodPatch:
		setHeader = true
	}
	log.Debugf("Header is set to %v", setHeader)

	url := fmt.Sprintf("%s:%s%s", c.ServerURL, c.Port, endpoint)
	// stupid - https://github.com/golang/go/issues/32897 can't pass a typed nil without panic, because its not an untyped nil
	// extra check in case you pass in []byte{}
	if len(payload) == 0 {
		req, err = http.NewRequest(methodType, url, nil)
	} else {
		req, err = http.NewRequest(methodType, url, bytes.NewBuffer(payload))
	}
	if err != nil {
		return []byte{}, err
	}

	if setHeader {
		req.Header.Set("Content-Type", "application/json")
	}
	log.Debugf("Using url %s", req.URL)
	log.Debugf("Headers from req %v", req.Header)
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
			return nil, fmt.Errorf("got status: %d -- error from body is %v", status, string(resp))
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
func (c *BeqClient) searchCatalog(m *models.SearchRequest) (models.BeqCatalog, error) {
	// url encode because of spaces and stuff
	code := urlEncode(m.Codec)
	var endpoint string
	// done this way to make it easier to add future authors
	switch m.PreferredAuthor {
	case "None", "none", "":
		endpoint = fmt.Sprintf("/api/1/search?audiotypes=%s&years=%d", code, m.Year)
	default:
		endpoint = fmt.Sprintf("/api/1/search?audiotypes=%s&years=%d&authors=%s", code, m.Year, urlEncode(m.PreferredAuthor))
	}
	log.Debugf("sending ezbeq search request to %s", endpoint)

	var payload []models.BeqCatalog
	res, err := c.makeReq(endpoint, nil, http.MethodGet)
	if err != nil {
		return models.BeqCatalog{}, err
	}

	err = json.Unmarshal(res, &payload)
	if err != nil {
		return models.BeqCatalog{}, fmt.Errorf("error: %v // response: %v", err, string(res))
	}

	// search through results and find match
	for _, val := range payload {
		log.Debugf("Beq results: Title: %v -- Codec %v, ID: %v", val.Title, val.AudioTypes, val.ID)
		// if we find a match, return it. Much easier to match on tmdb since plex provides it also
		if val.MovieDbID == m.TMDB && val.Year == m.Year && val.AudioTypes[0] == m.Codec {
			// if it matches, check edition
			if checkEdition(val, m.Edition) {
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
func (c *BeqClient) LoadBeqProfile(m *models.SearchRequest) error {
	if m.TMDB == "" {
		return errors.New("tmdb is empty. Can't find a match")
	}
	log.Debugf("beq payload is %#v", m)

	// if no devices provided, error
	if len(m.Devices) == 0 {
		return fmt.Errorf("no ezbeq devices provided. Can't load")
	}

	var err error
	var catalog models.BeqCatalog

	// if provided stuff is blank, we cant skip search
	if m.EntryID == "" || m.MVAdjust == 0 {
		m.SkipSearch = false
	}

	// skip searching when resuming for speed
	if !m.SkipSearch {
		// if AtmosMaybe, check if its really truehd 7.1. If fails, its atmos
		if m.Codec == "AtmosMaybe" {
			m.Codec = "TrueHD 7.1"
			catalog, err = c.searchCatalog(m)
			if err != nil {
				m.Codec = "Atmos"
				catalog, err = c.searchCatalog(m)
				if err != nil {
					return err
				}
			}
		} else {
			catalog, err = c.searchCatalog(m)
			if err != nil {
				return err
			}
		}
		// get the values from catalog search
		m.EntryID = catalog.ID
		m.MVAdjust = catalog.MvAdjust
	} else {
		log.Debug("Skipping search for extra speed")
	}

	// save the current stuff for later, used in media.resume
	c.CurrentMasterVolume = m.MVAdjust
	c.CurrentProfile = m.EntryID
	c.CurrentMediaType = m.MediaType

	if m.EntryID == "" {
		return errors.New("could not find catalog entry for ezbeq")
	}

	if m.DryrunMode {
		return fmt.Errorf("BEQ Dry run msg - Would load title %s -- codec %s -- edition: %s, ezbeq entry ID %s", catalog.SortTitle, m.Codec, catalog.Edition, m.EntryID)
	}

	// build payload
	var payload models.BeqPatchV2
	// for len m.Slots, add that many slots
	// if no slots, add one so it doesnt error
	if len(m.Slots) == 0 {
		m.Slots = []int{1}
	}
	for _, k := range m.Slots {
		// append a slot to payload for each
		payload.Slots = append(payload.Slots, models.SlotsV2{
			ID:     strconv.Itoa(k),
			Gains:  []float64{m.MVAdjust, m.MVAdjust},
			Active: true,
			Mutes:  []bool{false, false},
			Entry:  m.EntryID,
		})
	}
	log.Debugf("sending BEQ payload: %#v", payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// write payload to each device
	for _, v := range m.Devices {
		endpoint := fmt.Sprintf("/api/2/devices/%s", v)
		log.Debugf("json payload %v", string(jsonPayload))
		log.Debugf("using endpoint %s", endpoint)
		_, err = c.makeReq(endpoint, jsonPayload, http.MethodPatch)
		if err != nil {
			log.Debugf("json payload %v", string(jsonPayload))
			log.Debugf("using endpoint %s", endpoint)
			return err
		}
	}

	return publishWrapper("topicBeqCurrentProfile", catalog.SortTitle)
}

// UnloadBeqProfile will unload all profiles from all devices
func (c *BeqClient) UnloadBeqProfile(m *models.SearchRequest) error {
	if m.DryrunMode {
		return nil
	}
	log.Debug("Unloading ezBEQ profiles")

	for _, v := range m.Devices {
		for _, k := range m.Slots {
			endpoint := fmt.Sprintf("/api/1/devices/%s/filter/%v", v, k)
			log.Debugf("using endpoint %s", endpoint)
			_, err := c.makeReq(endpoint, nil, http.MethodDelete)
			if err != nil {
				return err
			}
		}
	}

	return publishWrapper("topicBeqCurrentProfile", "")
}
