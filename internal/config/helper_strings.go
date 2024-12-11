package config

import (
	"context"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

// Main
func GetMainListenPort() string {
	var config models.MainConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Main config", "error", err)
		return ""
	}
	return config.ListenPort
}

// HDMI
func GetHDMISyncSource() string {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return ""
	}
	return config.Source
}

func GetHDMISyncEnvyName() string {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return ""
	}
	return config.Envy
}

func GetHDMISyncSeconds() string {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return ""
	}
	return config.Time
}

func GetHDMISyncPlayerIP() string {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return ""
	}
	return config.PlayerIP
}

func GetHDMISyncMachineIdentifier() string {
	var config models.HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HDMISync config", "error", err)
		return ""
	}
	return config.PlayerMachineIdentifier
}

// Home Assistant

func GetHomeAssistantUrl() string {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return ""
	}
	return config.URL
}

func GetHomeAssistantToken() string {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return ""
	}
	return config.Token
}

func GetHomeAssistantPort() string {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return ""
	}
	return config.Port
}

func GetHomeAssistantRemoteEntityName() string {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load HomeAssistant config", "error", err)
		return ""
	}
	return config.RemoteEntityName
}

func GetHomeAssistantNotifyEndpointName() string {
	var config models.HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return ""
	}
	return config.NotifyEndpointName
}

// EZBeq
func GetEZBeqUrl() string {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return ""
	}
	return config.URL
}

func GetEZBeqScheme() string {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return ""
	}
	if config.Scheme == "" {
		return "http"
	}
	return config.Scheme
}

func GetEZBeqPort() string {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return ""
	}
	return config.Port
}

func GetEZBeqSlots() []int {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return []int{}
	}
	return config.Slots
}

func GetEZBeqPreferredAuthor() string {
	var config models.EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load EZBEQ config", "error", err)
		return ""
	}
	return config.PreferredAuthor
}

// Plex

func GetPlexUrl() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.URL
}

func GetPlexToken() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.Token
}

func GetPlexPort() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.Port
}

func GetPlexScheme() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.Scheme
}

func GetPlexDeviceUUIDFilter() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.DeviceUUIDFilter
}

func GetPlexOwnerNameFilter() string {
	var config models.PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		logger.Error("Failed to load Plex config", "error", err)
		return ""
	}
	return config.OwnerNameFilter
}
