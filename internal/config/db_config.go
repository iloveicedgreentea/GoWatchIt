package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var (
	globalConfig *Config
	once         sync.Once
)

type Config struct {
	db *sql.DB
}

// InitConfig initializes the global config instance
func InitConfig(db *sql.DB) error {
	var err error
	once.Do(func() {
		globalConfig = &Config{db: db}
		// Initialize tables if they don't exist
		err = globalConfig.initTables()
	})
	return err
}

func (c *Config) initTables() error {
	tables := []interface{}{
		&models.EZBEQConfig{},
		&models.HomeAssistantConfig{},
		&models.JellyfinConfig{},
		&models.MQTTConfig{},
		&models.HDMISyncConfig{},
	}

	for _, table := range tables {
		if err := c.CreateConfigTable(table); err != nil {
			return fmt.Errorf("failed to create table for %T: %v", table, err)
		}
	}
	return nil
}

func (c *Config) LoadConfig(ctx context.Context, cfg interface{}) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.db == nil {
		return fmt.Errorf("db is nil")
	}
	log := logger.GetLoggerFromContext(ctx)
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()
	tableName := t.Name()

	var columns []string
	var scanDest []interface{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)

			// Special handling for slices - use string placeholder
			if v.Field(i).Kind() == reflect.Slice {
				var jsonStr string
				scanDest = append(scanDest, &jsonStr)
				// Store the field index and scan destination for later processing
				if jsonStr != "" {
					var sliceValue reflect.Value
					slicePtr := reflect.New(v.Field(i).Type())
					if err := json.Unmarshal([]byte(jsonStr), slicePtr.Interface()); err == nil {
						sliceValue = slicePtr.Elem()
						v.Field(i).Set(sliceValue)
					}
				}
			} else {
				scanDest = append(scanDest, v.Field(i).Addr().Interface())
			}
		}
	}
	// yes this is theoretically vulnerable to SQL injection, but I control table names and its just not a major risk
	// its more complicated to use prepared statements here because I dynamically get columns from struct tags
	// #nosec
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", strings.Join(columns, ", "), tableName)
	err := c.db.QueryRowContext(ctx, query).Scan(scanDest...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No configuration found in table", "table", tableName)
			return nil
		}
		return fmt.Errorf("failed to scan row: %v", err)
	}

	log.Debug("Loaded configuration", "table", tableName)
	return nil
}

func (c *Config) SaveConfig(cfg interface{}) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.db == nil {
		return fmt.Errorf("db is nil")
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()
	tableName := t.Name()

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, "?")

			// Special handling for slices - convert to JSON string
			if v.Field(i).Kind() == reflect.Slice {
				jsonBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					return fmt.Errorf("failed to marshal slice field %s: %v", field.Name, err)
				}
				values = append(values, string(jsonBytes))
			} else {
				values = append(values, v.Field(i).Interface())
			}
		}
	}

	// yes this is theoretically vulnerable to SQL injection, but I control table names and its just not a major risk
	// #nosec
	query := fmt.Sprintf("INSERT OR REPLACE INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err := c.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	return nil
}

// CreateConfigTable creates a table for the given config struct if it doesn't exist
func (c *Config) CreateConfigTable(cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()
	tableName := t.Name()

	var columns []string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columns = append(columns, fmt.Sprintf("%s %s", dbTag, getSQLType(field.Type)))
		}
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))

	_, err := c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func getSQLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.String:
		return "TEXT"
	default:
		return "TEXT"
	}
}

// Example usage functions
func (c *Config) GetEzbeqConfig() (*models.EZBEQConfig, error) {
	cfg := &models.EZBEQConfig{}
	err := c.LoadConfig(context.Background(), cfg)
	return cfg, err
}

func (c *Config) SaveEzbeqConfig(cfg *models.EZBEQConfig) error {
	return c.SaveConfig(cfg)
}

// Implement similar methods for other config types

// GetConfig is a helper function to get the global config instance
func GetConfig() *Config {
	return globalConfig
}

// ResetConfig resets the global config instance (useful for testing)
func ResetConfig() {
	globalConfig = nil
	once = sync.Once{}
}
