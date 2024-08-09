package mediaplayer

import (
	"context"

	"github.com/iloveicedgreentea/go-plex/models"
)

// MediaPlayer defines the interface for direct player control
type MediaPlayer interface {
	// Play starts playback
	Play(ctx context.Context) error
	// Pause pauses playback
	Pause(ctx context.Context) error
	// Stop stops playback
	Stop(ctx context.Context) error
	// Add other control methods as needed
}

// MediaAPIClient defines the interface for media server client operations
type MediaAPIClient interface {
	DoPlaybackAction(ctx context.Context, action models.Action) error
	// Get edition from metadata like Extended, Director's Cut, etc.
	GetEdition(ctx context.Context, payload models.Event) (models.Edition, error)
	// Get media payload from a media server like the metadata for a movie
	// GetMediaData(ctx context.Context, key string) (models.DataMediaContainer, error)
	// Get codec from metadata not as reliable as AVR data but more portable
	// TODO: GetAudioCodecFromAVR/GetCodecFromSession should be a flag in this method and/or private
	GetAudioCodec(ctx context.Context, data models.Event) (models.CodecName, error)
	// Get codec from session data. Sometimes slower than metadata but could be concurrent and more reliable
	// GetCodecFromSession(ctx context.Context, uuid string) (string, error)
	// Get session data from a player like plex
	// GetSessionData(ctx context.Context) (models.MediaSession, error)
}
