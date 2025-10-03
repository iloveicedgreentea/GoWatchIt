package config

import (
	"context"
	"strings"

	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
)

func santizeURL(url string) string {
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	return url
}

func getHdmiSyncConfig(ctx context.Context) *configmodels.HDMISyncConfig {
	var config configmodels.HDMISyncConfig
	if err := globalConfig.LoadConfig(ctx, &config); err != nil {
		log := logger.GetLoggerFromContext(ctx)
		log.Error("Failed to load HDMISync config", zap.Error(err))
	}

	return &config
}

func getHAConfig(ctx context.Context) *configmodels.HomeAssistantConfig {
	var config configmodels.HomeAssistantConfig
	if err := globalConfig.LoadConfig(ctx, &config); err != nil {
		log := logger.GetLoggerFromContext(ctx)
		log.Error("Failed to load HomeAssistant config", zap.Error(err))
	}
	return &config
}

func getEZBEQConfig(ctx context.Context) *configmodels.EZBEQConfig {
	var config configmodels.EZBEQConfig
	if err := globalConfig.LoadConfig(ctx, &config); err != nil {
		log := logger.GetLoggerFromContext(ctx)
		log.Error("Failed to load EZBEQ config", zap.Error(err))
	}
	return &config
}

func getPlayerConfig(ctx context.Context) *configmodels.PlayerConfig {
	var config configmodels.PlayerConfig
	if err := globalConfig.LoadConfig(ctx, &config); err != nil {
		log := logger.GetLoggerFromContext(ctx)
		log.Error("Failed to load Player config", zap.Error(err))
	}
	return &config
}

// HDMI
func GetHDMISyncSource(ctx context.Context) string {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.Source
	}
	return ""
}

func GetHDMISyncEnvyName(ctx context.Context) string {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.Envy
	}
	return ""
}

func GetHDMISyncSeconds(ctx context.Context) string {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.Time
	}
	return ""
}

func GetHDMISyncPlayerIP(ctx context.Context) string {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.PlayerIP
	}
	return ""
}

func GetHDMISyncMachineIdentifier(ctx context.Context) string {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.PlayerMachineIdentifier
	}
	return ""
}

// Home Assistant

func GetHomeAssistantUrl(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		url := santizeURL(config.Url)
		if url == "" {
			return "homeassistant.local"
		}
		return url
	}
	return "homeassistant.local"
}

func GetHomeAssistantScheme(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		if config.Scheme == "" {
			return "http"
		}
		return string(config.Scheme)
	}
	return "http"
}

func GetHomeAssistantToken(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		return config.Token
	}
	return ""
}

func GetHomeAssistantPort(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		if config.Port == "" {
			return "8123"
		}
		return config.Port
	}
	return "8123"
}

func GetHomeAssistantRemoteEntityName(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		return config.RemoteEntityName
	}
	return ""
}

func GetHomeAssistantNotifyEndpointName(ctx context.Context) string {
	if config := getHAConfig(ctx); config != nil {
		return config.NotifyEndpointName
	}
	return ""
}

func GetHomeAssistantNotifyDisplayTime(ctx context.Context) int {
	if config := getHAConfig(ctx); config != nil {
		// set default to 15
		if config.NotifyDisplayTime == 0 {
			// Default to 15 seconds (15000 ms)
			return 15 * 1000
		}
		// Convert configured seconds to milliseconds
		return int(config.NotifyDisplayTime) * 1000
	}
	return 15 * 1000
}

// EZBeq
func GetEZBeqUrl(ctx context.Context) string {
	if config := getEZBEQConfig(ctx); config != nil {
		url := santizeURL(config.Url)
		if url == "" {
			return "ezbeq.local"
		}
		return url
	}
	return "ezbeq.local"
}

func GetEZBeqScheme(ctx context.Context) string {
	if config := getEZBEQConfig(ctx); config != nil {
		if config.Scheme == "" {
			return "http"
		}
		return string(config.Scheme)
	}
	return "http"
}

func GetEZBeqPort(ctx context.Context) string {
	if config := getEZBEQConfig(ctx); config != nil {
		if config.Port == "" {
			return "8080"
		}
		return config.Port
	}
	return "8080"
}

func GetEZBeqAvrURL(ctx context.Context) string {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.AvrURL
	}
	return ""
}

func GetEZBeqAvrBrand(ctx context.Context) configmodels.EZBEQConfigAvrBrand {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.AvrBrand
	}
	return ""
}

func GetEZBeqSlots(ctx context.Context) []int32 {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.Slots
	}
	return []int32{}
}

func GetEZBeqPreferredAuthors(ctx context.Context) []string {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.PreferredAuthors
	}
	return []string{}
}

// Player

func GetPlayerType(ctx context.Context) configmodels.Player {
	if config := getPlayerConfig(ctx); config != nil {
		return config.PlayerType
	}
	return ""
}

func GetPlayerURL(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		return santizeURL(config.Url)
	}
	return ""
}

func GetPlayerToken(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		return config.Token
	}
	return ""
}

func GetPlayerPort(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		if config.Port == "" {
			return "32400"
		}
		return config.Port
	}
	return "32400"
}

func GetPlayerScheme(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		if config.Scheme == "" {
			return "http"
		}
		return string(config.Scheme)
	}
	return "http"
}

func GetPlayerDeviceUUIDFilter(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		return config.DeviceUUIDFilter
	}
	return ""
}

func GetPlayerOwnerNameFilter(ctx context.Context) string {
	if config := getPlayerConfig(ctx); config != nil {
		return config.OwnerNameFilter
	}
	return ""
}

// HDMI
func IsHDMISyncEnabled(ctx context.Context) bool {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.Enabled
	}
	return false
}

func IsSignalSourceTime(ctx context.Context) bool {
	if config := getHdmiSyncConfig(ctx); config != nil {
		return config.Source == "time"
	}
	return false
}

// Player
func IsPlayerEnabled(ctx context.Context) bool {
	// PlayerConfig doesn't have Enabled field, need to check if it's configured
	if config := getPlayerConfig(ctx); config != nil {
		// Check if player is configured (has URL or token)
		return config.Url != "" || config.Token != ""
	}
	return false
}

// BEQ
func IsBeqEnabled(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.Enabled
	}
	return false
}

func IsBeqTVEnabled(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.EnableTVBEQ
	}
	return false
}

func IsBeqNotifyOnLoadEnabled(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.NotifyOnLoad
	}
	return false
}

func IsBeqNotifyOnUnLoadEnabled(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.NotifyOnUnload
	}
	return false
}

func IsBeqDryRun(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.DryRun
	}
	return false
}

func IsBeqLooseEditionMatching(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.LooseEditionMatching
	}
	return false
}

func IsBeqSkipEditionMatching(ctx context.Context) bool {
	if config := getEZBEQConfig(ctx); config != nil {
		return config.SkipEditionMatching
	}
	return false
}

// Home Assistant
func IsHomeAssistantEnabled(ctx context.Context) bool {
	if config := getHAConfig(ctx); config != nil {
		return config.Enabled
	}
	return false
}

func IsJellyfinSkipTMDB(ctx context.Context) bool {
	// Check PlayerConfig for SkipTMDB
	if config := getPlayerConfig(ctx); config != nil {
		return config.SkipTMDB
	}
	return false
}
