package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	db, err := GetDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to get in-memory database: %v", err)
	}
	defer func() {
		require.NoError(t, db.Close())
	}()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping in-memory database: %v", err)
	}
}
