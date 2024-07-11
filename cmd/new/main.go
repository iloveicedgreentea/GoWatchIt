package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/api"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/pkg/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/pkg/plex"
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
	// TODO: remove once templ is used
	r.Use(noCache())
	// time requests in debug mode
	if debug {
		r.Use(requestTimingMiddleware())
		log.Debug("Starting in debug mode with timing...")
	}

	/* ###############################
		Routes
	   ############################## */
	addRoutes(r)

	/* ###############################
	block until workers get ready
	############################## */
	factory := mediaplayer.NewMediaPlayerFactory()

	// Register Plex player
	plexPlayer := plex.NewPlexPlayer("plex.server.com", "32400")
	// TODO: use const
	factory.RegisterPlayer(PLEX_PLAYER, plexPlayer)

	log.Info("Waiting for workers to be ready...")
	// <-plexReady
	// <-minidspReady
	// <-jfReady
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
