package config

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	l "log"

	"github.com/iloveicedgreentea/go-plex/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/iloveicedgreentea/go-plex/internal/database"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func TestMain(m *testing.M) {
	var code int
	dbOnce.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			l.Fatalf("Failed to open database: %v", err)
		}

		// run migrations
		err = database.RunMigrations(db)
		if err != nil {
			l.Fatalf("Failed to run migrations: %v", err)
		}

		// Initialize the config with the database
		err = InitConfig(db)
		if err != nil {
			l.Fatalf("Failed to initialize config: %v", err)
		}
		// Run the tests
		code = m.Run()

		// Cleanup code after tests
		err = db.Close()
		if err != nil {
			l.Printf("Error closing database: %v", err)
		}
	})

	os.Exit(code)
}

func TestEZBEQConfig(t *testing.T) {
	// Test saving EZBEQ config
	ezbeqConfig := &models.EZBEQConfig{
		AdjustMasterVolumeWithProfile: true,
		DenonIP:                       "192.168.1.100",
		DenonPort:                     "8080",
		DryRun:                        false,
		Enabled:                       true,
		EnableTVBEQ:                   true,
		NotifyOnLoad:                  true,
		Port:                          "8081",
		PreferredAuthor:               "TestAuthor",
		Slots:                         []int{1, 2, 3},
		StopPlexIfMismatch:            false,
		URL:                           "http://ezbeq.example.com",
		UseAVRCodecSearch:             true,
		AVRBrand:                      "TestBrand",
		AVRURL:                        "http://avr.example.com",
	}

	err := GetConfig().SaveConfig(ezbeqConfig)
	assert.NoError(t, err)

	// Test loading EZBEQ config
	assert.True(t, IsBeqEnabled())
	assert.True(t, IsBeqTVEnabled())
	assert.True(t, IsBeqNotifyOnLoadEnabled())
	assert.False(t, IsBeqDryRun())

	// Test loading full config
	var loadedConfig models.EZBEQConfig
	err = GetConfig().LoadConfig(context.Background(), &loadedConfig)
	assert.NoError(t, err)
	assert.Equal(t, ezbeqConfig.AdjustMasterVolumeWithProfile, loadedConfig.AdjustMasterVolumeWithProfile)
	assert.Equal(t, ezbeqConfig.DenonIP, loadedConfig.DenonIP)
	assert.Equal(t, ezbeqConfig.DenonPort, loadedConfig.DenonPort)
	assert.Equal(t, ezbeqConfig.DryRun, loadedConfig.DryRun)
	assert.Equal(t, ezbeqConfig.Enabled, loadedConfig.Enabled)
	assert.Equal(t, ezbeqConfig.EnableTVBEQ, loadedConfig.EnableTVBEQ)
	assert.Equal(t, ezbeqConfig.NotifyOnLoad, loadedConfig.NotifyOnLoad)
	assert.Equal(t, ezbeqConfig.Port, loadedConfig.Port)
	assert.Equal(t, ezbeqConfig.PreferredAuthor, loadedConfig.PreferredAuthor)
	assert.Equal(t, ezbeqConfig.Slots, loadedConfig.Slots, "Slots should be equal to %v but got %v", ezbeqConfig.Slots, loadedConfig.Slots)
	assert.Equal(t, ezbeqConfig.StopPlexIfMismatch, loadedConfig.StopPlexIfMismatch)
	assert.Equal(t, ezbeqConfig.URL, loadedConfig.URL)
	assert.Equal(t, ezbeqConfig.UseAVRCodecSearch, loadedConfig.UseAVRCodecSearch)
	assert.Equal(t, ezbeqConfig.AVRBrand, loadedConfig.AVRBrand)
	assert.Equal(t, ezbeqConfig.AVRURL, loadedConfig.AVRURL)
}

func TestHomeAssistantConfig(t *testing.T) {
	// Test saving HomeAssistant config
	haConfig := &models.HomeAssistantConfig{
		Enabled:                         true,
		Port:                            "8123",
		RemoteEntityName:                "remote.living_room",
		Token:                           "test_token",
		TriggerAspectRatioChangeOnEvent: true,
		URL:                             "homeassistant.local",
		Scheme:                          "http",
		MediaPlayerEntityName:           "media_player.test",
		NotifyEndpointName:              "test_endpoint",
	}

	err := GetConfig().SaveConfig(haConfig)
	assert.NoError(t, err)

	// Test helper functions
	assert.True(t, IsHomeAssistantEnabled())

	// Test loading full config
	var loadedConfig models.HomeAssistantConfig
	err = GetConfig().LoadConfig(context.Background(), &loadedConfig)
	assert.NoError(t, err)
	assert.Equal(t, haConfig.Enabled, loadedConfig.Enabled)
	assert.Equal(t, haConfig.Port, loadedConfig.Port)
	assert.Equal(t, haConfig.RemoteEntityName, loadedConfig.RemoteEntityName)
	assert.Equal(t, haConfig.Token, loadedConfig.Token)
	assert.Equal(t, haConfig.TriggerAspectRatioChangeOnEvent, loadedConfig.TriggerAspectRatioChangeOnEvent)
	assert.Equal(t, haConfig.URL, loadedConfig.URL)
	assert.Equal(t, haConfig.Scheme, loadedConfig.Scheme)
	assert.Equal(t, haConfig.MediaPlayerEntityName, loadedConfig.MediaPlayerEntityName)
	assert.Equal(t, haConfig.NotifyEndpointName, loadedConfig.NotifyEndpointName)
}

// Add similar tests for other config types (JellyfinConfig, MQTTConfig, HDMISyncConfig)
