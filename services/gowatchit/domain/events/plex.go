package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/events"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/pkg/plex"
	"go.uber.org/zap"
)

// getMultipartPayload gets the payload from the multipart form and returns if ok
func getMultipartPayload(ctx context.Context, request *http.Request) ([]string, error) {
	log := logger.GetLoggerFromContext(ctx)
	if err := request.ParseMultipartForm(0); err != nil {
		log.Error("Error parsing multipart form",
			zap.Error(err),
		)
		return []string{}, fmt.Errorf("invalid multipart form: %s", err)
	}

	payload, ok := request.MultipartForm.Value["payload"]
	if !ok {
		log.Error("Error parsing multipart form",
			zap.String("error", "no payload found"),
		)
		return []string{}, errors.New("no payload found in request")
	}

	return payload, nil
}

func parsePlexMultipartForm(payload []string) (plex.PlexWebhookPayload, error) {
	var pwhPayload plex.PlexWebhookPayload

	err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// unmarshall error
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return pwhPayload, errors.New(msg + " " + err.Error())

		default:
			return pwhPayload, err
		}
	}

	return pwhPayload, nil
}

// Sends the payload to the channel for background processing
func processPlexWebhook(ctx context.Context, request *http.Request) (*events.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	payload, err := getMultipartPayload(ctx, request)
	if err != nil {
		return &events.Event{}, fmt.Errorf("error getting payload: %s", err)
	}

	// parse the payload
	log.Debug("decoding payload")
	decodedPayload, err := parsePlexMultipartForm(payload)
	if err != nil {
		return &events.Event{}, fmt.Errorf("error decoding payload: %s", err)
	}

	log.Debug("Got a request from UUID: %s",
		zap.String("player_uuid", decodedPayload.Player.UUID),
	)

	// check all filters
	checkRes, err := plex.CheckAllFilters(ctx, decodedPayload)
	if !checkRes {
		log.Warn("filters did not match",
			zap.Error(err),
		)
		return &events.Event{}, err
	}
	mediaType := decodedPayload.Metadata.Type
	log.Debug("Processed Webhook",
		zap.String("media_type", mediaType),
		zap.String("media_title", decodedPayload.Metadata.Title),
		zap.String("uuid", decodedPayload.Player.UUID),
		zap.String("username", decodedPayload.Account.Title),
	)

	// check if TV BEQ is enabled
	if strings.EqualFold(mediaType, string(plex.MediaTypeShow)) && !config.IsBeqTVEnabled(ctx) {
		log.Warn("TV BEQ is disabled",
			zap.String("media_type", mediaType),
		)
		return &events.Event{}, errors.New("TV BEQ is disabled but episode type found")
	}

	var action events.Action
	switch decodedPayload.Event {
	case string(plex.PlexActionPlay):
		action = events.ActionPlay
	case string(plex.PlexActionStop):
		action = events.ActionStop
	case string(plex.PlexActionPause):
		action = events.ActionPause
	// Pressing the 'resume' button in plex UI is media.play
	case string(plex.PlexActionResume):
		action = events.ActionResume
	case string(plex.PlexActionScrobble):
		action = events.ActionScrobble
	default:
		log.Debug("Received unsupported event",
			zap.String("event", decodedPayload.Event),
		)
		return &events.Event{}, EventNotSupportedError{Message: "event type not supported"}
	}

	var tmdb string
	// extract the tmdb ID from the GUID0 field
	for _, model := range decodedPayload.Metadata.GUID0 {
		if strings.Contains(model.ID, "tmdb") {
			log.Debug("getPlexMovieDb: Got tmdb ID from plex",
				zap.String("id", model.ID),
			)
			tmdb = strings.Split(model.ID, "tmdb://")[1]
		}
	}

	return &events.Event{
		Action:      action,
		EventType:   events.EventTypePlex,
		User:        decodedPayload.User,
		Owner:       decodedPayload.Owner,
		AccountID:   strconv.Itoa(decodedPayload.Account.ID),
		ServerUUID:  decodedPayload.Server.UUID,
		PlayerUUID:  decodedPayload.Player.UUID,
		PlayerTitle: decodedPayload.Player.Title,
		ServerTitle: decodedPayload.Server.Title,
		PlayerIP:    decodedPayload.Player.PublicAddress,
		Metadata: events.Metadata{
			TMDB:                tmdb,
			Year:                decodedPayload.Metadata.Year,
			LibrarySectionType:  decodedPayload.Metadata.LibrarySectionType,
			Key:                 decodedPayload.Metadata.Key,
			Type:                events.MediaType(decodedPayload.Metadata.Type),
			Title:               decodedPayload.Metadata.Title,
			LibrarySectionTitle: decodedPayload.Metadata.LibrarySectionTitle,
			LibrarySectionID:    decodedPayload.Metadata.LibrarySectionID,
			LibrarySectionKey:   decodedPayload.Metadata.LibrarySectionKey,
		},
	}, nil
}
