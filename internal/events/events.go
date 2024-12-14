package events

import (
	"context"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/models"
)

// RequestToEvent converts an http request to an event
func RequestToEvent(ctx context.Context, req *http.Request) (models.Event, error) {
	if req.Body == nil {
		return models.Event{}, EventNotSupportedError{Message: "Request body is empty"}
	}
	switch {
	case isPlexType(ctx, req):
		return processPlexWebhook(ctx, req)
	case isJellyfinType(req):
		return parseJellyfinWebhook(ctx, req)
	case isHomeassistantType(req):
		return models.Event{}, nil
	}

	return models.Event{}, EventNotSupportedError{Message: "Event type not supported"}
}

func isPlexType(ctx context.Context, req *http.Request) bool {
	_, err := getMultipartPayload(ctx, req)
	return err == nil
}

func isJellyfinType(req *http.Request) bool {
	return isJellyfinWebhook(req)
}

func isHomeassistantType(req *http.Request) bool {
	// TODO: implement
	return false
}
