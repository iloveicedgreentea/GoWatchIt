package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/models"
)

// RegisterRoutes registers the routes for the API contained in handlers.go
func RegisterRoutes(router *gin.Engine, webhookChan chan models.Event) {
	// router.GET("/config", GetConfig)
	// router.POST("/config", SaveConfig)
	// router.GET("/logs", GetLogs)
	router.GET("/health", processHealthcheckWebhookGin)
	router.POST("/webhook", func(c *gin.Context) {
		processWebhook(c.Request.Context(), webhookChan, c)
	})
	router.Static("/assets", "./assets")
	// router.GET("/config-exists", api.ConfigExists)
	// router.GET("/get-config", api.GetConfig)
	// router.POST("/save-config", api.SaveConfig)

	// TODO: route to export config as json
	// TODO: route to import config from json
}
