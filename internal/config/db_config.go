package config

import (
	"context"
	"database/sql"
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
			scanDest = append(scanDest, v.Field(i).Addr().Interface())
		}
	}

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
			values = append(values, v.Field(i).Interface())
		}
	}

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
