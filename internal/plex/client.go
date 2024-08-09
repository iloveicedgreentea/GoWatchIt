// client.go implements the API for plex itself
package plex

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/internal/utils"
	"github.com/iloveicedgreentea/go-plex/models"
)

// Stuff to interface directly with Plex
// of course their api is undocumented and worst of all, in xml. I had to piece it together reading various unofficial API implementations or trial and error

type PlexClient struct {
	URL        string
	Port       string
	HTTPClient http.Client
	MachineID  string
	ClientIP   string
	MediaType  string
}

// this needs to implement MediaAPIClient
var _ mediaplayer.MediaAPIClient = (*PlexClient)(nil)

// return a new instance of a plex client
func NewClient(scheme, url, port string) (*PlexClient, error) {
	// remove scheme
	if !utils.ValidateHttpScheme(scheme) {
		return nil, fmt.Errorf("invalid http scheme: %s", scheme)
	}
	url = strings.Replace(url, scheme, "", -1)
	return &PlexClient{
		URL:  url,
		Port: port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

// unmarshal xml into a struct
func parseMediaContainer(payload []byte) (models.MediaContainer, error) {
	var data models.MediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func parseSessionMediaContainer(payload []byte) (models.SessionMediaContainer, error) {
	var data models.SessionMediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *PlexClient) getRunningSession(ctx context.Context) (models.SessionMediaContainer, error) {
	log := logger.GetLoggerFromContext(ctx)
	// Get session object
	var data models.SessionMediaContainer
	var err error

	// loop until not empty for 30s
	for i := 0; i < 30; i++ {
		res, err := c.makePlexReq(ctx, string(APIStatusSession))
		if err != nil {
			return models.SessionMediaContainer{}, err
		}
		// if no response, keep trying
		if len(res) == 0 {
			log.Debug("Plex session empty, waiting",
				slog.Int("time_remaining", 30-i),
			)

			time.Sleep(1 * time.Second)
			continue
		}
		data, err = parseSessionMediaContainer(res)
		if err != nil {
			return models.SessionMediaContainer{}, err
		}

		break
	}

	return data, err
}

// GetCodecFromSession gets the codec from a running session
func (c *PlexClient) GetCodecFromSession(ctx context.Context, uuid string) (models.CodecName, error) {
	sess, err := c.getRunningSession(ctx)
	log := logger.GetLoggerFromContext(ctx)
	if err != nil {
		return "", err
	}
	// log.Debugf("Session data: %#v", sess.Video)
	// filter by uuid
	// try up to 15 times until session is active. webhook sends before session is ready
	for i := 0; i < 15; i++ {
		for _, video := range sess.Video {
			log.Debug("Machine identifier",
				slog.String("identifier", video.Player.MachineIdentifier),
			)
			if video.Player.MachineIdentifier == uuid {
				log.Debug("Found session matching uuid",
					slog.String("uuid", uuid),
				)
				for _, stream := range video.Media.Part.Stream {
					log.Debug("Stream data",
						slog.String("data", fmt.Sprintf("%#v", stream)),
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
func (c *PlexClient) getMediaData(ctx context.Context, payload models.Event) (models.MediaContainer, error) {
	libraryKey := payload.Metadata.Key
	res, err := c.makePlexReq(ctx, libraryKey)
	if err != nil {
		return models.MediaContainer{}, err
	}

	return parseMediaContainer(res)
}

func insensitiveContains(s string, sub models.CodecName) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(string(sub)))
}

// check if its DD+ codec
func containsDDP(s string) bool {
	//English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	ddPlusNames := []models.CodecName{models.CodecDDP, models.CodecDDPlus, models.CodecEAC3, models.CodecEAC3Alt}
	for _, name := range ddPlusNames {
		if insensitiveContains(strings.ToLower(s), name) {
			return true
		}
	}

	return false
}

// mapPlexToBeqAudioCodec maps a plex codec metadata to a beq catalog codec name
func MapPlexToBeqAudioCodec(ctx context.Context, codecTitle, codecExtendTitle string) models.CodecName {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Codecs from plex received",
		slog.String("codecTitle", codecTitle),
		slog.String("codecExtendTitle", codecExtendTitle),
	)

	// Titles are more likely to have atmos so check it first

	// Atmos logic
	atmosFlag := insensitiveContains(codecExtendTitle, models.CodecAtmos) || insensitiveContains(codecTitle, models.CodecAtmos)

	// check if contains DDP
	ddpFlag := containsDDP(codecTitle) || containsDDP(codecExtendTitle)

	log.Debug("Codec flags",
		slog.Bool("atmosFlag", atmosFlag),
		slog.Bool("ddpFlag", ddpFlag),
	)

	// if true and false, then Atmos
	if atmosFlag && !ddpFlag {
		return models.CodecAtmos
	}

	// if true and true, DD+ Atmos
	if atmosFlag && ddpFlag {
		return models.CodecDDPAtmos
	}

	// if not atmos and DD+, check later for DD+ Atmos, DD+ 7.1/5.1
	if !atmosFlag && ddpFlag {
		if insensitiveContains(codecTitle, "5.1") {
			return models.CodecDDPAtmos5Maybe
		}
		if insensitiveContains(codecTitle, "7.1") {
			return models.CodecDDPAtmos7Maybe
		}
	}

	// if False and false, then check others
	// TODO: simplify this like with jellyfin
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	case insensitiveContains(codecTitle, models.CodecTrueHD71) && insensitiveContains(codecExtendTitle, models.CodecTrueHD71):
		return models.CodecAtmosMaybe
	case insensitiveContains(codecTitle, models.CodecTrueHD71) && insensitiveContains(codecExtendTitle, models.CodecSurround71):
		return models.CodecAtmosMaybe
	// DTS:X
	case insensitiveContains(codecExtendTitle, models.CodecDTSX) || insensitiveContains(codecExtendTitle, models.CodecDTSXAlt):
		return models.CodecDTSX
	// DTS MA 7.1 containers but not DTS:X codecs
	case insensitiveContains(codecTitle, models.CodecDTSHDMA71) && !insensitiveContains(codecExtendTitle, models.CodecDTSX) && !insensitiveContains(codecExtendTitle, models.CodecDTSX):
		return models.CodecDTSHDMA71
	// DTS HA MA 5.1
	case insensitiveContains(codecTitle, models.CodecDTSHDMA51):
		return models.CodecDTSHDMA51
	// DTS 5.1
	case insensitiveContains(codecTitle, models.CodecDTS51):
		return models.CodecDTS51
	// TrueHD 5.1
	case insensitiveContains(codecTitle, models.CodecTrueHD51):
		return models.CodecTrueHD51
	// TrueHD 6.1
	case insensitiveContains(codecTitle, models.CodecTrueHD61):
		return models.CodecTrueHD61
	// DTS HRA
	case insensitiveContains(codecTitle, "DTS-HD HRA 7.1"):
		return models.CodecDTSHDHR71
	case insensitiveContains(codecTitle, "DTS-HD HRA 5.1"):
		return models.CodecDTSHDHR51
	// LPCM
	case insensitiveContains(codecTitle, models.CodecLPCM51):
		return models.CodecLPCM51
	case insensitiveContains(codecTitle, models.CodecLPCM71):
		return models.CodecLPCM71
	case insensitiveContains(codecTitle, models.CodecLPCM20):
		return models.CodecLPCM20
	case insensitiveContains(codecTitle, models.CodecAACStereo):
		return models.CodecAAC20
	case insensitiveContains(codecTitle, models.CodecAC351) || insensitiveContains(codecTitle, models.CodecEAC351):
		return models.CodecAC351
	case insensitiveContains(codecTitle, models.CodecEAC3) || insensitiveContains(codecExtendTitle, models.CodecEAC3):
		return models.CodecDDPlus
	default:
		return "Empty"
	}

}

// get the type of audio codec for BEQ purpose like atmos, dts-x, etc
func (c *PlexClient) GetAudioCodec(ctx context.Context, payload models.Event) (models.CodecName, error) {
	var plexAudioCodec models.CodecName
	log := logger.GetLoggerFromContext(ctx)
	data, err  := c.getMediaData(ctx, payload)
	if err != nil {
		return models.CodecAAC20, err
	}
	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
	// loop instead of index because of edge case with two or more video streams
	log.Debug("Data type",
		slog.String("type", fmt.Sprintf("%T", data)),
	)
	// TODO: better error handling
	if mc := data; len(mc.Video.Key) > 0 {
		// try to get Atmos from file because metadata with Truehd is usually misleading
		f := mc.Video.Media.Part.File
		if strings.Contains(strings.ToLower(f), string(models.CodecAtmos)) {
			log.Debug("Got atmos codec from filename")
			return MapPlexToBeqAudioCodec(ctx, f, f), nil
		}
		for _, val := range mc.Video.Media.Part.Stream {
			if val.StreamType == "2" {
				log.Debug("Found codecs",
					slog.String("displayTitle", val.DisplayTitle),
					slog.String("extendedDisplayTitle", val.ExtendedDisplayTitle),
				)
				return MapPlexToBeqAudioCodec(ctx, val.DisplayTitle, val.ExtendedDisplayTitle), nil
			}
		}

		if plexAudioCodec == "" {
			log.Error("did not find codec from plex metadata",
				slog.String("title", mc.Video.Title),
				slog.Any("raw_data", mc.Video.Media.Part.Stream),
			)
			return "", errors.New("no codec found")
		}
	} else {
		return "", errors.New("invalid data type")
	}
	return plexAudioCodec, nil
}

// TODO: rename
// GetPlexMovieDb is used because of the Client interface
func (c *PlexClient) GetPlexMovieDb(payload interface{}) string {
	return ""
}

// getEditionName tries to extract the edition from plex or file name. Assumes you have well named files
// Returned types, Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func (p *PlexClient) GetEdition(ctx context.Context, payload models.Event) (models.Edition, error) {
	data, err  := p.getMediaData(ctx, payload)
	if err != nil {
		return models.EditionUnknown, err
	}
	edition := strings.ToLower(data.Video.EditionTitle)
	fileName := strings.ToLower(data.Video.Media.Part.File)

	// First, check the edition from Plex metadata
	if edition != "" {
		mappedEdition := mapToEdition(edition)
		if mappedEdition != "" {
			return mappedEdition, nil
		}
		// If we couldn't map it, return it unknown
		return models.EditionUnknown, errors.New("could not map edition")
	}

	// If no edition in metadata, try to extract from file name
	mappedEdition := mapToEdition(fileName)
	if mappedEdition != "" {
		return mappedEdition, nil
	}

	// no edition found, so its standard
	return models.EditionNone, nil
}

func mapToEdition(s string) models.Edition {
	switch {
	case strings.Contains(s, "extended"):
		return models.EditionExtended
	case strings.Contains(s, "unrated"):
		return models.EditionUnrated
	case strings.Contains(s, "theatrical"):
		return models.EditionTheatrical
	case strings.Contains(s, "ultimate"):
		return models.EditionUltimate
	case strings.Contains(s, "director"):
		return models.EditionDirectorsCut
	case strings.Contains(s, "criterion"):
		return models.EditionCriterion
	default:
		return models.EditionUnknown
	}
}

func (c *PlexClient) makePlexReq(ctx context.Context, path string) ([]byte, error) {
	// Construct the URL with url.URL
	var u *url.URL
	log := logger.GetLoggerFromContext(ctx)

	// Add query parameters if needed
	if strings.Contains(path, "playback") {
		playerIP := config.GetString("signal.playerip")
		if playerIP == "" {
			return nil, errors.New("player IP not set in config")
		}
		log.Debug("Player IP",
			slog.String("playerIP", playerIP),
		)
		// this MUST use the CLIENT IP and 32500 port not server
		// god forbid plex makes any documentation for their APIs they dont want you using
		u = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%s", playerIP, "32500"),
			Path:   path,
		}
		params := url.Values{}
		// only X-Plex-Target-Client-Identifier MUST be sent and it MUST match the client machine id found in clientIP:32500/resources
		params.Add("X-Plex-Target-Client-Identifier", config.GetString("signal.playermachineidentifier"))
		// API docs says these must be sent, but thats not true at all
		// params.Add("commandID", "0")
		// params.Add("type", "video")
		u.RawQuery = params.Encode()

		log.Debug("using params for playback query",
			slog.String("params", u.RawQuery),
		)
	} else {
		u = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%s", c.URL, c.Port),
			Path:   path,
		}
	}
	// Create the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("X-Plex-Target-Client-Identifier", c.MachineID)
	log.Debug("Plex: sending request",
		slog.String("url", u.String()),
	)
	// Execute the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when calling plex API: %v", err)
	}
	defer res.Body.Close()

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

// DoPlaybackAction generic func to do playback - play, pause, stop
func (c *PlexClient) DoPlaybackAction(ctx context.Context, action models.Action) error {
	log := logger.GetLoggerFromContext(ctx)
	s := fmt.Sprintf("/player/playback/%s", action)
	log.Debug("Plex: sending request",
		slog.String("action", string(action)),
		slog.String("url", s),
	)
	_, err := c.makePlexReq(ctx, s)

	return err
}
