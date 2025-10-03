package config

import (
	"database/sql"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/pkg/database"
	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

func getDbModels() []interface{} {
	return []interface{}{
		&configmodels.EZBEQConfig{},
		&configmodels.HomeAssistantConfig{},
		&configmodels.PlayerConfig{},
		&configmodels.HDMISyncConfig{},
	}
}

// GetConfigTypeRegistry returns a map of config type names to their struct instances
// This is used to dynamically unmarshal config data without maintaining duplicate lists
func GetConfigTypeRegistry() map[string]interface{} {
	return map[string]interface{}{
		"ezbeq":         &configmodels.EZBEQConfig{},
		"homeassistant": &configmodels.HomeAssistantConfig{},
		"player":        &configmodels.PlayerConfig{},
		"hdmisync":      &configmodels.HDMISyncConfig{},
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

	// Migrate domain models
	if err := database.MigrateTable(db, &beq.BEQAuthor{}); err != nil {
		return fmt.Errorf("failed to migrate BEQAuthor table: %v", err)
	}

	return nil
}
