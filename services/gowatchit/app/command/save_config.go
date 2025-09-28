package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
)

type SaveConfigCommand struct {
	ConfigData map[string]json.RawMessage
}

type SaveConfigHandler struct{}

func NewSaveConfigHandler() (*SaveConfigHandler, error) {
	return &SaveConfigHandler{}, nil
}

func (h *SaveConfigHandler) Handle(ctx context.Context, cmd *SaveConfigCommand) error {
	log := logger.GetLoggerFromContext(ctx)

	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("config manager not initialized")
	}

	configTypes := map[string]interface{}{
		"ezbeq":         &configmodels.EZBEQConfig{},
		"homeassistant": &configmodels.HomeAssistantConfig{},
		"hdmisync":      &configmodels.HDMISyncConfig{},
		"main":          &configmodels.MainConfig{},
	}

	for name, data := range cmd.ConfigData {
		if name == "players" {
			var playerConfigs []configmodels.PlayerConfig
			if err := json.Unmarshal(data, &playerConfigs); err != nil {
				return fmt.Errorf("invalid player configs: %w", err)
			}

			for i := range playerConfigs {
				if err := cfg.SaveConfig(&playerConfigs[i]); err != nil {
					return fmt.Errorf("failed to save player config: %w", err)
				}
			}
			continue
		}

		configStruct, exists := configTypes[name]
		if !exists {
			log.Warn("Unknown config type", zap.String("name", name))
			continue
		}

		if err := json.Unmarshal(data, configStruct); err != nil {
			return fmt.Errorf("invalid %s config: %w", name, err)
		}

		if err := cfg.SaveConfig(configStruct); err != nil {
			return fmt.Errorf("failed to save %s config: %w", name, err)
		}
	}

	return nil
}
