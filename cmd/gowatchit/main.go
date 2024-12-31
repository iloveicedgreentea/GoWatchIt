package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
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
	logger.PanicLogger(func() {
		// init context
		ctx := context.Background()

		// reuse logger object in calls
		err := logger.InitLoggerFile()
		if err != nil {
			logger.Fatal("Failed to initialize the logger: ", err)
		}
		defer func() {
			if err := logger.CleanupLogger(); err != nil {
				logger.Fatal("Failed to close the logger: ", err)
			}
		}()

		log := logger.GetLogger()
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

		dbDir := fmt.Sprintf("%s/db.sqlite3", baseDir)
		log.Debug("Base directory", slog.String("dbDir", dbDir))

		// Create the database connection
		log.Info("Connecting to the database...")
		db, err := database.GetDB(dbDir)
		if err != nil {
			logger.Fatal("Failed to connect to the database: ", err)
		}
		if db == nil {
			logger.Fatal("db is nil")
		}

		// close db when done
		defer func() {
			if err := logger.CleanupLogger(); err != nil {
				logger.Fatal("Failed to close the logger: ", err)
			}
			log.Debug("Closing the database connection")
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

		// cors
		router.Use(func(c *gin.Context) {
			origin := c.Request.Header.Get("Origin")

			// Define allowed origins
			allowedOrigins := map[string]bool{
				"http://localhost:5173":            true, // bun
				"http://localhost:3000":            true, // nginx
				"http://host.docker.internal:3000": true, // docker
				"http://host.docker.internal:5173": true, // docker
			}

			// Check if origin is allowed and set the header
			if allowedOrigins[origin] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		})

		// init event channel
		log.Info("Creating workers...")
		// all webhook events are sent to this channel
		// set to one so new events block old ones and will be discarded
		eventChan := make(chan models.Event, 1)

		// register routes
		RegisterRoutes(router, eventChan)
		err = router.SetTrustedProxies(nil)
		if err != nil {
			logger.Fatal("Failed to set trusted proxies: ", err)
		}

		// run event loop in background
		go eventHandler(ctx, eventChan)

		// init the router
		port := "9999"
		log.Debug("Listening on port", slog.String("port", port))
		log.Info("Started server")
		if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatal(err.Error())
		}
	})
}
