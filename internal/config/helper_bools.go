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

// Plex
func IsPlexEnabled() bool {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return false
	}
	return config.Enabled
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

func IsBeqNotifyOnUnLoadEnabled() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.NotifyOnUnLoad
}

func IsBeqDryRun() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.DryRun
}

func IsBeqLooseEditionMatching() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.LooseEditionMatching
}

func IsBeqSkipEditionMatching() bool {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return false
	}
	return config.SkipEditionMatching
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

func IsJellyfinSkipTMDB() bool {
	var config models.JellyfinConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Jellyfin config", "error", err)
		return false
	}
	return config.SkipTMDB
}
