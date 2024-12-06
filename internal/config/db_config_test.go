package config

import (
	"context"
	"database/sql"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}

func TestEZBEQConfig(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	err := InitConfig(db)
	require.NoError(t, err)
	defer ResetConfig() // Reset the global config after the test

	// Test saving EZBEQ config
	ezbeqConfig := &models.EZBEQConfig{
		AdjustMasterVolumeWithProfile: true,
		DenonIP:                       "192.168.1.100",
		DenonPort:                     "8080",
		DryRun:                        false,
		Enabled:                       true,
		EnableTVBEQ:                   true,
		NotifyEndpointName:            "test_endpoint",
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

	err = GetConfig().SaveConfig(ezbeqConfig)
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
	assert.Equal(t, ezbeqConfig.NotifyEndpointName, loadedConfig.NotifyEndpointName)
	assert.Equal(t, ezbeqConfig.NotifyOnLoad, loadedConfig.NotifyOnLoad)
	assert.Equal(t, ezbeqConfig.Port, loadedConfig.Port)
	assert.Equal(t, ezbeqConfig.PreferredAuthor, loadedConfig.PreferredAuthor)
	assert.Equal(t, ezbeqConfig.Slots, loadedConfig.Slots)
	assert.Equal(t, ezbeqConfig.StopPlexIfMismatch, loadedConfig.StopPlexIfMismatch)
	assert.Equal(t, ezbeqConfig.URL, loadedConfig.URL)
	assert.Equal(t, ezbeqConfig.UseAVRCodecSearch, loadedConfig.UseAVRCodecSearch)
	assert.Equal(t, ezbeqConfig.AVRBrand, loadedConfig.AVRBrand)
	assert.Equal(t, ezbeqConfig.AVRURL, loadedConfig.AVRURL)
}

func TestHomeAssistantConfig(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	err := InitConfig(db)
	require.NoError(t, err)
	defer ResetConfig() // Reset the global config after the test

	// Test saving HomeAssistant config
	haConfig := &models.HomeAssistantConfig{
		Enabled:                             true,
		PauseScriptName:                     "pause_script",
		PlayScriptName:                      "play_script",
		Port:                                "8123",
		RemoteEntityName:                    "remote.living_room",
		StopScriptName:                      "stop_script",
		Token:                               "test_token",
		TriggerAspectRatioChangeOnEvent:     true,
		TriggerAVRMasterVolumeChangeOnEvent: true,
		TriggerLightsOnEvent:                true,
		URL:                                 "http://homeassistant.local",
	}

	err = GetConfig().SaveConfig(haConfig)
	assert.NoError(t, err)

	// Test helper functions
	assert.True(t, IsHomeAssistantEnabled())
	assert.True(t, IsHomeAssistantTriggerAVRMasterVolumeChangeOnEvent())
	assert.True(t, IsHomeAssistantTriggerLightsOnEvent())

	// Test loading full config
	var loadedConfig models.HomeAssistantConfig
	err = GetConfig().LoadConfig(context.Background(), &loadedConfig)
	assert.NoError(t, err)
	assert.Equal(t, haConfig.Enabled, loadedConfig.Enabled)
	assert.Equal(t, haConfig.PauseScriptName, loadedConfig.PauseScriptName)
	assert.Equal(t, haConfig.PlayScriptName, loadedConfig.PlayScriptName)
	assert.Equal(t, haConfig.Port, loadedConfig.Port)
	assert.Equal(t, haConfig.RemoteEntityName, loadedConfig.RemoteEntityName)
	assert.Equal(t, haConfig.StopScriptName, loadedConfig.StopScriptName)
	assert.Equal(t, haConfig.Token, loadedConfig.Token)
	assert.Equal(t, haConfig.TriggerAspectRatioChangeOnEvent, loadedConfig.TriggerAspectRatioChangeOnEvent)
	assert.Equal(t, haConfig.TriggerAVRMasterVolumeChangeOnEvent, loadedConfig.TriggerAVRMasterVolumeChangeOnEvent)
	assert.Equal(t, haConfig.TriggerLightsOnEvent, loadedConfig.TriggerLightsOnEvent)
	assert.Equal(t, haConfig.URL, loadedConfig.URL)
}

// Add similar tests for other config types (JellyfinConfig, MQTTConfig, HDMISyncConfig)
