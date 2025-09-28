package config

import (
	"database/sql"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/pkg/database"
	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
)

func getDbModels() []interface{} {
	return []interface{}{
		&configmodels.EZBEQConfig{},
		&configmodels.HomeAssistantConfig{},
		&configmodels.PlayerConfig{},
		&configmodels.MainConfig{},
		&configmodels.HDMISyncConfig{},
	}
}

// TODO: implement "github.com/jmoiron/sqlx" and parameterize the database connection

// RunMigrations runs the necessary migrations for the database
func RunMigrations(db *sql.DB) error {
	structs := getDbModels()
	for _, model := range structs {
		// create tables if they dont exist and add new columns if they exist
		if err := database.MigrateTable(db, model); err != nil {
			return fmt.Errorf("failed to migrate table for %T: %v", model, err)
		}
	}
	return nil
}
