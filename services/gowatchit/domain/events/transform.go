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

// RequestToEvent converts http to Event type
func RequestToEvent(ctx context.Context, player mediaplayer.MediaPlayer, req *http.Request) (*events.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	if req.Body == nil {
		return &events.Event{}, EventNotSupportedError{Message: "Request body is empty"}
	}

	if player == nil {
		return &events.Event{}, fmt.Errorf("player is nil")
	}

	switch player.(type) {
	case *plex.PlexClient:
		log.Debug("Using Plex event parser")
		return processPlexWebhook(ctx, req)
	default:
		log.Warn("Media player not supported for event parsing",
			zap.String("player", fmt.Sprintf("%T", player)),
		)
		return &events.Event{}, EventNotSupportedError{Message: "Media player not supported for event parsing"}
	}
	// TODO: JF event
	// TODO: HA event
}
