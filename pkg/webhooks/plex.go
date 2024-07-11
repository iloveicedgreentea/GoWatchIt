package webhooks

import (
    "net/http"
    "context"
    "github.com/iloveicedgreentea/go-plex/pkg/mediaplayer"
)

type PlexWebhookProcessor struct {
    Player mediaplayer.MediaPlayer
    EventChan <-chan interface{}
}

func (p *PlexWebhookProcessor) Process(r *http.Request) error {
    // Parse multipart form data
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        return err
    }
    
    // Extract and process Plex-specific data
    // Send processed data to EventChan
    
    return nil
}

func (p *PlexWebhookProcessor) StartWorker(ctx context.Context, ready chan<- bool) {
    ready <- true
    for {
        select {
        case event := <-p.EventChan:
            // Process Plex event
        case <-ctx.Done():
            return
        }
    }
}
