package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/api"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/handlers"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

// static files are cached which causes issues
func noCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Next()
	}
}

// requestTimingMiddleware captures the duration of each request to time things
func requestTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log or store the duration
		// For example, logging the duration:
		log.Infof("Request %s took %v", c.Request.RequestURI, duration)
	}
}

func main() {
	/*
		###############################
		Setups
		############################## */

	log.Info("Starting up please wait until the server is ready...")
	debug := os.Getenv("LOG_LEVEL") == "debug"
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// do not cache static files
	r.Use(noCache())
	// time requests in debug mode
	if debug {
		r.Use(requestTimingMiddleware())
		log.Debug("Starting in debug mode with timing...")
	}

	log.Info("Checking if a config exists...")
	_, err := api.GetConfigPath()
	if err != nil {
		log.Info("Config not found. Creating a new config file...")
		err = api.CreateConfig(&gin.Context{})
		if err != nil {
			log.Fatalf("Unable to create config file: %v", err)
		}
	}

	// you can copy this schema to create event handlers for any service
	// create channel to receive jobs
	var plexChan = make(chan models.PlexWebhookPayload, 5)
	var minidspChan = make(chan models.MinidspRequest, 5)
	var jfChan = make(chan models.JellyfinWebhook, 5)
	var webhookChan = make(chan models.Webhook, 5)

	// ready signals
	plexReady := make(chan bool)
	minidspReady := make(chan bool)
	jfReady := make(chan bool)
	webhookReady := make(chan bool)

	// run worker forever in background
	/*
		###############################
		handlers
		############################## */
	go handlers.PlexWorker(plexChan, plexReady)
	go handlers.MiniDspWorker(minidspChan, minidspReady)
	go handlers.JellyfinWorker(jfChan, jfReady)
	go handlers.WebhookWorker(webhookChan, webhookReady)

	/* ###############################
		Routes
	   ############################## */
	// healthcheck
	r.GET("/health", handlers.ProcessHealthcheckWebhookGin)
	// non-plex based webhook
	r.POST("/webhook", func(c *gin.Context) {
		log.Debug("plexwebhook received")
		handlers.ProcessPlainWebhook(c.Request.Context(), webhookChan, c)
	})
	// Add plex webhook handler
	r.POST("/plexwebhook", func(c *gin.Context) {
		log.Debug("plexwebhook received")
		handlers.ProcessWebhook(c.Request.Context(), plexChan, c)
	})
	r.POST("/minidspwebhook", func(c *gin.Context) {
		log.Debug("minidspwebhook received")
		handlers.ProcessMinidspWebhook(minidspChan, c)
	})
	r.POST("/jellyfinwebhook", func(c *gin.Context) {
		handlers.ProcessJfWebhook(jfChan, c)
	})
	r.Static("/assets", "./assets")
	r.GET("/config-exists", api.ConfigExists)
	r.GET("/get-config", api.GetConfig)
	r.POST("/save-config", api.SaveConfig)
	// TODO: add generic webhook endpoint, maybe mqtt?

	/*
		###############################
		block until workers get ready
		############################## */
	log.Info("Waiting for workers to be ready...")
	<-plexReady
	<-minidspReady
	<-jfReady
	log.Info("All workers are ready.")

	r.Static("/web", "./web")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Register routes
	api.RegisterRoutes(r)
	r.SetTrustedProxies(nil)
	port := config.GetString("main.listenPort")
	if port == "" {
		port = "9999"
	}
	log.Infof("Starting server on port %v", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
