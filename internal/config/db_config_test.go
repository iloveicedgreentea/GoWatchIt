package config

import (
	"database/sql"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "testdb.db")
	require.NoError(t, err)
	return db
}

func createTestTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE EZBEQConfig (
			id INTEGER PRIMARY KEY,
			adjust_master_volume_with_profile BOOLEAN,
			denon_ip TEXT,
			denon_port TEXT,
			dry_run BOOLEAN,
			enabled BOOLEAN,
			enable_tv_beq BOOLEAN,
			notify_endpoint_name TEXT,
			notify_on_load BOOLEAN,
			port TEXT,
			preferred_author TEXT,
			slots TEXT,
			stop_plex_if_mismatch BOOLEAN,
			url TEXT,
			use_avr_codec_search BOOLEAN,
			avr_brand TEXT,
			avr_url TEXT
		);
		
		CREATE TABLE HomeAssistantConfig (
			id INTEGER PRIMARY KEY,
			enabled BOOLEAN,
			pause_script_name TEXT,
			play_script_name TEXT,
			port TEXT,
			remote_entity_name TEXT,
			stop_script_name TEXT,
			token TEXT,
			trigger_aspect_ratio_change_on_event BOOLEAN,
			trigger_avr_master_volume_change_on_event BOOLEAN,
			trigger_lights_on_event BOOLEAN,
			url TEXT
		);
	`)
	require.NoError(t, err)
}

func TestEZBEQConfig(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	createTestTables(t, db)

	cfg, err := NewConfig(db)
	require.NoError(t, err)

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
		Slots:                         `{"slot1": "value1", "slot2": "value2"}`,
		StopPlexIfMismatch:            false,
		URL:                           "http://ezbeq.example.com",
		UseAVRCodecSearch:             true,
		AVRBrand:                      "TestBrand",
		AVRURL:                        "http://avr.example.com",
	}

	err = cfg.SaveConfig(ezbeqConfig)
	assert.NoError(t, err)

	// Test loading EZBEQ config
	loadedConfig, err := cfg.GetEzbeqConfig()
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
	defer db.Close()
	createTestTables(t, db)

	cfg, err := NewConfig(db)
	require.NoError(t, err)

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

	err = cfg.SaveConfig(haConfig)
	assert.NoError(t, err)

	// Test loading HomeAssistant config
	loadedConfig := &models.HomeAssistantConfig{}
	err = cfg.LoadConfig(loadedConfig)
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

// Add similar tests for other config models (JellyfinConfig, MainConfig, MQTTConfig, PlexConfig, HDMISyncConfig)
