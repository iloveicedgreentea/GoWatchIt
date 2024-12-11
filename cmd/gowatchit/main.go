package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"

	"github.com/iloveicedgreentea/go-plex/models"
)

// static files are cached which causes issues displaying configs
func noCache() gin.HandlerFunc {
	// TODO: cache specific resources just not
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Next()
	}
}

func main() {
	// init context
	ctx := context.Background()
	log := logger.GetLoggerFromContext(ctx)

	// reuse logger object in calls
	logger.AddLoggerToContext(ctx, log)

	log.Info("Starting up please wait until the server is ready...")

	// set program to debug mode
	debug := os.Getenv("LOG_LEVEL") == "debug"
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		baseDir = "."
	}

	// TODO: use const for sql file location
	// TODO: file needs to be docker-compatible
	// Create the database connection
	log.Info("Connecting to the database...")
	db, err := database.GetDB(fmt.Sprintf("%s/db.sqlite3", baseDir))
	if err != nil {
		logger.Fatal("Failed to connect to the database: ", err)
	}
	// close db when done
	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal("Failed to close the database: ", err)
		}
	}()

	// create or update tables
	log.Info("Running migrations...")
	err = database.RunMigrations(db)
	if err != nil {
		logger.Fatal("Failed to run migrations: ", err)
	}

	// init the config manager
	log.Info("Initializing config...")
	err = config.InitConfig(db)
	if err != nil {
		logger.Fatal("Failed to run init config: ", err)
	}

	// init router
	log.Info("Initializing router...")
	router := gin.Default()
	// do not cache static files
	router.Use(noCache()) // TODO: scope to specific routes

	// init event channel
	log.Info("Creating workers...")
	// all webhook events are sent to this channel
	eventChan := make(chan models.Event)

	router.Static("/web", "./web")
	router.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// register routes
	RegisterRoutes(router, eventChan)
	err = router.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatal("Failed to set trusted proxies: ", err)
	}

	// init clients
	log.Info("Creating clients...")
	beqClient, err := ezbeq.NewClient()
	if err != nil {
		log.Error("Error creating beq client",
			slog.Any("error", err),
		)
		return
	}

	homeAssistantClient, err := homeassistant.NewClient()
	if err != nil {
		log.Error("Error creating HA client",
			slog.Any("error", err),
		)
		return
	}

	// run event loop in background
	go eventHandler(ctx, eventChan, beqClient, homeAssistantClient)

	// init the router
	port := config.GetMainListenPort()
	if port == "" {
		port = "9999"
	}
	log.Info("Starting server")
	log.Debug("Listening on port", slog.String("port", port))
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal(err.Error())
	}
}
