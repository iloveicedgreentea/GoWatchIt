package homeassistant

import (
	"database/sql"
	l "log"
	"os"
	"sync"
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db       *sql.DB
	dbOnce   sync.Once
	haClient *HomeAssistantClient
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
		haCfg := models.HomeAssistantConfig{
			Enabled:               true,
			URL:                   "homeassistant.local",
			Port:                  "8123",
			Scheme:                "http",
			MediaPlayerEntityName: "media_player.test",
			Token:                 os.Getenv("HA_TOKEN"),
			NotifyEndpointName:    "notify.mobile_app_iphone",
		}
		err = cf.SaveConfig(&haCfg)
		if err != nil {
			l.Fatalf("Failed to save ezbeq config: %v", err)
		}

		haClient, err = NewClient()
		if err != nil {
			l.Fatalf("Failed to create home assistant client: %v", err)
		}

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

// make sure script can trigger
func TestScriptTrigger(t *testing.T) {
	t.Skip()
	t.Parallel()

	// trigger an empty script to verify client
	err := haClient.TriggerScript("test")
	assert.NoError(t, err)
}

// test sending a real notification
func TestNotification(t *testing.T) {
	t.Parallel()
	// trigger light and switch
	err := haClient.SendNotification("test from gowatchit")
	assert.NoError(t, err)
}

// TODO: test read state, other stuff
func TestReadAttributes(t *testing.T) {
	t.Parallel()
	type testStruct struct {
		entName string
		test    HAAttributeResponse
		entType models.HomeAssistantEntity
	}
	tt := []testStruct{
		{
			entName: "theater_2",
			test:    &models.HAMediaPlayerResponse{},
			entType: models.HomeAssistantEntityMediaPlayer,
		},
	}

	for _, k := range tt {
		attributes, err := haClient.ReadAttributes(k.entName, k.test, k.entType)
		require.NoError(t, err)

		state, err := haClient.ReadState(k.entName, k.test, k.entType)
		require.NoError(t, err)

		assert.NotEmpty(t, state)
		assert.NotEmpty(t, attributes)

	}
}
