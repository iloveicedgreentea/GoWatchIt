package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/api"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/handlers"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func main() {
	log := logger.GetLogger()
	log.Debug("Started in debug mode...")
	r := gin.Default()

	// you can copy this schema to create event handlers for any service
	// create channel to receive jobs
	var plexChan = make(chan models.PlexWebhookPayload, 5)
	var minidspChan = make(chan models.MinidspRequest, 5)

	// ready signals
	plexReady := make(chan bool)
	minidspReady := make(chan bool)
	
	// run worker forever in background
	go handlers.PlexWorker(plexChan, plexReady)
	go handlers.MiniDspWorker(minidspChan, minidspReady)

	// healthcheck
	r.GET("/health", handlers.ProcessHealthcheckWebhookGin)
	// Add plex webhook handler
	// TODO: split out non plex specific stuff into a library
	r.POST("/plexwebhook", func(c *gin.Context) {
		handlers.ProcessWebhook(plexChan, c)
	})
	r.POST("/minidspwebhook", func(c *gin.Context) {
		handlers.ProcessMinidspWebhook(minidspChan, c)
	})
	// TODO: add generic webhook endpoint, maybe mqtt?

	// wait for workers to get ready
	// TODO implement signal checking, error chan, etc
	<-plexReady
	<-minidspReady
	log.Info("All workers are ready.")

	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html") 
	})
    // Register routes
    api.RegisterRoutes(r)
	// TODO: Engine.SetTrustedProxies(nil)
	port := config.GetString("main.listenPort")
	if port == "" {
		log.Fatal("error getting port")
	}
	log.Infof("Starting server on port %v", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
