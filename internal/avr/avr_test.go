package avr

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	l "log"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

var (
	db     *sql.DB
	dbOnce sync.Once
	client AVRClient
)

// TestMain sets up the database and runs the tests
func TestMain(m *testing.M) {
	var code int
	dbOnce.Do(func() {
		// Setup code before tests
		var err error

		// Open SQLite database connection
		db, err = database.GetDB(":memory:")
		if err != nil {
			l.Fatalf("Failed to open database: %v", err)
		}

		// run migrations
		err = database.RunMigrations(db)
		if err != nil {
			l.Fatalf("Failed to run migrations: %v", err)
		}

		// Initialize the config with the database
		err = config.InitConfig(db)
		if err != nil {
			l.Fatalf("Failed to initialize config: %v", err)
		}

		cf := config.GetConfig()

		// populate test data
		beqCfg := models.EZBEQConfig{
			Enabled:           true,
			UseAVRCodecSearch: true,
			AVRURL:            "192.168.88.40",
			AVRBrand:          "denon",
		}
		err = cf.SaveConfig(&beqCfg)
		if err != nil {
			l.Fatalf("Failed to save ezbeq config: %v", err)
		}
		client = GetAVRClient()
		// Run the tests
		code = m.Run()

		// Cleanup code after tests
		err = db.Close()
		if err != nil {
			l.Printf("Error closing database: %v", err)
		}
	})
	// Exit with the test result code
	os.Exit(code)
}

func TestAvrGetAudioMode(t *testing.T) {
	t.SkipNow()
	mode, err := client.GetCodec()
	assert.NoError(t, err)
	t.Log(mode)
	assert.NotEmpty(t, mode)
}

// TODO: test failure to init avr if brand or url empty
