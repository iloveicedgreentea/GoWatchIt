package config

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

func santizeURL(url string) string {
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	return url
}

// HDMI
func GetHDMISyncSource() string {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync source", zap.Error(err))
		return ""
	}
	return config.Source
}

func GetHDMISyncEnvyName() string {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync name", zap.Error(err))
		return ""
	}
	return config.Envy
}

func GetHDMISyncSeconds() string {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync seconds", zap.Error(err))
		return ""
	}
	return config.Time
}

func GetHDMISyncPlayerIP() string {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync ip", zap.Error(err))
		return ""
	}
	return config.PlayerIP
}

func GetHDMISyncMachineIdentifier() string {
	var config HDMISyncConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HDMISync id", zap.Error(err))
		return ""
	}
	return config.PlayerMachineIdentifier
}

// Home Assistant

func GetHomeAssistantUrl() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant url", zap.Error(err))
		return ""
	}

	config.URL = santizeURL(config.URL)

	if config.URL == "" {
		return "homeassistant.local"
	}

	return config.URL
}

func GetHomeAssistantScheme() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant scheme", zap.Error(err))
		return ""
	}

	if config.Scheme == "" {
		return "http"
	}

	return config.Scheme
}

func GetHomeAssistantToken() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant token", zap.Error(err))
		return ""
	}
	return config.Token
}

func GetHomeAssistantPort() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant port", zap.Error(err))
		return ""
	}

	if config.Port == "" {
		return "8123"
	}

	return config.Port
}

func GetHomeAssistantRemoteEntityName() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant remote entity", zap.Error(err))
		return ""
	}
	return config.RemoteEntityName
}

func GetHomeAssistantNotifyEndpointName() string {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant notify endpoint", zap.Error(err))
		return ""
	}
	return config.NotifyEndpointName
}

func GetHomeAssistantNotifyDisplayTime() int {
	var config HomeAssistantConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load HomeAssistant notify display time", zap.Error(err))
		return 0
	}
	// set default to 15
	if config.NotifyDisplayTime == 0 {
		// Default to 15 seconds (15000 ms)
		return 15 * 1000
	}
	// Convert configured seconds to milliseconds
	return config.NotifyDisplayTime * 1000
}

// EZBeq
func GetEZBeqUrl() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ url", zap.Error(err))
		return ""
	}

	config.URL = santizeURL(config.URL)

	if config.URL == "" {
		return "ezbeq.local"
	}

	return config.URL
}

func GetEZBeqScheme() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ scheme", zap.Error(err))
		return ""
	}

	if config.Scheme == "" {
		return "http"
	}

	return config.Scheme
}

func GetEZBeqPort() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ port", zap.Error(err))
		return ""
	}

	if config.Port == "" {
		return "8080"
	}

	return config.Port
}

func GetEZBeqAvrURL() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ avr url", zap.Error(err))
		return ""
	}
	return config.AVRURL
}

func GetEZBeqAvrBrand() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ avr brand", zap.Error(err))
		return ""
	}
	return config.AVRBrand
}

func GetEZBeqSlots() []int {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ slots", zap.Error(err))
		return []int{}
	}
	return config.Slots
}

func GetEZBeqPreferredAuthor() string {
	var config EZBEQConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load EZBEQ author", zap.Error(err))
		return ""
	}
	return config.PreferredAuthor
}

// Plex

func GetPlexUrl() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex url", zap.Error(err))
		return ""
	}

	config.URL = santizeURL(config.URL)

	return config.URL
}

func GetPlexToken() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex token", zap.Error(err))
		return ""
	}
	return config.Token
}

func GetPlexPort() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex port", zap.Error(err))
		return ""
	}

	if config.Port == "" {
		return "32400"
	}

	return config.Port
}

func GetPlexScheme() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex scheme", zap.Error(err))
		return ""
	}

	if config.Scheme == "" {
		return "http"
	}

	return config.Scheme
}

func GetPlexDeviceUUIDFilter() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex uuid", zap.Error(err))
		return ""
	}
	return config.DeviceUUIDFilter
}

func GetPlexOwnerNameFilter() string {
	var config PlexConfig
	if err := globalConfig.LoadConfig(context.Background(), &config); err != nil {
		gLog.Error("Failed to load Plex name", zap.Error(err))
		return ""
	}
	return config.OwnerNameFilter
}
