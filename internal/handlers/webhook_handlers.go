package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/models"
)

func WebhookWorker(webhookChan <-chan models.Webhook, ready chan<- bool) {
	ready <- true
	for i := range webhookChan {
		log.Debugf("webhook received %v", i)
	}
}

func ProcessPlainWebhook(ctx context.Context, webhookChan chan<- models.Webhook, c *gin.Context) {
	var webhook models.Webhook
	err := c.BindJSON(&webhook)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	webhookChan <- webhook
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
