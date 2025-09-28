package query

import (
	"context"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	configmodels "github.com/iloveicedgreentea/gowatchit/pkg/gen/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
)

type GetConfigQuery struct{}

type GetConfigHandler struct{}

func NewGetConfigHandler() (*GetConfigHandler, error) {
	return &GetConfigHandler{}, nil
}

func (h *GetConfigHandler) Handle(ctx context.Context, query *GetConfigQuery) (map[string]interface{}, error) {
	log := logger.GetLoggerFromContext(ctx)

	cfg := config.GetConfig()
	if cfg == nil {
		return nil, nil
	}

	configMap := make(map[string]interface{})

	ezbeqConfig := &configmodels.EZBEQConfig{}
	haConfig := &configmodels.HomeAssistantConfig{}
	hdmiConfig := &configmodels.HDMISyncConfig{}
	mainConfig := &configmodels.MainConfig{}

	configs := map[string]interface{}{
		"ezbeq":         ezbeqConfig,
		"homeassistant": haConfig,
		"hdmisync":      hdmiConfig,
		"main":          mainConfig,
	}

	for name, conf := range configs {
		if err := cfg.LoadConfig(ctx, conf); err != nil {
			log.Error("Failed to load config", zap.String("name", name), zap.Error(err))
			continue
		}
		configMap[name] = conf
	}

	// Load player configs if they exist
	playerConfig := &configmodels.PlayerConfig{}
	if err := cfg.LoadConfig(ctx, playerConfig); err == nil {
		configMap["player"] = playerConfig
	}

	return configMap, nil
}
