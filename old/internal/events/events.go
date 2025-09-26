package events

import (
	"context"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

// RequestToEvent converts an http request to an event
func RequestToEvent(ctx context.Context, req *http.Request) (models.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	if req.Body == nil {
		return models.Event{}, EventNotSupportedError{Message: "Request body is empty"}
	}
	switch {
	case IsPlexType(ctx, req):
		log.Debug("Plex event")
		return processPlexWebhook(ctx, req)
	case IsJellyfinType(req):
		log.Debug("Jellyfin event")
		return parseJellyfinWebhook(ctx, req)
	case IsHomeassistantType(req):
		log.Debug("Homeassistant event")
		return parseHAWebhook(ctx, req)
	}

	return models.Event{}, EventNotSupportedError{Message: "Event type not supported"}
}

func IsPlexType(ctx context.Context, req *http.Request) bool {
	_, err := getMultipartPayload(ctx, req)
	return err == nil
}

func IsJellyfinType(req *http.Request) bool {
	return isJellyfinWebhook(req)
}

func IsHomeassistantType(req *http.Request) bool {
	// TODO: implement
	return false
}
