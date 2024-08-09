package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/events"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func main() {
	ctx := context.Background()
	log := logger.GetLoggerFromContext(ctx)
	logger.AddLoggerToContext(ctx, log)

	eventChan := make(chan models.Event)

	// TODO: handler and shit
	// go eventHandler(ctx)

	// Process request
	// TODO: obviously do this in the API worker async
	event, err := events.RequestToEvent(ctx, &http.Request{})
	if err != nil {
		logger.Fatal("Error processing event",
			slog.Any("error", err),
		)
	}

	log.Info("Found event",
		slog.Any("event", event),
	)

	// send event to chan
	// TODO: func inside the worker
	eventChan <- event
}
