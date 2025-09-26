package events

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iloveicedgreentea/gowatchit/pkg/events"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/pkg/plex"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/mediaplayer"
	"go.uber.org/zap"
)

func ProcessWebhook(ctx context.Context, player mediaplayer.MediaPlayer, req *http.Request) error {
	log := logger.GetLoggerFromContext(ctx)
	if req == nil {
		log.Error("Request is nil")
		return fmt.Errorf("request is nil")
	}
	if player == nil {
		log.Error("Media player is nil")
		return fmt.Errorf("media player is nil")
	}
	event, err := requestToEvent(ctx, player, req)
	if err != nil {
		return fmt.Errorf("failed to convert webhook to event: %w", err)
	}
	if event == nil {
		log.Error("Event is nil after conversion")
		return fmt.Errorf("event is nil after conversion")
	}

	log.Debug("Received event",
		zap.Any("event", event),
	)
	// TODO: CONTINUE HERE - from event, trigger a load
	// TODO: support hdmi sync

	return nil
}

// RequestToEvent converts http to Event type
func requestToEvent(ctx context.Context, player mediaplayer.MediaPlayer, req *http.Request) (*events.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	if req.Body == nil {
		return &events.Event{}, EventNotSupportedError{Message: "Request body is empty"}
	}

	switch player.(type) {
	case *plex.PlexClient:
		log.Debug("Using Plex event parser")
		return processPlexWebhook(ctx, req)
	}
	// TODO: JF event
	// TODO: HA event
	return &events.Event{}, EventNotSupportedError{Message: "Event type not supported"}
}
