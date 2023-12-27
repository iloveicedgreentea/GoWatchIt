package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)
func ProcessHealthcheckWebhookGin(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}