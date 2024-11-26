package ezbeq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

type BeqClient struct {
	Scheme              string
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
func NewClient() (*BeqClient, error) {
	if !config.IsBeqEnabled() {
		return nil, nil
	}

	port := config.GetEZBeqPort()
	parsedUrl, err := url.ParseRequestURI(fmt.Sprintf("%s://%s", config.GetEZBeqScheme(),config.GetEZBeqUrl()))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	c := &BeqClient{
		ServerURL: parsedUrl.Host,
		Scheme:    parsedUrl.Scheme,
		Port:      port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
	log := logger.GetLogger()
	log.Debug("created new beq client",
		slog.String("server_url", c.ServerURL),
		slog.String("scheme", c.Scheme),
		slog.String("port", c.Port),
	)

	// update client with latest metadata from minidsp
	err = c.GetStatus()
	if err != nil {
		return c, errors.New("error initializing beq client")
	}

	if len(c.DeviceInfo) == 0 {
		return c, errors.New("no ezBEQ hardware devices found. Check your settings and devices")
	}

	return c, nil
}

// NewRequest returns a new request for ezbeq
func (c *BeqClient) NewRequest(ctx context.Context, skipSearch bool, year int, mediaType models.MediaType, edition models.Edition, TMDB string, codec models.CodecName) *models.BeqSearchRequest {
	log := logger.GetLoggerFromContext(ctx)
	var deviceNames []string

	for _, k := range c.DeviceInfo {
		log.Debug("adding device to request",
			slog.String("device", k.Name),
		)
		deviceNames = append(deviceNames, k.Name)
	}
	if len(deviceNames) == 0 {
		log.Error("no devices found in DeviceInfo")
		// TODO: ensure caller knows to check for nil
		return nil
	}

	// TODO: test and assert device names len

	return &models.BeqSearchRequest{
		DryrunMode:      config.IsBeqDryRun(),
		Slots:           config.GetEZBeqSlots(),
		PreferredAuthor: config.GetEZBeqPreferredAuthor(),
		Devices:         deviceNames,
		SkipSearch:      skipSearch,
		Year:            year,
		MediaType:       mediaType,
		TMDB:            TMDB,
		Codec:           codec,
		Edition:         edition,
	}

}

// GetStatus will get metadata from ezbeq and load into client
func (c *BeqClient) GetStatus() error {
	if c == nil {
		return errors.New("beq client is nil")
	}
	// get all devices
	res, err := c.makeReq("/api/2/devices", nil, http.MethodGet)
	if err != nil {
		return err
	}
	payload, err := mapToBeqDevice(res)
	if err != nil {
		return err
	}
	log.Debug("BEQ payload",
		slog.Any("payload", payload),
	)

	log.Debug("Payload length",
		slog.Int("length", len(payload)),
	)
	// add devices to client, it returns as a map not list
	for _, v := range payload {
		log.Debug("BEQ device",
			slog.String("name", v.Name),
		)
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

// MuteCommand sends a mute on/off true = muted, false = not muted
func (c *BeqClient) MuteCommand(status bool) error {
	if c == nil {
		return errors.New("beq client is nil")
	}
	log.Debug("Running mute command")
	for _, v := range c.DeviceInfo {
		endpoint := fmt.Sprintf("/api/1/devices/%s/mute", v.Name)
		log.Debug("Muting device",
			slog.String("endpoint", endpoint),
		)
		var method string
		switch status {
		case true:
			method = http.MethodPut
		case false:
			method = http.MethodDelete
		}
		log.Debug("Running request",
			slog.String("method", method),
		)
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
		log.Info("Device mute status set",
			slog.String("device", v.Name),
			slog.Bool("mute_status", out.Mute),
		)

		if out.Mute != status {
			return fmt.Errorf("mute value %v requested but mute status is now %v", status, out.Mute)
		}

	}

	return nil
}

// MakeCommand sends the command of payload
func (c *BeqClient) MakeCommand(payload []byte) error {
	if c == nil {
		return errors.New("beq client is nil")
	}
	for _, v := range c.DeviceInfo {
		endpoint := fmt.Sprintf("/api/1/devices/%s", v.Name)
		_, err := c.makeReq(endpoint, payload, http.MethodPatch)
		if err != nil {
			return err
		}
	}

	return nil
}

// generic func for beq requests. Payload should be nil
func (c *BeqClient) makeReq(endpoint string, payload []byte, methodType string) ([]byte, error) {
	if c == nil {
		return nil, errors.New("beq client is nil")
	}
	var setHeader bool
	var req *http.Request
	var err error

	// log.Debugf("Using method %s", methodType)
	switch methodType {
	case http.MethodPut:
		setHeader = true
	case http.MethodPatch:
		setHeader = true
	}
	// log.Debugf("Header is set to %v", setHeader)

	u := url.URL{
		Scheme: c.Scheme,
		Host:   fmt.Sprintf("%s:%s", c.ServerURL, c.Port),
		Path:  endpoint,
	}
	// stupid - https://github.com/golang/go/issues/32897 can't pass a typed nil without panic, because its not an untyped nil
	// extra check in case you pass in []byte{}
	if len(payload) == 0 {
		req, err = http.NewRequest(methodType, u.String(), nil)
	} else {
		req, err = http.NewRequest(methodType, u.String(), bytes.NewBuffer(payload))
	}
	if err != nil {
		return []byte{}, err
	}

	if setHeader {
		req.Header.Set("Content-Type", "application/json")
	}
	// log.Debugf("Using url %s", req.URL)
	// log.Debugf("Headers from req %v", req.Header)
	// simple retry
	res, err := c.makeCallWithRetry(req, 20, endpoint)

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
			log.Debug("Error with request - Retrying",
				slog.String("error", err.Error()),
			)
			time.Sleep(time.Second * 2)
			continue
		}
		defer res.Body.Close()

		resp, err = io.ReadAll(res.Body)
		if err != nil {
			log.Debug("Reading body failed - Retrying",
				slog.String("error", err.Error()),
			)
			time.Sleep(time.Second * 2)
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
			log.Debug("Response status",
				slog.String("body", string(resp)),
				slog.Int("status", status),
			)
			log.Debug("Retrying request...")
			err = fmt.Errorf("error in response: %v", res.Status)
			time.Sleep(time.Second * 2)
			continue
		}
	}

	return resp, err
}

// authorCompare returns true if there is an author
func hasAuthor(s string) bool {
	hasAuthor := strings.ToLower(strings.TrimSpace(s))
	return hasAuthor != "none" && hasAuthor != ""
}

// buildAuthorWhitelist returns a string of authors to search for
func buildAuthorWhitelist(preferredAuthors string, endpoint string) string {
	authors := strings.Split(preferredAuthors, ",")
	for _, v := range authors {
		endpoint += fmt.Sprintf("&authors=%s", strings.TrimSpace(v))
	}
	return endpoint
}

// searchCatalog will use ezbeq to search the catalog and then find the right match. tmdb data comes from plex, matched to ezbeq catalog
func (c *BeqClient) searchCatalog(m *models.BeqSearchRequest) (models.BeqCatalog, error) {
	// url encode because of spaces and stuff
	code := urlEncode(string(m.Codec))
	endpoint := fmt.Sprintf("/api/1/search?audiotypes=%s&years=%d&tmdbid=%s", code, m.Year, m.TMDB)

	// this is an author whitelist for each non-empty author append it to search
	if hasAuthor(m.PreferredAuthor) {
		endpoint = buildAuthorWhitelist(m.PreferredAuthor, endpoint)
	}
	log.Debug("Sending ezbeq search request",
		slog.String("endpoint", endpoint),
	)

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
		// if skipping TMDB, set the IDs to match
		if config.IsJellyfinSkipTMDB() {
			if m.Title == "" {
				return models.BeqCatalog{}, errors.New("title is blank, can't skip TMDB")
			}
			log.Debug("Skipping TMDB for search")
			val.MovieDbID = m.TMDB
			if !strings.EqualFold(val.Title, m.Title) {
				log.Debug("Title mismatch",
					slog.String("catalog_title", val.Title),
					slog.String("search_title", m.Title),
				)
				continue
			}
			log.Debug("Title match",
				slog.String("catalog_title", val.Title),
				slog.String("search_title", m.Title),
			)
		}
		log.Debug("Beq results",
			slog.String("title", val.Title),
			slog.Any("codec", val.AudioTypes),
			slog.String("id", val.ID),
		)
		// if we find a match, return it. Much easier to match on tmdb since plex provides it also
		var audioMatch bool
		// rationale here is some BEQ entries have multiple audio types in one entry
		for _, v := range val.AudioTypes {
			if strings.EqualFold(v, string(m.Codec)) {
				audioMatch = true
				break
			}
		}
		// TODO: matching is probably failing here
		if val.MovieDbID == m.TMDB && val.Year == m.Year && audioMatch {
			log.Debug("Potential match found",
				slog.String("title", val.Title),
				slog.Any("codecs", val.AudioTypes),
			)
			// if it matches, check edition
			if checkEdition(val, m.Edition) {
				log.Info("Found a match in catalog",
					slog.String("author", val.Author),
				)
				return val, nil
			} else {
				log.Error("Found a potential match but editions did not match entry. Not loading")
			}
		}
	}

	return models.BeqCatalog{}, errors.New("beq profile was not found in catalog")
}

// map to Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func checkEdition(val models.BeqCatalog, edition models.Edition) bool {
	// skip if matching disabled
	if config.IsBeqSkipEditionMatching() {
		return true
	}
	valLower := strings.ToLower(val.Edition)
	editionLower := strings.ToLower(string(edition))

	// if edition from beq is empty, any match will do
	if len(val.Edition) == 0 {
		return true
	}

	// if the beq edition contains the string like Extended for "Extended Cut", its ok
	if strings.Contains(valLower, editionLower) {
		return true
	}

	// Some BEQ have short hand names
	switch {
	case strings.Contains(valLower, "dc"):
		return edition == models.EditionDirectorsCut
	case strings.Contains(valLower, "se"):
		return edition == models.EditionSpecialEdition
	case strings.Contains(valLower, "tc"):
		return edition == models.EditionTheatrical
	case strings.Contains(valLower, "uc"):
		return edition == models.EditionUltimate
	case strings.Contains(valLower, "cr"):
		return edition == models.EditionCriterion
	case strings.Contains(valLower, "ur"):
		return edition == models.EditionUnrated
	case strings.Contains(valLower, "ex"):
		return edition == models.EditionExtended
	}

	// if BEQ returns an edition but we have none, and loose matching is enabled, let it match
	if config.IsBeqLooseEditionMatching() {
		if edition == models.EditionNone {
			return true
		}
	}

	return false
}

// Edition support doesn't seem important ATM, might revisit later
// LoadBeqProfile will load a profile into slot 1. If skipSearch true, rest of the params will be used (good for quick reload)
func (c *BeqClient) LoadBeqProfile(m *models.BeqSearchRequest) error {
	// TODO: config isEditionEnabled to ignore edition matching
	if !config.IsBeqEnabled() {
		log.Debug("BEQ is disabled, skipping")
		return nil
	}

	log.Debug("BEQ payload",
		slog.Any("payload", m),
	)

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
			// most metadata contains DD+5.1 or something but its actually DD+ Atmos, so try a few options
		} else if m.Codec == "DD+Atmos5.1Maybe" {
			m.Codec = "DD+ Atmos"
			catalog, err = c.searchCatalog(m)
			// else try DD+ 5.1
			if err != nil {
				m.Codec = "DD+ 5.1"
				catalog, err = c.searchCatalog(m)
				if err != nil {
					m.Codec = "DD+"
					catalog, err = c.searchCatalog(m)
					if err != nil {
						return err
					}
				}
			}
		} else if m.Codec == "DD+Atmos7.1Maybe" {
			m.Codec = "DD+ Atmos"
			catalog, err = c.searchCatalog(m)
			// else try DD+ 7.1
			if err != nil {
				m.Codec = "DD+ 7.1"
				catalog, err = c.searchCatalog(m)
				if err != nil {
					m.Codec = "DD+"
					catalog, err = c.searchCatalog(m)
					if err != nil {
						return err
					}
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
	c.CurrentMediaType = string(m.MediaType)

	if m.EntryID == "" {
		return errors.New("could not find catalog entry for ezbeq")
	}

	if m.DryrunMode {
		return fmt.Errorf("BEQ Dry run msg - Would load title %s -- codec %s -- edition: %s, ezbeq entry ID %s - author %s", catalog.Title, m.Codec, catalog.Edition, m.EntryID, catalog.Author)
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
	log.Debug("Sending BEQ payload",
		slog.Any("payload", payload),
	)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// write payload to each device
	for _, v := range m.Devices {
		endpoint := fmt.Sprintf("/api/2/devices/%s", v)
		_, err = c.makeReq(endpoint, jsonPayload, http.MethodPatch)
		if err != nil {
			log.Debug("Error sending payload",
				slog.String("json_payload", string(jsonPayload)),
				slog.String("endpoint", endpoint),
			)
			return err
		}
	}

	return nil
}

// UnloadBeqProfile will unload all profiles from all devices
func (c *BeqClient) UnloadBeqProfile(m *models.BeqSearchRequest) error {
	if !config.IsBeqEnabled() {
		log.Debug("BEQ is disabled, skipping")
		return nil
	}
	if m.DryrunMode {
		return nil
	}
	log.Debug("Unloading ezBEQ profiles")

	for _, v := range m.Devices {
		for _, k := range m.Slots {
			endpoint := fmt.Sprintf("/api/1/devices/%s/filter/%v", v, k)
			log.Debug("Unloading profile",
				slog.String("endpoint", endpoint),
			)
			_, err := c.makeReq(endpoint, nil, http.MethodDelete)
			if err != nil {
				return err
			}
		}
	}

	// TODO: use a wrapper function
	return nil
}
