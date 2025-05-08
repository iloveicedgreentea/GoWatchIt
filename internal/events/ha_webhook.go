package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func parseHAWebhook(ctx context.Context, req *http.Request) (models.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	var webhook models.HomeAssistantWebhookPayload

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return models.Event{}, err
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Error("Failed to close request body",
				slog.Any("error", err),
			)
		}
	}()

	// Unmarshal the JSON into the webhook struct
	err = json.Unmarshal(body, &webhook)
	if err != nil {
		return models.Event{}, err
	}

	c, err := homeassistant.NewClient()
	if err != nil {
		return models.Event{}, fmt.Errorf("failed to create homeassistant client: %w", err)
	}

	// TODO: get current state
	state, err := c.ReadState(webhook.EntityID, models.HomeAssistantEntityMediaPlayer)
	if err != nil {
		return models.Event{}, fmt.Errorf("failed to read state: %w", err)
	}

	return models.Event{
		Action:    state.State.StateToAction(),
		Client:    c,
		EventType: models.EventTypeHomeAssistant,
		Metadata: models.Metadata{
			Title: state.Attributes.MediaTitle,
			TMDB:  state.Attributes.TMDB,
			Type:  models.MediaType(state.Attributes.MediaContentType),
			// TODO: modfiy search if no year, search by tmdb filter
			// TODO: need to get codec
			// TODO: add an option to assume codec?
		},
	}, nil
}
