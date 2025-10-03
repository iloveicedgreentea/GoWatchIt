package config

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	l "log"

	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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
		err = RunMigrations(db)
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
	ezbeqConfig := &configmodels.EZBEQConfig{
		AdjustMasterVolumeWithProfile: true,
		DenonIP:                       "192.168.1.100",
		DenonPort:                     "8080",
		DryRun:                        false,
		Enabled:                       true,
		EnableTVBEQ:                   true,
		NotifyOnLoad:                  true,
		NotifyOnUnload:                true,
		Port:                          "8081",
		PreferredAuthors:              []string{"Author1", "Author2"},
		Slots:                         []int32{1, 2, 3},
		StopPlexIfMismatch:            false,
		Url:                           "http://ezbeq.example.com",
		UseAVRCodecSearch:             true,
		AvrBrand:                      "TestBrand",
		AvrURL:                        "http://avr.example.com",
	}

	err := GetConfig().SaveConfig(ezbeqConfig)
	assert.NoError(t, err)

	// Test loading EZBEQ config
	ctx := context.Background()
	assert.True(t, IsBeqEnabled(ctx))
	assert.True(t, IsBeqTVEnabled(ctx))
	assert.True(t, IsBeqNotifyOnLoadEnabled(ctx))
	assert.True(t, IsBeqNotifyOnUnLoadEnabled(ctx))
	assert.False(t, IsBeqDryRun(ctx))

	// Test loading full config
	var loadedConfig configmodels.EZBEQConfig
	err = GetConfig().LoadConfig(context.Background(), &loadedConfig)
	assert.NoError(t, err)
	assert.Equal(t, ezbeqConfig.AdjustMasterVolumeWithProfile, loadedConfig.AdjustMasterVolumeWithProfile)
	assert.Equal(t, ezbeqConfig.DenonIP, loadedConfig.DenonIP)
	assert.Equal(t, ezbeqConfig.DenonPort, loadedConfig.DenonPort)
	assert.Equal(t, ezbeqConfig.DryRun, loadedConfig.DryRun)
	assert.Equal(t, ezbeqConfig.Enabled, loadedConfig.Enabled)
	assert.Equal(t, ezbeqConfig.EnableTVBEQ, loadedConfig.EnableTVBEQ)
	assert.Equal(t, ezbeqConfig.NotifyOnLoad, loadedConfig.NotifyOnLoad)
	assert.Equal(t, ezbeqConfig.NotifyOnUnload, loadedConfig.NotifyOnUnload)
	assert.Equal(t, ezbeqConfig.Port, loadedConfig.Port)
	assert.Equal(t, ezbeqConfig.PreferredAuthors, loadedConfig.PreferredAuthors)
	assert.Equal(t, ezbeqConfig.Slots, loadedConfig.Slots, "Slots should be equal to %v but got %v", ezbeqConfig.Slots, loadedConfig.Slots)
	assert.Equal(t, ezbeqConfig.StopPlexIfMismatch, loadedConfig.StopPlexIfMismatch)
	assert.Equal(t, ezbeqConfig.Url, loadedConfig.Url)
	assert.Equal(t, ezbeqConfig.UseAVRCodecSearch, loadedConfig.UseAVRCodecSearch)
	assert.Equal(t, ezbeqConfig.AvrBrand, loadedConfig.AvrBrand)
	assert.Equal(t, ezbeqConfig.AvrURL, loadedConfig.AvrURL)
}

func TestHomeAssistantConfig(t *testing.T) {
	// Test saving HomeAssistant config
	haConfig := &configmodels.HomeAssistantConfig{
		Enabled:                         true,
		Port:                            "8123",
		RemoteEntityName:                "remote.living_room",
		Token:                           "test_token",
		TriggerAspectRatioChangeOnEvent: true,
		Url:                             "homeassistant.local",
		Scheme:                          configmodels.HomeAssistantConfigSchemeHttp,
		NotifyEndpointName:              "test_endpoint",
		NotifyDisplayTime:               5000,
	}

	err := GetConfig().SaveConfig(haConfig)
	assert.NoError(t, err)

	// Test helper functions
	ctx := context.Background()
	assert.True(t, IsHomeAssistantEnabled(ctx))

	// Test loading full config
	var loadedConfig configmodels.HomeAssistantConfig
	err = GetConfig().LoadConfig(context.Background(), &loadedConfig)
	assert.NoError(t, err)
	assert.Equal(t, haConfig.Enabled, loadedConfig.Enabled)
	assert.Equal(t, haConfig.Port, loadedConfig.Port)
	assert.Equal(t, haConfig.RemoteEntityName, loadedConfig.RemoteEntityName)
	assert.Equal(t, haConfig.Token, loadedConfig.Token)
	assert.Equal(t, haConfig.TriggerAspectRatioChangeOnEvent, loadedConfig.TriggerAspectRatioChangeOnEvent)
	assert.Equal(t, haConfig.Url, loadedConfig.Url)
	assert.Equal(t, haConfig.Scheme, loadedConfig.Scheme)
	assert.Equal(t, haConfig.NotifyEndpointName, loadedConfig.NotifyEndpointName)
	assert.Equal(t, haConfig.NotifyDisplayTime, loadedConfig.NotifyDisplayTime)
}

// Add similar tests for other config types (JellyfinConfig, MQTTConfig, HDMISyncConfig)
