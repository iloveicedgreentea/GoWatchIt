// File: mediaplayer/mediaplayer.go

package mediaplayer

import (
	"context"
	"net/http"

	"github.com/iloveicedgreentea/go-plex/models"
)

// MediaPlayer is a generic interface for media players like plex or jellyfin
type MediaPlayer interface {
	// TODO: add play/pause functions for webhook response? probably not since its in event router
	Play(ctx context.Context) error
	Pause(ctx context.Context) error
	Stop(ctx context.Context) error
	Resume(ctx context.Context) error
	Scrobble(ctx context.Context) error
	GetAudioCodec(ctx context.Context, payload models.MediaPayload) (string, error)
	GetEdition(ctx context.Context, payload models.DataMediaContainer) (models.Edition, error)
	RouteEvent(ctx context.Context, eventType string, payload models.MediaPayload) error
	GetMediaData(ctx context.Context, key string) (interface{}, error)
	ProcessWebhook(ctx context.Context, payload *http.Request) error
}

// MediaPlayerFactory manages different media player instances
type MediaPlayerFactory struct {
	players map[string]interface{}
}

// NewMediaPlayerFactory creates a new MediaPlayerFactory
func NewMediaPlayerFactory() *MediaPlayerFactory {
	return &MediaPlayerFactory{
		players: make(map[string]interface{}),
	}
}

// RegisterPlayer registers a media player with the factory
func (f *MediaPlayerFactory) RegisterPlayer(playerType string, player interface{}) {
	f.players[playerType] = player
}

// GetPlayer retrieves a media player instance
func (f *MediaPlayerFactory) GetPlayer(playerType string) (interface{}, bool) {
	player, ok := f.players[playerType]
	return player, ok
}
