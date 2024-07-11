package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/api"
	"github.com/iloveicedgreentea/go-plex/internal/handlers"
)

func addRoutes(r *gin.Engine) {
	/* ###############################
		Routes
	   ############################## */
	// healthcheck
	r.GET("/health", handlers.ProcessHealthcheckWebhookGin)
	// non-plex based webhook
	r.POST("/webhook", func(c *gin.Context) {
		log.Debug("webhook received")
		// handlers.ProcessPlainWebhook(c.Request.Context(), webhookChan, c)
	})
	// Add plex webhook handler
	r.POST("/plexwebhook", func(c *gin.Context) {
		log.Debug("plexwebhook received")
		// handlers.ProcessWebhook(c.Request.Context(), plexChan, c)
	})
	r.POST("/minidspwebhook", func(c *gin.Context) {
		log.Debug("minidspwebhook received")
		// handlers.ProcessMinidspWebhook(minidspChan, c)
	})
	r.POST("/jellyfinwebhook", func(c *gin.Context) {
		// handlers.ProcessJfWebhook(jfChan, c)
	})
	r.Static("/assets", "./assets")
	r.GET("/config-exists", api.ConfigExists)
	r.GET("/get-config", api.GetConfig)
	r.POST("/save-config", api.SaveConfig)
	// TODO: add generic webhook endpoint, maybe mqtt?
}
