package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
)

// add routes to router
func addRoutes(router *gin.Engine, appInst *app.App) {
	router.GET("/api/v1/config", func(ctx *gin.Context) {
		GetConfig(ctx, appInst)
	})
	router.POST("/api/v1/config", func(ctx *gin.Context) {
		SaveConfig(ctx, appInst)
	})
	router.GET("/api/v1/logs", func(ctx *gin.Context) {
		GetLogs(ctx, appInst)
	})
	router.GET("/api/v1/health", HealthCheck)
	router.GET("/api/v1/debug", DebugRoute)
	router.POST("/api/v1/webhook", func(ctx *gin.Context) {
		PostWebhook(ctx, appInst)
	})
	// get beq profile
	router.GET("/api/v1/profile", func(ctx *gin.Context) {
		GetBeqProfile(ctx, appInst)
	})
}

// GetConfig returns all configurations from the database
func GetConfig(c *gin.Context, appInst *app.App) {
	if appInst == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	qu := &query.GetConfigQuery{}
	configs, err := appInst.Queries.GetConfig.Handle(c.Request.Context(), qu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// SaveConfig saves configurations to the database
func SaveConfig(c *gin.Context, appInst *app.App) {
	if appInst == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	var configMap map[string]json.RawMessage
	if err := c.ShouldBindJSON(&configMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	cmd := &command.SaveConfigCommand{ConfigData: configMap}
	err := appInst.Commands.SaveConfig.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configurations saved successfully"})
}

// GetLogs returns application logs
func GetLogs(c *gin.Context, appInst *app.App) {
	if appInst == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	qu := &query.GetLogsQuery{}
	logs, err := appInst.Queries.GetLogs.Handle(c.Request.Context(), qu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// HealthCheck returns a simple health status
func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func DebugRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// PostWebhook receives the webhook request and routes it
func PostWebhook(c *gin.Context, appInst *app.App) {
	if appInst == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	cmd := &command.ProcessWebhookCommand{Request: c.Request}
	err := appInst.Commands.ProcessWebhook.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func GetBeqProfile(c *gin.Context, appInst *app.App) {
	if appInst == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	qu := &query.GetBeqQuery{}
	profile, err := appInst.Queries.GetBeq.Handle(c.Request.Context(), qu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
