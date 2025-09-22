package beq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/mediaplayer"
	"go.uber.org/zap"
)

// return a new instance of a plex client
func NewClient(ctx context.Context) (*BeqClient, error) {
	if !config.IsBeqEnabled() {
		return &BeqClient{}, errors.New("ezBEQ is not enabled")
	}

	port := config.GetEZBeqPort()
	// safely parse the url
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

	// update client with latest metadata from minidsp
	err = c.GetStatus(ctx)
	if err != nil {
		return c, fmt.Errorf("error initializing beq client while getting status - %w", err)
	}

	if len(c.DeviceInfo) == 0 {
		return c, errors.New("no ezBEQ hardware devices found. Check your settings and devices")
	}

	return c, nil
}

// GetStatus will get metadata from ezbeq and load into client
func (c *BeqClient) GetStatus(ctx context.Context) error {
	if c == nil {
		return errors.New("beq client is nil")
	}
	log := logger.GetLoggerFromContext(ctx)
	// get all devices using V2 API
	endpoint := fmt.Sprintf("%s/devices", apiV2Prefix)
	res, err := c.makeReq(ctx, endpoint, nil, http.MethodGet)
	if err != nil {
		return err
	}
	payload, err := mapToBeqDevice(res)
	if err != nil {
		return err
	}

	// add devices to client, it returns as a map not list
	for _, v := range payload {
		log.Debug("Found ezBEQ device", zap.String("name", v.Name))
		c.DeviceInfo = append(c.DeviceInfo, v)
	}

	if len(c.DeviceInfo) == 0 || c.DeviceInfo == nil {
		return errors.New("no devices found")
	}

	return nil
}

// GetCurrentProfile gets profiles loaded per device
func (c *BeqClient) GetCurrentProfile(ctx context.Context) (map[string]string, error) {
	if c == nil {
		return nil, errors.New("beq client is nil")
	}
	out := make(map[string]string)
	for _, v := range c.DeviceInfo {
		for _, slot := range v.Slots {
			if slot.Active {
				out[v.Name] = slot.Last
			}
		}
	}

	return out, nil
}

func mapToBeqDevice(jsonData []byte) (beqPayload map[string]BeqDevices, err error) {
	err = json.Unmarshal(jsonData, &beqPayload)

	return beqPayload, err
}

// generic func for beq requests. Payload should be nil
func (c *BeqClient) makeReq(ctx context.Context, endpoint string, payload []byte, methodType string) ([]byte, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Making request",
		zap.String("endpoint", endpoint),
		zap.String("method", methodType),
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

	// retry
	return c.makeCallWithRetry(ctx, req)
}

// makeCallWithRetry returns response body and err
func (c *BeqClient) makeCallWithRetry(ctx context.Context, req *retryablehttp.Request) ([]byte, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if c == nil {
		return nil, errors.New("beq client is nil")
	}

	log := logger.GetLoggerFromContext(ctx)
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
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	log.Debug("Response from BEQ",
		zap.String("response", string(resp)),
		zap.String("endpoint", req.URL.String()),
	)

	return resp, err
}

// MuteCommand sends a mute on/off true = muted, false = not muted
func (c *BeqClient) MuteCommand(ctx context.Context, status bool) error {
	if c == nil {
		return errors.New("beq client is nil")
	}

	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Running mute command")
	for _, v := range c.DeviceInfo {
		// Note: Still using V1 API endpoint for mute as V2 equivalent might not exist or behave differently.
		endpoint := fmt.Sprintf("%s/devices/%s/mute", apiV1Prefix, v.Name)
		log.Debug("Muting device",
			zap.String("endpoint", endpoint),
		)
		var method string
		switch status {
		case true:
			method = http.MethodPut
		case false:
			method = http.MethodDelete
		}
		log.Debug("Running request",
			zap.String("method", method),
		)
		resp, err := c.makeReq(ctx, endpoint, nil, method)
		if err != nil {
			return err
		}

		// ensure we changed the status
		var out BeqDevices
		err = json.Unmarshal(resp, &out)
		if err != nil {
			return err
		}
		log.Info("Device mute status set",
			zap.String("device", v.Name),
			zap.Bool("mute_status", out.Mute),
		)

		if out.Mute != status {
			return fmt.Errorf("mute value %v requested but mute status is now %v", status, out.Mute)
		}

	}

	return nil
}

// MakeCommand sends the command of payload
func (c *BeqClient) MakeCommand(ctx context.Context, payload []byte) error {
	if c == nil {
		return errors.New("beq client is nil")
	}
	for _, v := range c.DeviceInfo {
		// Note: Still using V1 API endpoint for generic patch as V2 equivalent might not exist or behave differently.
		endpoint := fmt.Sprintf("%s/devices/%s", apiV1Prefix, v.Name)
		_, err := c.makeReq(ctx, endpoint, payload, http.MethodPatch)
		if err != nil {
			return err
		}
	}

	return nil
}

// authorCompare returns true if there is an author
func hasAuthor(s string) bool {
	hasAuthor := strings.ToLower(strings.TrimSpace(s))
	// Use constant for "none" author check
	return hasAuthor != authorNone && hasAuthor != ""
}

// buildAuthorWhitelist returns a string of authors to search for
func buildAuthorWhitelist(preferredAuthors string, q url.Values) url.Values {
	for _, author := range strings.Split(preferredAuthors, ",") {
		q.Add("authors", strings.TrimSpace(author))
	}

	return q
}

// searchCatalog will use ezbeq to search the catalog and then find the right match. tmdb data comes from plex, matched to ezbeq catalog
func (c *BeqClient) searchCatalog(ctx context.Context, m *BEQPayload) (BeqCatalog, error) {
	// build query
	log := logger.GetLogger()
	log.Debug("Searching ezbeq catalog",
		zap.String("title", m.Title),
		zap.String("codec", string(m.Codec)),
		zap.Int("year", m.Year),
		zap.String("tmdb", m.TMDB),
	)
	q := url.Values{}
	q.Add("audiotypes", string(m.Codec))
	// dont add blank year
	if m.Year != 0 {
		q.Add("years", strconv.Itoa(m.Year))
	}
	// add tmdb id
	q.Add("tmdbid", m.TMDB)

	// Add authors if present
	if hasAuthor(m.PreferredAuthor) {
		q = buildAuthorWhitelist(m.PreferredAuthor, q)
	}

	// Note: Still using V1 API endpoint for search as V2 equivalent might not exist or behave differently.
	endpoint := fmt.Sprintf("%s/search?%s", apiV1Prefix, q.Encode())

	var payload []BeqCatalog
	res, err := c.makeReq(ctx, endpoint, nil, http.MethodGet)
	if err != nil {
		return BeqCatalog{}, err
	}

	err = json.Unmarshal(res, &payload)
	if err != nil {
		return BeqCatalog{}, fmt.Errorf("error: %v // response: %v", err, string(res))
	}

	// search through results and find match
	for i := range payload {
		val := payload[i]
		// if skipping TMDB, set the IDs to match
		if config.IsJellyfinSkipTMDB() {
			if m.Title == "" {
				return BeqCatalog{}, errors.New("title is blank, can't skip TMDB")
			}
			log.Debug("Skipping TMDB for search")
			val.MovieDbID = m.TMDB
			if !strings.EqualFold(val.Title, m.Title) {
				log.Debug("Title mismatch",
					zap.String("catalog_title", val.Title),
					zap.String("search_title", m.Title),
				)
				continue
			}
			log.Debug("Title match",
				zap.String("catalog_title", val.Title),
				zap.String("search_title", m.Title),
			)
		}
		log.Debug("Beq results",
			zap.String("title", val.Title),
			zap.Any("codec", val.AudioTypes),
			zap.String("id", val.ID),
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

		// ensure tmdb and codec match
		if val.MovieDbID == m.TMDB && audioMatch {
			log.Debug("Potential match found",
				zap.String("title", val.Title),
				zap.Any("codecs", val.AudioTypes),
				zap.String("found_edition", val.Edition),
				zap.String("requested_edition", string(m.Edition)),
			)
			// if it matches, check edition
			if checkEdition(&val, m.Edition) {
				log.Info("Found a match in catalog",
					zap.String("author", val.Author),
				)
				return val, nil
			} else {
				log.Warn("Found a potential match but editions did not match. Not loading",
					zap.String("title", val.Title),
					zap.Any("codecs", val.AudioTypes),
					zap.String("found_edition", val.Edition),
					zap.String("requested_edition", string(m.Edition)),
				)
			}
		}
	}

	return BeqCatalog{}, errors.New("beq profile was not found in catalog")
}

// map to Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func checkEdition(val *BeqCatalog, edition mediaplayer.Edition) bool {
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
		return edition == mediaplayer.EditionDirectorsCut
	case strings.Contains(valLower, "se"):
		return edition == mediaplayer.EditionSpecialEdition
	case strings.Contains(valLower, "tc"):
		return edition == mediaplayer.EditionTheatrical
	case strings.Contains(valLower, "uc"):
		return edition == mediaplayer.EditionUltimate
	case strings.Contains(valLower, "cr"):
		return edition == mediaplayer.EditionCriterion
	case strings.Contains(valLower, "ur"):
		return edition == mediaplayer.EditionUnrated
	case strings.Contains(valLower, "ex"):
		return edition == mediaplayer.EditionExtended
	}

	// if BEQ returns an edition but we have none, and loose matching is enabled, let it match
	if config.IsBeqLooseEditionMatching() {
		if edition == mediaplayer.EditionNone {
			return true
		}
	}

	return false
}

// LoadBeqProfile will load a profile into slot 1. If skipSearch true, rest of the params will be used (good for quick reload)
func (c *BeqClient) LoadBeqProfile(ctx context.Context, m *BEQPayload) error {
	log := logger.GetLoggerFromContext(ctx)
	if !config.IsBeqEnabled() {
		log.Debug("BEQ is disabled, skipping")
		return nil
	}

	log.Debug("BEQ search request",
		zap.Any("request", m),
	)

	// if no devices provided, error
	if len(m.Devices) == 0 {
		return fmt.Errorf("no ezbeq devices provided. Can't load")
	}

	var err error
	var catalog BeqCatalog
	// TODO: cache these in DB for faster lookup. Purge cache on new BEQ entry
	// if provided stuff is blank, we cant skip search
	if m.EntryID == "" || m.MVAdjust == 0 {
		m.SkipSearch = false
	}

	// skip searching when resuming for speed
	if !m.SkipSearch {
		// Attempting codec variations due to potential metadata ambiguity from source.
		// Use constants for codec names.
		switch m.Codec {
		case mediaplayer.CodecAtmosMaybe:
			m.Codec = mediaplayer.CodecTrueHD71
			catalog, err = c.searchCatalog(ctx, m)
			if err != nil {
				m.Codec = mediaplayer.CodecAtmos
				catalog, err = c.searchCatalog(ctx, m)
				if err != nil {
					return fmt.Errorf("failed to find BEQ %w", err)
				}
			}
		case mediaplayer.CodecDDPlusAtmos51Maybe:
			m.Codec = mediaplayer.CodecDDPAtmos
			catalog, err = c.searchCatalog(ctx, m)
			if err != nil {
				m.Codec = mediaplayer.CodecDDPlus51
				catalog, err = c.searchCatalog(ctx, m)
				if err != nil {
					m.Codec = mediaplayer.CodecDDPlus
					catalog, err = c.searchCatalog(ctx, m)
					if err != nil {
						return fmt.Errorf("failed to find BEQ %w", err)
					}
				}
			}
		case mediaplayer.CodecDDPlusAtmos71Maybe:
			m.Codec = mediaplayer.CodecDDPAtmos
			catalog, err = c.searchCatalog(ctx, m)
			if err != nil {
				m.Codec = mediaplayer.CodecDDPlus71
				catalog, err = c.searchCatalog(ctx, m)
				if err != nil {
					m.Codec = mediaplayer.CodecDDPlus
					catalog, err = c.searchCatalog(ctx, m)
					if err != nil {
						return fmt.Errorf("failed to find BEQ  %w", err)
					}
				}
			}
		default:
			catalog, err = c.searchCatalog(ctx, m)
			if err != nil {
				return fmt.Errorf("failed to find BEQ for %s: %w", m.Codec, err)
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

	// build payload to target ALL configured slots
	var payload BeqPatchV2
	if len(m.Slots) == 0 {
		log.Warn("No BEQ slots configured (m.Slots is empty), attempting to load to slot 1 by default.")
		// If no slots are configured by the user, ezbeq typically defaults to slot 1 for some operations.
		// However, the V2 PATCH expects explicit slot IDs. We'll default to slot 1 here.
		// This situation should ideally be handled by ensuring EZBEQ_SLOTS is always populated if BEQ is enabled.
		payload.Slots = append(payload.Slots, SlotsV2{
			ID:     "1", // Default to slot "1"
			Gains:  []float64{m.MVAdjust, m.MVAdjust},
			Active: true,
			Mutes:  []bool{false, false},
			Entry:  m.EntryID,
		})
	} else {
		log.Info("Targeting all configured BEQ slots for load operation", zap.Any("slots", m.Slots))
		for _, slotNum := range m.Slots {
			payload.Slots = append(payload.Slots, SlotsV2{
				ID:     strconv.Itoa(slotNum),
				Gains:  []float64{m.MVAdjust, m.MVAdjust},
				Active: true,
				Mutes:  []bool{false, false},
				Entry:  m.EntryID,
			})
		}
	}
	log.Debug("Sending BEQ payload",
		zap.Any("payload", payload),
	)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// write payload to each device using V2 API
	for _, deviceName := range m.Devices {
		endpoint := fmt.Sprintf("%s/devices/%s", apiV2Prefix, deviceName)
		_, err = c.makeReq(ctx, endpoint, jsonPayload, http.MethodPatch)
		if err != nil {
			log.Debug("Error sending payload",
				zap.String("json_payload", string(jsonPayload)),
				zap.String("endpoint", endpoint),
			)
			return err
		}
	}

	return nil
}

// UnloadBeqProfile will unload the profile from the first configured slot on all devices using the V2 API.
func (c *BeqClient) UnloadBeqProfile(ctx context.Context, m *BEQPayload) error {
	log := logger.GetLoggerFromContext(ctx)
	if !config.IsBeqEnabled() {
		log.Debug("BEQ is disabled, skipping unload")
		return nil
	}
	if m.DryrunMode {
		log.Info("BEQ Dry run: Would unload profile from first configured slot")
		return nil
	}
	log.Debug("Unloading ezBEQ profile using V2 API")

	// Build V2 payload to deactivate ALL configured slots
	var payload BeqPatchV2
	if len(m.Slots) == 0 {
		log.Warn("No BEQ slots configured (m.Slots is empty), attempting to unload from slot 1 by default.")
		// Defaulting to slot "1" for unload if no slots are configured.
		payload.Slots = append(payload.Slots, SlotsV2{
			ID:     "1",
			Active: false,
			Entry:  "",
			Gains:  []float64{0, 0},
			Mutes:  []bool{false, false},
		})
	} else {
		log.Info("Targeting all configured BEQ slots for unload operation", zap.Any("slots", m.Slots))
		for _, slotNum := range m.Slots {
			payload.Slots = append(payload.Slots, SlotsV2{
				ID:     strconv.Itoa(slotNum),
				Active: false,
				Entry:  "",
				Gains:  []float64{0, 0},
				Mutes:  []bool{false, false},
			})
		}
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal unload payload: %w", err)
	}

	log.Debug("Sending BEQ unload payload", zap.Any("payload", payload))

	var allErrors error

	// Send PATCH request to each device using V2 API, collect errors
	for _, deviceName := range m.Devices {
		endpoint := fmt.Sprintf("%s/devices/%s", apiV2Prefix, deviceName)
		_, err = c.makeReq(ctx, endpoint, jsonPayload, http.MethodPatch)
		if err != nil {
			deviceErr := fmt.Errorf("failed to unload profile from device %s: %w", deviceName, err)
			log.Error("Error sending unload payload to device",
				zap.String("device", deviceName),
				zap.String("json_payload", string(jsonPayload)),
				zap.String("endpoint", endpoint),
				zap.String("error", deviceErr.Error()), // Log the wrapped error
			)
			allErrors = errors.Join(allErrors, deviceErr) // Collect errors
		} else {
			log.Info("Successfully sent unload command to device for configured slots", zap.String("device", deviceName), zap.Any("slots_targeted_in_payload", m.Slots))
		}
	}

	return allErrors // Return combined errors, nil if none occurred
}
