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
	once.Do(func() {
		globalConfig = &Config{db: db}
	})
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
	// Map to store slice field indices and their corresponding scan destinations
	sliceFields := make(map[int]*string)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)

			// Special handling for slices because sqlite doesn't support them
			if v.Field(i).Kind() == reflect.Slice {
				log.Debug("Loading slice field", "field", field.Name)
				var jsonStr string
				// Store the pointer to the JSON string in the map
				sliceFields[i] = &jsonStr
				// data will be written to the pointer
				scanDest = append(scanDest, &jsonStr)
			} else {
				scanDest = append(scanDest, v.Field(i).Addr().Interface())
			}
		}
	}
	// #nosec - yes this is theoretically vulnerable to SQL injection, but I control table names and its just not a major risk
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", strings.Join(columns, ", "), tableName)
	err := c.db.QueryRowContext(ctx, query).Scan(scanDest...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No configuration found in table", "table", tableName)
			return nil
		}
		return fmt.Errorf("failed to scan row: %v", err)
	}

	// Process slice fields after successful scan
	for fieldIndex, jsonStrPtr := range sliceFields {
		if jsonStrPtr != nil && *jsonStrPtr != "" {
			sliceField := v.Field(fieldIndex)
			// Create a new slice value of the correct type
			newSlice := reflect.New(sliceField.Type())
			if err := json.Unmarshal([]byte(*jsonStrPtr), newSlice.Interface()); err != nil {
				return fmt.Errorf("failed to unmarshal slice for field %s: %v", t.Field(fieldIndex).Name, err)
			}
			// Set the slice value in the struct
			sliceField.Set(newSlice.Elem())
		}
	}

	log.Debug("Loaded configuration", "table", tableName)
	return nil
}

func (c *Config) SaveConfig(cfg interface{}) error {
	// TODO: add safety checks for URLs make sure they are reachable
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

// TODO: Implement similar methods for other config types
// GetConfig is a helper function to get the global config instance
func GetConfig() *Config {
	return globalConfig
}

// ResetConfig resets the global config instance (useful for testing)
func ResetConfig() {
	globalConfig = nil
	once = sync.Once{}
}
