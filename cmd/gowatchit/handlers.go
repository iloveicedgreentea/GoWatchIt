package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/events"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func processWebhook(ctx context.Context, eventChan chan models.Event, c *gin.Context) {
	log := logger.GetLoggerFromContext(ctx)
	event, err := events.RequestToEvent(ctx, c.Request)
	if err != nil {
		log.Error("Error processing webhook",
			slog.Any("error", err),
		)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Try to send event, discard old one if channel is full
	select {
	case eventChan <- event: // Try to send
	default:
		// Channel is full, remove old event and send new one
		select {
		case <-eventChan: // Remove old event
		default:
		}
		eventChan <- event
	}

	c.JSON(200, gin.H{"message": "Webhook decoded successfully"})
}

func processHealthcheckWebhookGin(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// GetConfig returns all configurations from the database
func GetConfig(c *gin.Context) {
	cfg := config.GetConfig()
	if cfg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config manager not initialized"})
		return
	}

	// Create a map to store all configs
	configMap := make(map[string]interface{})

	// Load each config type
	ezbeqConfig := &models.EZBEQConfig{}
	haConfig := &models.HomeAssistantConfig{}
	jellyfinConfig := &models.JellyfinConfig{}
	hdmiConfig := &models.HDMISyncConfig{}
	plexConfig := &models.PlexConfig{}
	mainConfig := &models.MainConfig{}

	configs := map[string]interface{}{
		"ezbeq":         ezbeqConfig,
		"homeassistant": haConfig,
		"jellyfin":      jellyfinConfig,
		"hdmisync":      hdmiConfig,
		"plex":          plexConfig,
		"main":          mainConfig,
	}

	// Load each config
	for name, conf := range configs {
		if err := cfg.LoadConfig(c.Request.Context(), conf); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("Failed to load %s config: %v", name, err))
			continue
		}
		configMap[name] = conf
	}

	c.JSON(http.StatusOK, configMap)
}

// SaveConfig saves configurations to the database
func SaveConfig(c *gin.Context) {
	cfg := config.GetConfig()
	if cfg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config manager not initialized"})
		return
	}

	// Parse incoming JSON
	var configMap map[string]json.RawMessage
	if err := c.ShouldBindJSON(&configMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	// Map of config types
	configTypes := map[string]interface{}{
		"ezbeq":         &models.EZBEQConfig{},
		"homeassistant": &models.HomeAssistantConfig{},
		"jellyfin":      &models.JellyfinConfig{},
		"hdmisync":      &models.HDMISyncConfig{},
		"plex":          &models.PlexConfig{},
		"main":          &models.MainConfig{},
	}

	// Process each config section
	for name, data := range configMap {
		configStruct, exists := configTypes[name]
		if !exists {
			logger.GetLogger().Error(fmt.Sprintf("Unknown config type: %s", name))
			continue
		}

		// Unmarshal the config data
		if err := json.Unmarshal(data, configStruct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid %s config: %v", name, err)})
			return
		}

		// Save the config
		if err := cfg.SaveConfig(configStruct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save %s config: %v", name, err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configurations saved successfully"})
}

// GetLogs returns application logs (implement based on your logging system)
func GetLogs(c *gin.Context) {
	entries, err := logger.GetLogEntries()
	if err != nil {
		logger.GetLogger().Error("Failed to retrieve logs",
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}
