package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iloveicedgreentea/gowatchit/pkg/beq"
	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	eventType "github.com/iloveicedgreentea/gowatchit/pkg/events"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/pkg/plex"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/events"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/mediaplayer"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/mappers"
	"go.uber.org/zap"
)

type ProcessWebhookCommand struct {
	Request *http.Request
}

type ProcessWebhookHandler struct {
	BeqClient *beq.BeqClient
}

func NewProcessWebhookHandler(beqClient *beq.BeqClient) (*ProcessWebhookHandler, error) {
	if beqClient == nil {
		return nil, errors.New("beqClient is nil")
	}
	return &ProcessWebhookHandler{
		BeqClient: beqClient,
	}, nil
}

// Handle loads BEQ to device
func (h *ProcessWebhookHandler) Handle(ctx context.Context, cmd *ProcessWebhookCommand) error {
	return ProcessWebhook(ctx, cmd, h)
}

func getClient(ctx context.Context) (mediaplayer.MediaPlayer, error) {
	log := logger.GetLoggerFromContext(ctx)

	// from config, choose which player is configured
	playerConfig := config.GetPlayerType(ctx)
	log.Info("Processing request with media player",
		zap.String("player", string(playerConfig)),
	)

	var player mediaplayer.MediaPlayer
	var err error

	switch playerConfig {
	case config.PlayerPlex:
		player, err = plex.NewClient(ctx)
	// TODO add other players
	// case config.PlayerJellyfin:
	// 	player, err = jellyfin.NewClient(config.GetJellyfinConfig())
	// case config.PlayerHomeAssistant:
	// 	player, err = homeassistant.NewClient(config.GetHomeAssistantConfig())
	default:
		return nil, errors.New("unsupported media player")
	}
	if err != nil {
		return nil, err
	}

	return player, nil
}

func ProcessWebhook(ctx context.Context, cmd *ProcessWebhookCommand, handler *ProcessWebhookHandler) error {
	log := logger.GetLoggerFromContext(ctx)

	player, err := getClient(ctx)
	if err != nil {
		return err
	}

	if cmd.Request == nil {
		log.Error("Request is nil")
		return fmt.Errorf("request is nil")
	}
	if player == nil {
		log.Error("Media player is nil")
		return fmt.Errorf("media player is nil")
	}
	event, err := events.RequestToEvent(ctx, player, cmd.Request)
	if err != nil {
		return fmt.Errorf("failed to convert webhook to event: %w", err)
	}
	if event == nil {
		log.Error("Event is nil after conversion")
		return fmt.Errorf("event is nil after conversion")
	}

	log.Debug("Received event",
		zap.Any("event", event),
	)

	// Map event to BEQ payload
	beqPayload, err := mappers.EventToBEQPayload(ctx, event, player)
	if err != nil {
		return fmt.Errorf("failed to map event to BEQ payload: %w", err)
	}

	if handler == nil {
		log.Error("ProcessWebhookHandler is nil, skipping BEQ operations")
		return nil
	}
	if handler.BeqClient == nil {
		log.Error("BeqClient in handler is nil, skipping BEQ operations")
		return nil
	}

	// Handle BEQ profile based on event action
	switch event.Action {
	case eventType.ActionPlay, eventType.ActionResume:
		// Load BEQ profile
		if err := handler.BeqClient.LoadBeqProfile(ctx, beqPayload); err != nil {
			log.Error("Failed to load BEQ profile",
				zap.Error(err),
				zap.String("title", beqPayload.Title),
			)
			return fmt.Errorf("failed to load BEQ profile: %w", err)
		}

		log.Info("Successfully loaded BEQ profile",
			zap.String("title", beqPayload.Title),
			zap.String("codec", string(beqPayload.Codec)),
			zap.String("edition", string(beqPayload.Edition)),
		)

	case eventType.ActionPause, eventType.ActionStop:
		// Unload BEQ profile
		if err := handler.BeqClient.UnloadBeqProfile(ctx, beqPayload); err != nil {
			log.Error("Failed to unload BEQ profile",
				zap.Error(err),
				zap.String("title", beqPayload.Title),
			)
			return fmt.Errorf("failed to unload BEQ profile: %w", err)
		}

		log.Info("Successfully unloaded BEQ profile",
			zap.String("title", beqPayload.Title),
		)

	default:
		log.Info("Event action not supported for BEQ operations",
			zap.String("action", string(event.Action)),
		)
		return nil
	}

	// TODO: support hdmi sync

	return nil
}
