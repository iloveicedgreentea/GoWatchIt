package mappers

import (
	"context"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/pkg/beq"
	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/events"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/mediaplayer"
	"go.uber.org/zap"
)

// EventToBEQPayload maps an Event and MediaPlayer to a BEQPayload for loading profiles
func EventToBEQPayload(ctx context.Context, event *events.Event, player mediaplayer.MediaPlayer) (*beq.BEQPayload, error) {
	log := logger.GetLoggerFromContext(ctx)

	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}
	if player == nil {
		return nil, fmt.Errorf("player is nil")
	}

	// Get codec from player
	codec, err := player.GetAudioCodec(ctx, event)
	if err != nil {
		log.Error("Failed to get audio codec",
			zap.Error(err),
			zap.String("title", event.Metadata.Title),
		)
		return nil, fmt.Errorf("failed to get audio codec: %w", err)
	}

	// Get edition from player
	edition, err := player.GetEdition(ctx, event)
	if err != nil {
		log.Warn("Failed to get edition, using None",
			zap.Error(err),
			zap.String("title", event.Metadata.Title),
		)
		// Continue with EditionNone if we can't determine edition
	}

	// Build the BEQ payload
	payload := &beq.BEQPayload{
		Codec:           codec,
		DryrunMode:      config.IsBeqDryRun(ctx),
		Edition:         edition,
		EntryID:         "", // Will be populated by search
		MediaType:       event.Metadata.Type,
		MVAdjust:        0, // Will be populated by search
		PreferredAuthor: config.GetEZBeqPreferredAuthor(ctx),
		SkipSearch:      false, // Always search on initial load
		Slots:           config.GetEZBeqSlots(ctx),
		Title:           event.Metadata.Title,
		TMDB:            event.Metadata.TMDB,
		Year:            event.Metadata.Year,
	}

	log.Debug("Mapped Event to BEQPayload",
		zap.String("title", payload.Title),
		zap.String("codec", string(payload.Codec)),
		zap.String("edition", string(payload.Edition)),
		zap.String("mediaType", string(payload.MediaType)),
		zap.String("tmdb", payload.TMDB),
		zap.Int("year", payload.Year),
	)

	return payload, nil
}
