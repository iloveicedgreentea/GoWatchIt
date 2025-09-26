package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStruct is a sample struct for testing
type TestStruct struct {
	ID        int     `db:"id"`
	Name      string  `db:"name"`
	Age       int     `db:"age"`
	IsActive  bool    `db:"is_active"`
	Score     float64 `db:"score"`
	ExtraInfo string  `db:"extra_info" index:"true"`
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}

func TestMigrateTable(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	err := migrateTable(db, &TestStruct{})
	assert.NoError(t, err)

	// Check if table exists
	exists, err := tableExists(db, "teststruct")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check if columns are created correctly
	columns, err := getExistingColumns(db, "teststruct")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(columns))
	assert.True(t, columns["id"])
	assert.True(t, columns["name"])
	assert.True(t, columns["age"])
	assert.True(t, columns["is_active"])
	assert.True(t, columns["score"])
	assert.True(t, columns["extra_info"])

	// Check if index is created
	var indexName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='teststruct'").Scan(&indexName)
	assert.NoError(t, err)
	assert.Equal(t, "idx_teststruct_extra_info", indexName)
}

func TestUpdateTableSchema(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	// Create initial table
	_, err := db.Exec("CREATE TABLE teststruct (id INTEGER PRIMARY KEY, name TEXT);")
	require.NoError(t, err)

	// Update schema
	err = updateTableSchema(db, &TestStruct{})
	assert.NoError(t, err)

	// Check if new columns are added
	columns, err := getExistingColumns(db, "teststruct")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(columns))
	assert.True(t, columns["age"])
	assert.True(t, columns["is_active"])
	assert.True(t, columns["score"])
	assert.True(t, columns["extra_info"])
}

func TestGenerateTableSQL(t *testing.T) {
	s := generateTableSQL(&TestStruct{})
	expectedSQL := `CREATE TABLE IF NOT EXISTS teststruct (
	id INTEGER,
	name TEXT,
	age INTEGER,
	is_active BOOLEAN,
	score REAL,
	extra_info TEXT,
	PRIMARY KEY (id)
);`
	assert.Equal(t, expectedSQL, s)
}

func TestGenerateIndexSQL(t *testing.T) {
	s := generateIndexSQL(&TestStruct{})
	expectedSQL := "CREATE INDEX IF NOT EXISTS idx_teststruct_extra_info ON teststruct(extra_info);"
	assert.Equal(t, expectedSQL, s)
}

func TestGetSQLType(t *testing.T) {
	assert.Equal(t, "INTEGER", getSQLType(reflect.TypeOf(0)))
	assert.Equal(t, "REAL", getSQLType(reflect.TypeOf(0.0)))
	assert.Equal(t, "BOOLEAN", getSQLType(reflect.TypeOf(true)))
	assert.Equal(t, "TEXT", getSQLType(reflect.TypeOf("")))
}

func TestRunMigrations(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	// Set up our test models
	testModels := []interface{}{&TestStruct{}}

	// Modify RunMigrations to use our testModels
	err := RunMigrationsForModels(db, testModels)
	assert.NoError(t, err)

	// Check if table exists
	exists, err := tableExists(db, "teststruct")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check if all columns are created
	columns, err := getExistingColumns(db, "teststruct")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(columns))
}

// Add this function to your main database package
func RunMigrationsForModels(db *sql.DB, models []interface{}) error {
	for _, model := range models {
		if err := migrateTable(db, model); err != nil {
			return fmt.Errorf("failed to migrate table for %T: %v", model, err)
		}
	}
	return nil
}
