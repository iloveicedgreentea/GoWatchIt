package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

// GetDB returns a connection to the database
func GetDB(path string) (*sql.DB, error) {
	log := logger.GetLogger()
	if path != ":memory:" {
		log.Debug("not using in-memory database", slog.String("path", path))
		// Check if the file exists
		_, err := os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to stat database file: %w", err)
			}
			log.Debug("Database file does not exist", slog.String("path", path))
		}
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
			log.Debug("Creating database directory", slog.String("dir", dir))

			// Create an empty file
			// TODO: is path ever user supplied? potential directory traversal
			file, err := os.Create(path) // #nosec
			if err != nil {
				return nil, fmt.Errorf("failed to create database file: %w", err)
			}
			log.Debug("File created", slog.String("path", path))
			err = file.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to close created database file: %w", err)
			}
		}
	}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	log.Debug("Database connection opened", slog.String("path", path), slog.Any("db", db))

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Error("Failed to ping the database", slog.Any("error", err))
		err2 := db.Close()
		return nil, fmt.Errorf("failed to ping the database: %w %w", err, err2)
	}

	var count int
	err = db.QueryRow("SELECT 1").Scan(&count)
	if err != nil {
		logger.Fatal("Failed to run test query: ", err)
	}
	log.Debug("Successfully ran test query")

	return db, nil
}

// TODO: implement "github.com/jmoiron/sqlx" and parameterize the database connection

// RunMigrations runs the necessary migrations for the database
func RunMigrations(db *sql.DB) error {
	structs := getDbModels()
	for _, model := range structs {
		// create tables if they dont exist and add new columns if they exist
		if err := migrateTable(db, model); err != nil {
			return fmt.Errorf("failed to migrate table for %T: %v", model, err)
		}
	}
	return nil
}
