package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/database"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/ports/webhooks"
	"go.uber.org/zap"
)

// Register handlers from ports - events
func run(ctx context.Context) error {
	// create logger
	log := logger.GetLoggerFromContext(ctx)
	defer log.Sync()

	log.Info("Starting service")
	debug := strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug"
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		baseDir = "."
	}

	dbDir := fmt.Sprintf("%s/db.sqlite3", baseDir)
	log.Debug("Base directory", zap.String("dbDir", dbDir))

	// Create the database connection
	log.Info("Connecting to the database...")
	db, err := database.GetDB(dbDir)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	if db == nil {
		return errors.New("db is nil")
	}

	// close db when done
	defer func() {
		if err := logger.CleanupLogger(); err != nil {
			log.Error("Failed to close the logger", zap.Error(err))
		}
		log.Debug("Closing the database connection")
		if err := db.Close(); err != nil {
			log.Error("Failed to close the database", zap.Error(err))
		}
	}()

	// create or update tables
	log.Info("Running migrations...")
	err = config.RunMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// init the config manager
	log.Info("Initializing config...")
	err = config.InitConfig(db)
	if err != nil {
		return fmt.Errorf("failed to run init config: %w", err)
	}

	// init beq client
	// TODO: must support being loaded while off in config
	beqClient, err := beq.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create beq client: %w", err)
	}

	// initialize app
	app, err := app.NewApplication(ctx, beqClient)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	// set up routes
	router, err := webhooks.NewRouter(ctx, app)
	if err != nil {
		return err
	}

	// start server
	port := "9999"
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Debug("Listening on port", zap.String("port", port))
	log.Info("Service listening...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	log.Info("Shutting down service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("service failed: %v", err)
	}
}
