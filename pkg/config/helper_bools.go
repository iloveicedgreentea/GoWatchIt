package config

import (
	"context"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
)

var gLog = logger.GetLogger()

// HDMI
func IsHDMISyncEnabled() bool {
	var config HDMISyncConfig
	// TODO: have a worker to refresh config for each enabled service and modify these functions to use a global object instead
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync config", zap.Error(err))
		return false
	}
	return config.Enabled
}

func IsSignalSourceTime() bool {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync config", zap.Error(err))
		return false
	}
	return config.Source == "time"
}

// Plex
func IsPlexEnabled() bool {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex config", zap.Error(err))
		return false
	}
	return config.Enabled
}

// BEQ
func IsBeqEnabled() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.Enabled
}

func IsBeqTVEnabled() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.EnableTVBEQ
}

func IsBeqNotifyOnLoadEnabled() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.NotifyOnLoad
}

func IsBeqNotifyOnUnLoadEnabled() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.NotifyOnUnLoad
}

func IsBeqDryRun() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.DryRun
}

func IsBeqLooseEditionMatching() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.LooseEditionMatching
}

func IsBeqSkipEditionMatching() bool {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ config", zap.Error(err))
		return false
	}
	return config.SkipEditionMatching
}

// Jellyfin
func IsJellyfinEnabled() bool {
	var config JellyfinConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Jellyfin config", zap.Error(err))
		return false
	}
	return config.Enabled
}

// Home Assistant
func IsHomeAssistantEnabled() bool {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant config", zap.Error(err))
		return false
	}
	return config.Enabled
}

func IsJellyfinSkipTMDB() bool {
	var config JellyfinConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Jellyfin config", zap.Error(err))
		return false
	}
	return config.SkipTMDB
}
