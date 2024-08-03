package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/events"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
)

func main() {
	ctx := context.Background()
	log := logger.GetLoggerFromContext(ctx)
	logger.AddLoggerToContext(ctx, log)

	// TODO: handler and shit

	// Process request
	// TODO: obviously do this in a handler
	event, err := events.RequestToEvent(ctx, &http.Request{})
	if err != nil {
		logger.Fatal("Error processing event",
			slog.Any("error", err),
		)
	}

	log.Info("Found event",
		slog.Any("event", event),
	)
}
