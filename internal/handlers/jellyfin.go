package handlers

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/jellyfin"
	"github.com/iloveicedgreentea/go-plex/models"
)

func ProcessJfWebhook(jfChan chan<- models.JellyfinWebhook, c *gin.Context) {
	// send payload to worker
	log.Debug("Sending payload to JellyfinWorker")
	r := c.Request.Body
	defer r.Close()
	read, err := io.ReadAll(r)
	if err != nil {
		log.Errorf("Error reading request body: %v", err)
	}

	// log.Debugf("ProcessJfWebhook Request: %v", string(read))
	var payload models.JellyfinWebhook
	err = json.Unmarshal(read, &payload)
	if err != nil {
		log.Errorf("Error decoding payload: %v", err)
	}
	log.Debugf("Payload: %#v", payload)
	// respond to request with 200
	c.JSON(200, gin.H{"status": "ok"})
	// send payload to worker
	jfChan <- payload
}

// entry point for background tasks
func JellyfinWorker(jfChan <-chan models.JellyfinWebhook, readyChan chan<- bool) {
	log.Info("JellyfinWorker started")

	// Server Info
	jellyfinClient := jellyfin.NewClient(config.GetString("jellyfin.url"), config.GetString("jellyfin.port"), config.GetString("jellyfin.playerMachineIdentifier"), config.GetString("jellyfin.playerIP"))

	readyChan <- true
	log.Info("JellyfinWorker is ready")
	// block forever until closed so it will wait in background for work
	for i := range jfChan {
		log.Debug("Sending new payload to eventRouter")
		log.Debug(i)
		// if its not an empty struct
		if i != (models.JellyfinWebhook{}) {
			// get metadata
			metadata, err := jellyfinClient.GetMetadata(i.UserID, i.ItemID)
			if err != nil {
				log.Errorf("Error getting metadata from jellyfin API: %v", err)
			}
			codec, displayTitle, err := jellyfinClient.GetCodec(metadata)
			if err != nil {
				log.Errorf("Error getting codec: %v", err)
			}
			log.Debugf("Response: %v, %v", codec, displayTitle)
		}
		log.Debug("eventRouter done processing payload")
	}

	log.Debug("JellyfinWorker worker stopped")
}
