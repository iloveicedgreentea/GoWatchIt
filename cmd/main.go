package main

import (
	"fmt"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/handlers"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
)

func main() {
	log := logger.GetLogger()
	log.Debug("Started in debug mode...")
	mux := http.NewServeMux()

	// you can copy this schema to create event handlers for any service
	// create channel to receive jobs
	var plexChan = make(chan models.PlexWebhookPayload, 5)
	var minidspChan = make(chan models.MinidspRequest, 5)

	// run worker forever in background
	go handlers.PlexWorker(plexChan)
	go handlers.MiniDspWorker(minidspChan)

	// pass the chan to the handlers
	plexWh := handlers.ProcessWebhook(plexChan)
	dspWh := handlers.ProcessMinidspWebhook(minidspChan)

	// healthcheck
	health := handlers.ProcessHealthcheckWebhook()

	// Add plex webhook handler
	// TODO: split out non plex specific stuff into a library
	mux.Handle("/plexwebhook", plexWh)
	// TODO: add generic webhook endpoint, maybe mqtt?

	// minidsp
	mux.Handle("/minidspwebhook", dspWh)

	// healthcheck
	mux.Handle("/health", health)

	log.Info("Starting server")
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetString("main.listenPort")), mux)
	log.Fatal(err)
}
