// implement the MediaPlayer interface for Plex
package plex

import (
	"context"
	"fmt"

	// "github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/models"
)

// PlexPlayer implements the MediaAPIClient interface
type PlexPlayer struct {
	Client mediaplayer.MediaAPIClient
	// TODO: all beq operations should be moved to orchestrator
	// HaClient  *homeassistant.HomeAssistantClient // TODO: interface
	// SearchRequest is a shared state for BEQ searching
	// SearchRequest *models.BeqSearchRequest
	// EventChan is a channel for sending processed webhook data to a background worker
	EventChan chan<- models.MediaPayload
	// SkipActions is a flag for skipping processing events while HDMI sync is running
	SkipActions *bool
}

// Ensure PlexPlayer implements MediaPlayer interface at compile time
var _ mediaplayer.MediaPlayer = (*PlexPlayer)(nil)

func NewPlexPlayer(scheme, serverURL, port string) (*PlexPlayer, error) {
	client, err := NewClient(scheme, serverURL, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create Plex client: %v", err)
	}

	return &PlexPlayer{
		Client: client,
	}, nil
}

// TODO: anything that needs to access BEQ and stuff needs to be moved to the player

// Implement MediaPlayerControlHandler methods
func (p *PlexPlayer) Play(ctx context.Context) error {
	// Implement Plex-specific play logic
	return p.Client.DoPlaybackAction(ctx, models.ActionPlay)
}

func (p *PlexPlayer) Pause(ctx context.Context) error {
	// Implement Plex-specific pause logic
	return p.Client.DoPlaybackAction(ctx, models.ActionPause)
}

func (p *PlexPlayer) Stop(ctx context.Context) error {
	// Implement Plex-specific stop logic
	return p.Client.DoPlaybackAction(ctx, models.ActionStop)
}

