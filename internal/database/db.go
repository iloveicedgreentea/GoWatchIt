package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// GetDB returns a connection to the database
func GetDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
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
