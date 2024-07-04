// File: plex/plex.go

package plex

import (
	"context"
	"fmt"

	"github.com/iloveicedgreentea/go-plex/models"
)

// PlexPlayer implements the MediaPlayer interface for Plex
type PlexPlayer struct {
	plexClient *PlexClient
	// TODO: add all clients and stuff here or to base mediaplayer
	// Add any Plex-specific fields here
}

// NewPlexPlayer creates a new PlexPlayer instance
func NewPlexPlayer(serverURL, port string) *PlexPlayer {
	return &PlexPlayer{
		plexClient: NewClient(serverURL, port),
	}
}

// Play implements the Play method for Plex
func (p *PlexPlayer) Play(ctx context.Context) error {
	return p.plexClient.DoPlaybackAction(models.ActionPlay)
}

// Pause implements the Pause method for Plex
func (p *PlexPlayer) Pause(ctx context.Context) error {
	return p.plexClient.DoPlaybackAction(models.ActionPause)
}

// Implement other methods of the MediaPlayer interface...

// RouteEvent routes the event to the appropriate method
func (p *PlexPlayer) RouteEvent(ctx context.Context, eventType string, payload models.MediaPayload) error {
	// TODO pass this to handler or reimplement hanlder here
	switch eventType {
	case "media.play":
		return p.Play(ctx)
	case "media.pause":
		return p.Pause(ctx)
	// Add cases for other event types
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}
}

// GetMediaData retrieves media data from Plex
func (p *PlexPlayer) GetMediaData(ctx context.Context, key string) (interface{}, error) {
	// Implement logic to fetch media data from Plex
	// This is a placeholder implementation
	return map[string]string{"key": key}, nil
}
