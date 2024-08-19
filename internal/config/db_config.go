package config

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iloveicedgreentea/go-plex/models"
)

type Config struct {
	db *sql.DB
}

func NewConfig(db *sql.DB) (*Config, error) {
	return &Config{db: db}, nil
}

// TODO: verify its loading correctly

// LoadConfig loads a configuration struct from the database
func (c *Config) LoadConfig(cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	tableName := t.Name()
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := c.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query table %s: %v", tableName, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(sql.RawBytes)
		}

		err = rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		for i, colName := range columns {
			field := v.FieldByName(colName)
			if field.IsValid() && field.CanSet() {
				err = setField(field, string(*values[i].(*sql.RawBytes)))
				if err != nil {
					return fmt.Errorf("failed to set field %s: %v", colName, err)
				}
			}
		}

		// We've loaded one complete row, so we can break here
		// If you expect multiple rows, you might need to handle this differently
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating over rows: %v", rows.Err())
	}

	return nil
}

// SaveConfig saves a configuration struct to the database
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

// Helper function to set field value from string
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	// Add more types as needed
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}

// Example usage functions
func (c *Config) GetEzbeqConfig() (*models.EZBEQConfig, error) {
	cfg := &models.EZBEQConfig{}
	err := c.LoadConfig(cfg)
	return cfg, err
}

func (c *Config) SaveEzbeqConfig(cfg *models.EZBEQConfig) error {
	return c.SaveConfig(cfg)
}

// Implement similar methods for other config types
