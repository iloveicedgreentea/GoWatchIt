// this file implements MediaPlayer interface for Plex

package plex

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/iloveicedgreentea/go-plex/pkg/utils"
)

// PlexPlayer implements the MediaPlayer interface for Plex
type PlexPlayer struct {
	PlexClient *PlexClient
	BeqClient *ezbeq.BeqClient
	HaClient *homeassistant.HomeAssistantClient
	SearchRequest *models.SearchRequest
	// TODO: chan
	plexChan chan<- models.PlexWebhookPayload
	skipActions *bool
	// Add any Plex-specific fields here
}

// NewPlexPlayer creates a new PlexPlayer instance
func NewPlexPlayer(scheme, serverURL, port string, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, c chan<- models.PlexWebhookPayload ) (*PlexPlayer, error) {
	if !utils.ValidateHttpScheme(scheme) {
		return nil, fmt.Errorf("invalid http scheme: %s", scheme)

	}
	return &PlexPlayer{
		PlexClient: NewClient(scheme, serverURL, port),
		BeqClient: beqClient,
		HaClient: haClient,
		SearchRequest: &models.SearchRequest{},
		plexChan: c,

	}, nil
}

// Play implements the Play method for Plex
func (p *PlexPlayer) Play(ctx context.Context) error {
	return p.PlexClient.DoPlaybackAction(models.ActionPlay)
}

// Pause implements the Pause method for Plex
func (p *PlexPlayer) Pause(ctx context.Context) error {
	return p.PlexClient.DoPlaybackAction(models.ActionPause)
}

// Play implements the Play method for Plex
func (p *PlexPlayer) Stop(ctx context.Context) error {
	return p.PlexClient.DoPlaybackAction(models.ActionStop)
}

func (p *PlexPlayer) GetEdition(ctx context.Context, payload models.DataMediaContainer) (models.Edition, error) {
	return p.getEditionName(ctx, payload)
}

// Implement other methods of the MediaPlayer interface...

// RouteEvent routes the event to the appropriate method
func (p *PlexPlayer) RouteEvent(ctx context.Context, eventType string, payload models.MediaPayload) error {
	// TODO pass this to handler or reimplement hanlder here
	switch eventType {
	case "media.play":
		// TODO: handle play event
	case "media.pause":
		// TODO: handle pause event
	case "media.stop":
	case "media.resume":
		// TODO: handle stop event
	case "media.scrobble":
		// TODO: handle scrobble event
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	return nil
}

// GetMediaData retrieves media data from Plex
func (p *PlexPlayer) GetMediaData(ctx context.Context, key string) (interface{}, error) {
	// Implement logic to fetch media data from Plex
	// This is a placeholder implementation
	return map[string]string{"key": key}, nil
}

// TODO: return this to the endpoint handler
func (p *PlexPlayer) ProcessWebhook(ctx context.Context, request *http.Request) error {
	// Implement logic to process webhook payload
	// This is a placeholder implementation
	return p.ProcessPlexWebhook(ctx, request)
}