package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func parsePlexMultipartForm(payload []string) (models.PlexWebhookPayload, error) {
	var pwhPayload models.PlexWebhookPayload

	err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// unmarshall error
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return pwhPayload, errors.New(msg)

		default:
			return pwhPayload, err
		}
	}

	return pwhPayload, nil
}

// TODO: test this
// getMultipartPayload gets the payload from the multipart form and returns if ok
func getMultipartPayload(ctx context.Context, request *http.Request) ([]string, error) {
	log := logger.GetLoggerFromContext(ctx)
	if err := request.ParseMultipartForm(0); err != nil {
		log.Error("Error parsing multipart form",
			slog.String("error", err.Error()),
		)
		return []string{}, fmt.Errorf("invalid multipart form: %s", err)
	}

	payload, ok := request.MultipartForm.Value["payload"]
	if !ok {
		log.Error("Error parsing multipart form",
			slog.String("error", "no payload found"),
		)
		return []string{}, errors.New("no payload found in request")
	}

	return payload, nil
}

// Sends the payload to the channel for background processing
// TODO: this takes an HTTP request and sends it to a channel
func processPlexWebhook(ctx context.Context, request *http.Request) (models.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	payload, err := getMultipartPayload(ctx, request)
	if err != nil {
		return models.Event{}, fmt.Errorf("error getting payload: %s", err)
	}

	// parse the payload
	log.Debug("decoding payload")
	decodedPayload, err := parsePlexMultipartForm(payload)
	if err != nil {
		return models.Event{}, fmt.Errorf("error decoding payload: %s", err)
	}

	log.Debug("Got a request from UUID: %s",
		slog.String("player_uuid", decodedPayload.Player.UUID),
	)

	mediaType := decodedPayload.Metadata.Type

	log.Debug("Processed Webhook",
		slog.String("media_type", mediaType),
		slog.String("media_title", decodedPayload.Metadata.Title),
	)
	// check filter for user if not blank
	userID := config.GetString("plex.ownerNameFilter")
	// only respond to events on a particular account if you share servers and only for movies and shows
	// TODO: decodedPayload.Account.Title seems to always map to server owner not player account
	if len(userID) == 0 || strings.EqualFold(decodedPayload.Account.Title, userID) {
		if strings.EqualFold(mediaType, string(models.MediaTypeMovie)) || strings.EqualFold(mediaType, string(models.MediaTypeShow)) {
			log.Debug("adding item to plexChan")

		} else {
			log.Debug("Media type not supported",
				slog.String("media_type", mediaType),
			)
		}
	} else {
		// TODO: this seems to be hitting even when the filter matches
		log.Debug("userID does not match filter",
			slog.String("account_title", decodedPayload.Account.Title),
			slog.String("filter", userID),
		)

	}

	return models.Event{}, nil
}
