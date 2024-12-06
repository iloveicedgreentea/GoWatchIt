package config

import (
	"context"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

// HDMI
func IsHDMISyncEnabled() bool {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return false
	}
	return config.Enabled
}

func IsSignalSourceTime() bool {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return false
	}
	return config.Source == "time"
}

// BEQ
func IsBeqEnabled() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.Enabled
}

func IsBeqTVEnabled() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.EnableTVBEQ
}

func IsBeqNotifyOnLoadEnabled() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.NotifyOnLoad
}

func IsBeqDryRun() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.DryRun
}

// TODO: implement a way to automatically grab configs and add to UI
// TODO: implement in UI
func IsBeqLooseEditionMatching() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.LooseEditionMatching
}

// TODO: implement in ui
func IsBeqSkipEditionMatching() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.SkipEditionMatching
}

// MQTT
func IsMQTTEnabled() bool {
	var config models.MQTTConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load MQTT config", "error", err)
		return false
	}
	return config.Enabled
}

// Jellyfin
func IsJellyfinEnabled() bool {
	var config models.JellyfinConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Jellyfin config", "error", err)
		return false
	}
	return config.Enabled
}

// Home Assistant
func IsHomeAssistantEnabled() bool {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return false
	}
	return config.Enabled
}

func IsHomeAssistantTriggerAVRMasterVolumeChangeOnEvent() bool {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return false
	}
	return config.TriggerAVRMasterVolumeChangeOnEvent && config.Enabled
}

func IsHomeAssistantTriggerLightsOnEvent() bool {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return false
	}
	return config.TriggerLightsOnEvent && config.Enabled
}

func IsJellyfinSkipTMDB() bool {
	var config models.JellyfinConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Jellyfin config", "error", err)
		return false
	}
	return config.SkipTMDB
}
