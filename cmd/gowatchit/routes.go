package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/models"
)

// RegisterRoutes registers the routes for the API contained in handlers.go
func RegisterRoutes(router *gin.Engine, webhookChan chan models.Event) {
	router.GET("/api/config", GetConfig)
	router.POST("/api/config", SaveConfig)
	router.GET("/api/logs", GetLogs)
	router.GET("/api/health", processHealthcheckWebhookGin)
	router.GET("/api/currentprofile", func(c *gin.Context) {
		profileHandler(c.Request.Context(), c)
	})
	router.POST("/api/minidsp", processMinidspHookGin)
	router.POST("/api/webhook", func(c *gin.Context) {
		processWebhook(c.Request.Context(), webhookChan, c)
	})
}
