// client.go implements the API for plex itself
package plex

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/iloveicedgreentea/gowatchit/pkg/codecs"
	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/editions"
	"github.com/iloveicedgreentea/gowatchit/pkg/events"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
)

// Stuff to interface directly with Plex
// of course their api is undocumented and worst of all, in xml. I had to piece it together reading various unofficial API implementations or trial and error

const (
	APIStatusSession string = "/status/sessions"
)

type PlexClient struct {
	URL        string
	Port       string
	Scheme     string
	HTTPClient *retryablehttp.Client
	MachineID  string
	ClientIP   string
	MediaType  string
}

// return a new instance of a plex client
func NewClient(ctx context.Context) (*PlexClient, error) {
	c := &PlexClient{
		URL:        config.GetPlayerURL(ctx),
		Port:       config.GetPlayerPort(ctx),
		Scheme:     config.GetPlayerScheme(ctx),
		HTTPClient: retryablehttp.NewClient(),
	}
	// set timeout
	c.HTTPClient.HTTPClient.Timeout = time.Second * 10

	return c, nil
}

// unmarshal xml into a struct
func parseMediaContainer(payload []byte) (MediaContainer, error) {
	var data MediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, parseXMLError(err, string(payload))
	}

	return data, nil
}

func parseSessionMediaContainer(payload []byte) (SessionMediaContainer, error) {
	var data SessionMediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, fmt.Errorf("error unmarshalling parseSessionMediaContainer xml: %v", err)
	}

	return data, nil
}

func (c *PlexClient) getRunningSession(ctx context.Context) (SessionMediaContainer, error) {
	// Get session object
	var data SessionMediaContainer
	var err error

	res, err := c.makePlexReq(ctx, string(APIStatusSession))
	if err != nil {
		return SessionMediaContainer{}, fmt.Errorf("error getting session data: %v", err)
	}
	data, err = parseSessionMediaContainer(res)
	if err != nil {
		return SessionMediaContainer{}, fmt.Errorf("error parsing getRunningSession session data: %v", err)
	}

	return data, err
}

// GetCodecFromSession gets the codec from a running session
func (c *PlexClient) GetCodecFromSession(ctx context.Context, uuid string) (codecs.Codec, error) {
	log := logger.GetLoggerFromContext(ctx)
	sess, err := c.getRunningSession(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting GetCodecFromSession session data: %v", err)
	}
	// log.Debugf("Session data: %#v", sess.Video)
	// filter by uuid
	// try up to 15 times until session is active. webhook sends before session is ready
	for i := 0; i < 15; i++ {
		for i := range sess.Video {
			video := &sess.Video[i]
			log.Debug("Machine identifier",
				zap.String("identifier", video.Player.MachineIdentifier),
			)
			if video.Player.MachineIdentifier == uuid {
				log.Debug("Found session matching uuid",
					zap.String("uuid", uuid),
				)
				for i := range video.Media.Part.Stream {
					stream := &video.Media.Part.Stream[i]
					log.Debug("Stream data",
						zap.String("data", fmt.Sprintf("%#v", stream)),
					)
					if stream.StreamType == "2" {
						return MapPlexToBeqAudioCodec(ctx, stream.DisplayTitle, stream.ExtendedDisplayTitle), nil
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
func (c *PlexClient) getMediaData(ctx context.Context, payload *events.Event) (MediaContainer, error) {
	libraryKey := payload.Metadata.Key
	res, err := c.makePlexReq(ctx, libraryKey)
	if err != nil {
		return MediaContainer{}, err
	}

	return parseMediaContainer(res)
}

func insensitiveContains(s string, sub codecs.Codec) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(string(sub)))
}

// check if its DD+ codec
func containsDDP(s string) bool {
	// English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	ddPlusNames := []codecs.Codec{codecs.CodecDDP, codecs.CodecDDPlus, codecs.CodecEAC3, codecs.CodecEAC3Alt}
	for _, name := range ddPlusNames {
		if insensitiveContains(strings.ToLower(s), name) {
			return true
		}
	}

	return false
}

// mapPlexToBeqAudioCodec maps a plex codec metadata to a beq catalog codec name
func MapPlexToBeqAudioCodec(ctx context.Context, codecTitle, codecExtendTitle string) codecs.Codec {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Codecs from plex received",
		zap.String("codecTitle", codecTitle),
		zap.String("codecExtendTitle", codecExtendTitle),
	)

	// Titles are more likely to have atmos so check it first

	// Atmos logic
	atmosFlag := insensitiveContains(codecExtendTitle, codecs.CodecAtmos) || insensitiveContains(codecTitle, codecs.CodecAtmos)

	// check if contains DDP
	ddpFlag := containsDDP(codecTitle) || containsDDP(codecExtendTitle)

	log.Debug("Codec flags",
		zap.Bool("atmosFlag", atmosFlag),
		zap.Bool("ddpFlag", ddpFlag),
	)

	// if true and false, then Atmos
	if atmosFlag && !ddpFlag {
		return codecs.CodecAtmos
	}

	// if true and true, DD+ Atmos
	if atmosFlag && ddpFlag {
		return codecs.CodecDDPAtmos
	}

	// if not atmos and DD+, check later for DD+ Atmos, DD+ 7.1/5.1
	if !atmosFlag && ddpFlag {
		if insensitiveContains(codecTitle, "5.1") {
			return codecs.CodecDDPlusAtmos51Maybe
		}
		if insensitiveContains(codecTitle, "7.1") {
			return codecs.CodecDDPlusAtmos71Maybe
		}
	}

	// if False and false, then check others
	// TODO: simplify this like with jellyfin
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	case insensitiveContains(codecTitle, codecs.CodecTrueHD71) && insensitiveContains(codecExtendTitle, codecs.CodecTrueHD71):
		return codecs.CodecAtmosMaybe
	case insensitiveContains(codecTitle, codecs.CodecTrueHD71) && insensitiveContains(codecExtendTitle, codecs.CodecSurround71):
		return codecs.CodecAtmosMaybe
	// DTS:X
	case insensitiveContains(codecExtendTitle, codecs.CodecDTSX) || insensitiveContains(codecExtendTitle, codecs.CodecDTSXAlt):
		return codecs.CodecDTSX
	// DTS MA 7.1 containers but not DTS:X codecs
	case insensitiveContains(codecTitle, codecs.CodecDTSHDMA71) && !insensitiveContains(codecExtendTitle, codecs.CodecDTSX) && !insensitiveContains(codecExtendTitle, codecs.CodecDTSX):
		return codecs.CodecDTSHDMA71
	// DTS HA MA 5.1
	case insensitiveContains(codecTitle, codecs.CodecDTSHDMA51):
		return codecs.CodecDTSHDMA51
	// DTS 5.1
	case insensitiveContains(codecTitle, codecs.CodecDTS51):
		return codecs.CodecDTS51
	// TrueHD 5.1
	case insensitiveContains(codecTitle, codecs.CodecTrueHD51):
		return codecs.CodecTrueHD51
	// TrueHD 6.1
	case insensitiveContains(codecTitle, codecs.CodecTrueHD61):
		return codecs.CodecTrueHD61
	// DTS HRA
	case insensitiveContains(codecTitle, "DTS-HD HRA 7.1"):
		return codecs.CodecDTSHDHR71
	case insensitiveContains(codecTitle, "DTS-HD HRA 5.1"):
		return codecs.CodecDTSHDHR51
	// LPCM
	case insensitiveContains(codecTitle, codecs.CodecLPCM51):
		return codecs.CodecLPCM51
	case insensitiveContains(codecTitle, codecs.CodecLPCM71):
		return codecs.CodecLPCM71
	case insensitiveContains(codecTitle, codecs.CodecLPCM20):
		return codecs.CodecLPCM20
	case insensitiveContains(codecTitle, codecs.CodecAACStereo):
		return codecs.CodecAAC20
	case insensitiveContains(codecTitle, codecs.CodecAC351) || insensitiveContains(codecTitle, codecs.CodecEAC351):
		return codecs.CodecAC351
	case insensitiveContains(codecTitle, codecs.CodecEAC3) || insensitiveContains(codecExtendTitle, codecs.CodecEAC3):
		return codecs.CodecDDPlus
	default:
		return "Empty"
	}
}

// get the type of audio codec for BEQ purpose like atmos, dts-x, etc
func (c *PlexClient) GetAudioCodec(ctx context.Context, payload *events.Event) (codecs.Codec, error) {
	var plexAudioCodec codecs.Codec
	log := logger.GetLoggerFromContext(ctx)
	data, err := c.getMediaData(ctx, payload)
	if err != nil {
		return codecs.CodecAAC20, err
	}
	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
	// loop instead of index because of edge case with two or more video streams
	log.Debug("Data type",
		zap.String("type", fmt.Sprintf("%T", data)),
	)
	// TODO: better error handling
	if mc := data; mc.Video.Key != "" {
		// try to get Atmos from file because metadata with Truehd is usually misleading
		f := mc.Video.Media.Part.File
		if strings.Contains(strings.ToLower(f), string(codecs.CodecAtmos)) {
			log.Debug("Got atmos codec from filename")
			return MapPlexToBeqAudioCodec(ctx, f, f), nil
		}

		// index because of performance
		for i := range mc.Video.Media.Part.Stream {
			val := &mc.Video.Media.Part.Stream[i]
			if val.StreamType == "2" {
				log.Debug("Found codecs",
					zap.String("displayTitle", val.DisplayTitle),
					zap.String("extendedDisplayTitle", val.ExtendedDisplayTitle),
				)
				return MapPlexToBeqAudioCodec(ctx, val.DisplayTitle, val.ExtendedDisplayTitle), nil
			}
		}

		if plexAudioCodec == "" {
			log.Error("did not find codec from plex metadata",
				zap.String("title", mc.Video.Title),
				zap.Any("raw_data", mc.Video.Media.Part.Stream),
			)
			return "", errors.New("no codec found")
		}
	} else {
		return "", errors.New("invalid data type; mc.Video.Key is empty")
	}
	return plexAudioCodec, nil
}

// getEditionName tries to extract the edition from plex or file name. Assumes you have well named files
// Returned types, Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func (c *PlexClient) GetEdition(ctx context.Context, payload *events.Event) (editions.Edition, error) {
	data, err := c.getMediaData(ctx, payload)
	if err != nil {
		return editions.EditionUnknown, err
	}
	return getEdition(&data)
}

func getEdition(data *MediaContainer) (editions.Edition, error) {
	edition := data.Video.EditionTitle
	fileName := data.Video.Media.Part.File

	// First, check the edition from Plex metadata
	if edition != "" {
		mappedEdition := editions.MapSToEdition(edition)
		if mappedEdition != "" {
			return mappedEdition, nil
		}
		// If we couldn't map it, return it unknown
		return editions.EditionNone, errors.New("could not map edition")
	}

	// If no edition in metadata, try to extract from file name
	mappedEdition := editions.MapSToEdition(fileName)
	if mappedEdition != "" {
		return mappedEdition, nil
	}

	// no edition found, so its standard
	return editions.EditionNone, nil
}

func (c *PlexClient) makePlexReq(ctx context.Context, path string) ([]byte, error) {
	// Construct the URL with url.URL
	var u *url.URL
	log := logger.GetLoggerFromContext(ctx)

	// Add query parameters if needed
	if strings.Contains(path, "playback") {
		playerIP := config.GetHDMISyncPlayerIP(ctx)
		if playerIP == "" {
			return nil, errors.New("player IP not set in config")
		}
		log.Debug("Player IP",
			zap.String("playerIP", playerIP),
		)
		// this MUST use the CLIENT IP and 32500 port not server
		// god forbid plex makes any documentation for their APIs they dont want you using
		u = &url.URL{
			Scheme: c.Scheme,
			Host:   fmt.Sprintf("%s:%s", playerIP, "32500"),
			Path:   path,
		}
		params := url.Values{}
		// only X-Plex-Target-Client-Identifier MUST be sent and it MUST match the client machine id found in clientIP:32500/resources
		params.Add("X-Plex-Target-Client-Identifier", config.GetHDMISyncMachineIdentifier(ctx))
		// API docs says these must be sent, but thats not true at all
		// params.Add("commandID", "0")
		// params.Add("type", "video")
		u.RawQuery = params.Encode()

		log.Debug("using params for playback query",
			zap.String("params", u.RawQuery),
		)
	} else {
		u = &url.URL{
			Scheme: c.Scheme,
			Host:   fmt.Sprintf("%s:%s", c.URL, c.Port),
			Path:   path,
		}
	}
	// Create the request
	req, err := retryablehttp.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	// "X-Plex-Target-Client-Identifier"
	log.Debug("Plex: sending request",
		zap.String("url", u.String()),
	)
	// Execute the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when calling plex API: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("error closing response body",
				zap.Any("error", err),
			)
		}
	}()

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

func (c *PlexClient) ActionPause(ctx context.Context) error {
	return c.doPlaybackAction(ctx, events.ActionPause)
}

func (c *PlexClient) ActionPlay(ctx context.Context) error {
	return c.doPlaybackAction(ctx, events.ActionPlay)
}

func (c *PlexClient) ActionStop(ctx context.Context) error {
	return c.doPlaybackAction(ctx, events.ActionStop)
}

// DoPlaybackAction generic func to do playback - play, pause, stop
func (c *PlexClient) doPlaybackAction(ctx context.Context, action events.Action) error {
	log := logger.GetLoggerFromContext(ctx)
	s := fmt.Sprintf("/player/playback/%s", action)
	log.Debug("Plex: sending request",
		zap.String("action", string(action)),
		zap.String("url", s),
	)
	_, err := c.makePlexReq(ctx, s)

	return err
}
