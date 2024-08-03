package events

import (
	"context"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/models"
)

// RequestToEvent converts an http request to an event
func RequestToEvent(ctx context.Context, req *http.Request) (models.Event, error) {
	// TODO: read body and stuff, determine webhook type
	switch {
	case isPlexType(ctx, req):
		return processPlexWebhook(ctx, req)
	case isJellyfinType(req):
	case isHomeassistantType(req):
	}

	return models.Event{}, EventNotSupportedError{Message: "Event type not supported"}

}

// TODO: send to channel

func isPlexType(ctx context.Context, req *http.Request) bool {
	_, err := getMultipartPayload(ctx, req)
	return err == nil
}

func isJellyfinType(req *http.Request) bool {
	return false
}

func isHomeassistantType(req *http.Request) bool {
	return false
}
