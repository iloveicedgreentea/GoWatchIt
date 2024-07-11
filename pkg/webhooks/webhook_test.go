package webhooks

import (
	"net/http"
	"context"
	// "github.com/iloveicedgreentea/go-plex/pkg/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/pkg/plex"
	// "github.com/iloveicedgreentea/go-plex/pkg/shield"
)

func main() {
    webhookHandler := NewWebhookHandler()

    // Set up Plex
    plexPlayer := plex.NewPlexPlayer()
    plexEventChan := make(chan interface{}, 10)
    plexProcessor := &PlexWebhookProcessor{
        Player: plexPlayer,
        EventChan: plexEventChan,
    }
    webhookHandler.RegisterProcessor("plex", plexProcessor)

    // // Set up Shield
    // shieldPlayer := shield.NewShieldPlayer()
    // shieldEventChan := make(chan interface{}, 10)
    // shieldProcessor := &ShieldWebhookProcessor{
    //     Player: shieldPlayer,
    //     EventChan: shieldEventChan,
    // }
    // webhookHandler.RegisterProcessor("shield", shieldProcessor)

    // Start workers
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    readyChan := make(chan bool)
    go plexProcessor.StartWorker(ctx, readyChan)
    // go shieldProcessor.StartWorker(ctx, readyChan)

    // Wait for workers to be ready
    <-readyChan
    <-readyChan

    // Set up HTTP server
    http.HandleFunc("/webhook/", webhookHandler.HandleWebhook)
    
    // Start server
    http.ListenAndServe(":8080", nil)
}