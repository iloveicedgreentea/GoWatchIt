package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/iloveicedgreentea/go-plex/models"
)

// migrateTable will create or update a given table based on a struct
func migrateTable(db *sql.DB, structType interface{}) error {
	tableName := getTableName(structType)

	// Check if table exists
	exists, err := tableExists(db, tableName)
	if err != nil {
		return err
	}

	if !exists {
		// If table doesn't exist, create it
		createSQL := generateTableSQL(structType)
		_, err := db.Exec(createSQL)
		if err != nil {
			return fmt.Errorf("failed to create table %s: %v", tableName, err)
		}
	} else {
		// If table exists, update schema
		if err := updateTableSchema(db, structType); err != nil {
			return err
		}
	}

	// Create or update indices
	indexSQL := generateIndexSQL(structType)
	_, err = db.Exec(indexSQL)
	if err != nil {
		return fmt.Errorf("failed to create/update indices for %s: %v", tableName, err)
	}

	return nil
}

// tableExists checks if a table exists in the database
func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// updateTableSchema will add new columns to an existing table
func updateTableSchema(db *sql.DB, structType interface{}) error {
	tableName := getTableName(structType)
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Get existing columns
	existingColumns, err := getExistingColumns(db, tableName)
	if err != nil {
		return err
	}

	// Compare and add new columns
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		// Skip fields without db tag
		if dbTag == "" {
			continue
		}

		if _, exists := existingColumns[dbTag]; !exists {
			// Add new column
			alterSQL := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
				tableName, dbTag, getSQLType(field.Type))
			_, err := db.Exec(alterSQL)
			if err != nil {
				return fmt.Errorf("failed to add column %s to %s: %v", dbTag, tableName, err)
			}
		}
	}

	return nil
}

// getExistingColumns returns a map of existing columns in a table
func getExistingColumns(db *sql.DB, tableName string) (map[string]bool, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err) // TODO: use logger
		}
	}()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid, notnull, pk int
		var name, type_ string
		var dflt_value sql.NullString
		err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk)
		if err != nil {
			return nil, err
		}
		columns[name] = true
	}

	return columns, nil
}

// getTableName returns the table name for a given struct
func getTableName(structType interface{}) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

// Generates a CREATE TABLE SQL statement for the given struct
func generateTableSQL(structType interface{}) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var columns []string
	var primaryKey string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		columnDef := fmt.Sprintf("%s %s", dbTag, getSQLType(field.Type))

		if dbTag == "id" {
			primaryKey = fmt.Sprintf("PRIMARY KEY (%s)", dbTag)
		} else if field.Tag.Get("unique") == "true" {
			columnDef += " UNIQUE"
		}

		columns = append(columns, columnDef)
	}

	if primaryKey != "" {
		columns = append(columns, primaryKey)
	}

	tableName := strings.ToLower(t.Name())
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);",
		tableName, strings.Join(columns, ",\n\t"))

	return createTableSQL
}

// Generates SQL statements for creating indices on the given struct
func generateIndexSQL(structType interface{}) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var indices []string
	tableName := strings.ToLower(t.Name())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("index") == "true" {
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				continue
			}
			indexName := fmt.Sprintf("idx_%s_%s", tableName, dbTag)
			indexSQL := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(%s);",
				indexName, tableName, dbTag)
			indices = append(indices, indexSQL)
		}
	}

	return strings.Join(indices, "\n")
}

func getSQLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Bool:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}

func getDbModels() []interface{} {
	return []interface{}{
		&models.EZBEQConfig{},
		&models.HomeAssistantConfig{},
		&models.JellyfinConfig{},
		&models.MainConfig{},
		&models.MQTTConfig{},
		&models.PlexConfig{},
		&models.HDMISyncConfig{},
	}
}
