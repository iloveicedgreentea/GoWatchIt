package database

import (
	"testing"
)

func TestGetDB(t *testing.T) {
	db, err := GetDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to get in-memory database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping in-memory database: %v", err)
	}
}

