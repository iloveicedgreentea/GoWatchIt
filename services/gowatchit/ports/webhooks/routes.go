package webhooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
)

// add routes to router
func addRoutes(router *gin.Engine, appInst *app.App) {
	router.GET("/api/v1/config", GetConfig)
	router.GET("/api/v1/debug", DebugRoute)
	router.POST("/api/v1/webhook", func(ctx *gin.Context) {
		PostWebhook(ctx, appInst)
	})
	// get beq profile
	router.GET("/api/v1/profile", func(ctx *gin.Context) {
		GetBeqProfile(ctx, appInst)
	})
}

// TODO: make route for configs
func GetConfig(c *gin.Context) {}

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
