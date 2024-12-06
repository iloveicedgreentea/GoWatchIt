package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// GetDB returns a connection to the database
func GetDB(path string) (*sql.DB, error) {
	// Check if the file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return nil, err
		}

		// Create an empty file
		// TODO: is path ever user supplied? potential directory traversal
		file, err := os.Create(path) // #nosec
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		err = db.Close()
		return nil, err
	}

	return db, nil
}

// TODO: implement "github.com/jmoiron/sqlx" and parameterize the database connection

// RunMigrations runs the necessary migrations for the database
func RunMigrations(db *sql.DB) error {
	structs := getDbModels()
	for _, model := range structs {
		if err := migrateTable(db, model); err != nil {
			return fmt.Errorf("failed to migrate table for %T: %v", model, err)
		}
	}
	return nil
}
