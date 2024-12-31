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

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

type BeqClient struct {
	Scheme              string
	ServerURL           string
	Port                string
	CurrentMasterVolume float64
	CurrentMediaType    string
	MuteStatus          bool
	MasterVolume        float64
	HTTPClient          *retryablehttp.Client
	DeviceInfo          []models.BeqDevices
}

// return a new instance of a plex client
func NewClient() (*BeqClient, error) {
	if !config.IsBeqEnabled() {
		return &BeqClient{}, nil
	}

	port := config.GetEZBeqPort()
	parsedUrl, err := url.ParseRequestURI(fmt.Sprintf("%s://%s", config.GetEZBeqScheme(), config.GetEZBeqUrl()))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	c := &BeqClient{
		ServerURL:  parsedUrl.Host,
		Scheme:     parsedUrl.Scheme,
		Port:       port,
		HTTPClient: retryablehttp.NewClient(),
	}

	// set timeout
	c.HTTPClient.HTTPClient.Timeout = time.Second * 10
	log := logger.GetLogger()
	c.HTTPClient.Logger = log

	log.Debug("created new beq client",
		slog.String("server_url", c.ServerURL),
		slog.String("scheme", c.Scheme),
		slog.String("port", c.Port),
	)

	// update client with latest metadata from minidsp
	err = c.GetStatus()
	if err != nil {
		return c, fmt.Errorf("error initializing beq client while getting status - %w", err)
	}

	if len(c.DeviceInfo) == 0 {
		return c, errors.New("no ezBEQ hardware devices found. Check your settings and devices")
	}

	return c, nil
}

func (c *BeqClient) GetLoadedProfile() (out map[string]string) {
	out = make(map[string]string)
	if c == nil {
		return out
	}

	if config.IsBeqEnabled() {
		// map the current profile to each device for any active slots
		// assuming the slots would have the same profile because why wouldnt they
		for _, k := range c.DeviceInfo {
			for _, v := range k.Slots {
				if v.Active {
					out[k.Name] = v.Last
				}
			}
		}

		if len(out) == 0 {
			for _, k := range c.DeviceInfo {
				out[k.Name] = "No profile loaded"
			}
		}
	}

	return out
}

// NewRequest returns a new request for ezbeq
func (c *BeqClient) NewRequest(ctx context.Context, skipSearch bool, year int, mediaType models.MediaType, edition models.Edition, tmdb string, codec models.CodecName) *models.BeqSearchRequest {
	log := logger.GetLoggerFromContext(ctx)
	deviceNames := make([]string, 0, len(c.DeviceInfo))

	for _, k := range c.DeviceInfo {
		log.Debug("adding device to request",
			slog.String("device", k.Name),
		)
		deviceNames = append(deviceNames, k.Name)
	}
	if len(deviceNames) == 0 {
		log.Error("no devices found in DeviceInfo")
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
		TMDB:            tmdb,
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
	log := logger.GetLogger()
	log.Debug("Making request",
		slog.String("endpoint", endpoint),
		slog.String("method", methodType),
	)
	if c == nil {
		return nil, errors.New("beq client is nil")
	}
	var setHeader bool
	var req *retryablehttp.Request
	var err error

	switch methodType {
	case http.MethodPut:
		setHeader = true
	case http.MethodPatch:
		setHeader = true
	}
	// caller encodes stuff by using url.Values
	fullURL := fmt.Sprintf("%s://%s:%s%s", c.Scheme, c.ServerURL, c.Port, endpoint)

	// stupid - https://github.com/golang/go/issues/32897 can't pass a typed nil without panic, because its not an untyped nil
	// extra check in case you pass in []byte{}
	if len(payload) == 0 {
		req, err = retryablehttp.NewRequest(methodType, fullURL, nil)
	} else {
		req, err = retryablehttp.NewRequest(methodType, fullURL, bytes.NewBuffer(payload))
	}
	if err != nil {
		return []byte{}, err
	}

	if setHeader {
		req.Header.Set("Content-Type", "application/json")
	}
	log.Debug("Created request object",
		slog.String("url", req.URL.String()),
		slog.String("method", methodType),
	)

	// retry
	res, err := c.makeCallWithRetry(req)

	return res, err
}

// makeCallWithRetry returns response body and err
func (c *BeqClient) makeCallWithRetry(req *retryablehttp.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.GetLogger().Warn("error closing response body: %v")
		}
	}()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Debug("Reading body failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	log.Debug("Response from BEQ",
		slog.String("response", string(resp)),
		slog.String("endpoint", req.URL.String()),
	)

	return resp, err
}

// authorCompare returns true if there is an author
func hasAuthor(s string) bool {
	hasAuthor := strings.ToLower(strings.TrimSpace(s))
	return hasAuthor != "none" && hasAuthor != ""
}

// buildAuthorWhitelist returns a string of authors to search for
func buildAuthorWhitelist(preferredAuthors string, q url.Values) url.Values {
	for _, author := range strings.Split(preferredAuthors, ",") {
		q.Add("authors", strings.TrimSpace(author))
	}

	return q
}

// searchCatalog will use ezbeq to search the catalog and then find the right match. tmdb data comes from plex, matched to ezbeq catalog
func (c *BeqClient) searchCatalog(m *models.BeqSearchRequest) (models.BeqCatalog, error) {
	// build query
	log := logger.GetLogger()
	log.Debug("Searching ezbeq catalog",
		slog.String("title", m.Title),
		slog.String("codec", string(m.Codec)),
		slog.Int("year", m.Year),
		slog.String("tmdb", m.TMDB),
	)
	q := url.Values{}
	q.Add("audiotypes", string(m.Codec))
	q.Add("years", strconv.Itoa(m.Year))
	log.Debug("converted year to string",
		slog.String("year", strconv.Itoa(m.Year)),
	)
	q.Add("tmdbid", m.TMDB)

	// Add authors if present
	if hasAuthor(m.PreferredAuthor) {
		q = buildAuthorWhitelist(m.PreferredAuthor, q)
	}

	endpoint := fmt.Sprintf("/api/1/search?%s", q.Encode())

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
	for i := range payload {
		val := payload[i]
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
			if checkEdition(&val, m.Edition) {
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
func checkEdition(val *models.BeqCatalog, edition models.Edition) bool {
	// skip if matching disabled
	if config.IsBeqSkipEditionMatching() {
		return true
	}
	valLower := strings.ToLower(val.Edition)
	editionLower := strings.ToLower(string(edition))

	// if edition from beq is empty, any match will do
	if val.Edition == "" {
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
	// TODO: cache these in DB for faster lookup. Purge cache on new BEQ entry
	// if provided stuff is blank, we cant skip search
	if m.EntryID == "" || m.MVAdjust == 0 {
		m.SkipSearch = false
	}

	// skip searching when resuming for speed
	if !m.SkipSearch {
		// if AtmosMaybe, check if its really truehd 7.1. If fails, its atmos
		switch m.Codec {
		case "AtmosMaybe":
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
		case "DD+Atmos5.1Maybe":
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
		case "DD+Atmos7.1Maybe":
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
		default:
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
	log.Debug("Sending BEQ payload to devices",
		slog.Int("device_count", len(m.Devices)),
	)
	for _, v := range m.Devices {
		log.Debug("Sending payload to device",
			slog.String("device_name", v),
		)
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

	return nil
}
