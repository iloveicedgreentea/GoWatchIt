package command

import (
	"context"
	"errors"
	"net/http"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/pkg/plex"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/mediaplayer"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/ports/events"
	"go.uber.org/zap"
)

type ProcessesWebhookCommand struct {
	Request *http.Request
}

type ProcessesWebhookHandler struct{}

func NewProcessesWebhookHandler(ctx context.Context) (*ProcessesWebhookHandler, error) {
	return &ProcessesWebhookHandler{}, nil
}

func getClient(ctx context.Context) (mediaplayer.MediaPlayer, error) {
	log := logger.GetLoggerFromContext(ctx)
	// from config, choose which player is configured
	playerConfig := config.GetPlayerType(ctx)
	log.Debug("Using media player",
		zap.String("player", string(playerConfig)),
	)

	var player mediaplayer.MediaPlayer
	var err error

	switch playerConfig {
	case config.PlayerPlex:
		player, err = plex.NewClient(ctx)
	// TODO add other players
	// case config.PlayerJellyfin:
	// 	player, err = mediaplayer.NewJellyfinPlayer(config.GetJellyfinConfig())
	// case config.PlayerHomeAssistant:
	// 	player, err = mediaplayer.NewHomeAssistantPlayer(config.GetHomeAssistantConfig())
	default:
		return nil, errors.New("unsupported media player")
	}
	if err != nil {
		return nil, err
	}

	return player, nil
}

// Handle loads BEQ to device
func (h *ProcessesWebhookHandler) Handle(ctx context.Context, cmd *ProcessesWebhookCommand) error {
	player, err := getClient(ctx)
	if err != nil {
		return err
	}
	return events.ProcessWebhook(ctx, player, cmd.Request)
}
