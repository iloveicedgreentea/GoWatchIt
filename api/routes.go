package api

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes registers the routes for the API contained in handlers.go
func RegisterRoutes(router *gin.Engine) {
	router.GET("/config", GetConfig)
	router.POST("/config", SaveConfig)
	router.GET("/logs", GetLogs)
}