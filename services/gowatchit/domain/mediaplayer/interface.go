package mediaplayer

import (
	"context"

	"github.com/iloveicedgreentea/gowatchit/pkg/codecs"
	"github.com/iloveicedgreentea/gowatchit/pkg/editions"
	"github.com/iloveicedgreentea/gowatchit/pkg/events"
)

// MediaPlayer defines actions and methods every player must implement like Plex, JF, etc
type MediaPlayer interface {
	// GetStatus returns a status object
	// GetStatus(ctx context.Context) (*Status, error)

	// GetAudioCodec returns the audio codec used in the media player
	GetAudioCodec(ctx context.Context, event *events.Event) (codecs.Codec, error)
	// GetEdition returns the edition of the media being played
	GetEdition(ctx context.Context, event *events.Event) (editions.Edition, error)
	// ActionPause pauses the media player
	ActionPause(ctx context.Context) error
	// ActionPlay plays the media player
	ActionPlay(ctx context.Context) error
	// ActionStop stops the media player
	ActionStop(ctx context.Context) error
}
