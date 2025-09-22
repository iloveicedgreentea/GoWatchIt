package events

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/events"
)

func ProcessWebhook(ctx context.Context, req *http.Request) error {
	_, err := requestToEvent(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to convert webhook to event: %w", err)
	}

	// TODO: route event to correct cmd handler

	return nil
}

// RequestToEvent converts http to Event type
func requestToEvent(ctx context.Context, req *http.Request) (*events.Event, error) {
	log := logger.GetLoggerFromContext(ctx)
	if req.Body == nil {
		return &events.Event{}, EventNotSupportedError{Message: "Request body is empty"}
	}

	switch {
	case IsPlexType(ctx, req):
		log.Debug("Plex")
		return processPlexWebhook(ctx, req)
	}
	// TODO: JF event
	// TODO: HA event
	return &events.Event{}, EventNotSupportedError{Message: "Event type not supported"}
}

// IsPlexType checks if it can extract a multipart payload. If true, its a plex payload
func IsPlexType(ctx context.Context, req *http.Request) bool {
	_, err := getMultipartPayload(ctx, req)
	return err == nil
}
