package webhooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
)

// add routes to router
func addRoutes(router *gin.Engine, app *app.App) {
	router.GET("/api/v1/config", GetConfig)
	router.GET("/api/v1/debug", DebugRoute)
	router.POST("/api/v1/webhook", func(ctx *gin.Context) {
		GetWebhook(ctx, app)
	})
	// get beq profile
	router.GET("/api/v1/profile", func(ctx *gin.Context) {
		GetBeqProfile(ctx, app)
	})
}

// TODO: make route for configs
func GetConfig(c *gin.Context) {}

func DebugRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// GetWebhook receives the webhook request and routes it
func GetWebhook(c *gin.Context, app *app.App) {
	if app == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	cmd := &command.ProcessesWebhookCommand{Request: c.Request}
	err := app.Commands.ProcessWebhook.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func GetBeqProfile(c *gin.Context, app *app.App) {
	if app == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app is nil"})
		return
	}

	qu := &query.GetBeqQuery{}
	profile, err := app.Queries.GetBeq.Handle(c.Request.Context(), qu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
