package mediaplayer

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

// Implement MediaEventHandler methods
func HandlePlay(ctx context.Context, payload *models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) error {
	log := logger.GetLoggerFromContext(ctx)
	// Check if context is already cancelled before starting lets say you play but then stop, this should stop processing
	if ctx.Err() != nil {
		log.Debug("Context already cancelled, stopping processing",
			"error", ctx.Err(),
			"func", "HandlePlay",
		)
		return nil
	}
	log.Debug("handleplay searchRequest",
		slog.Any("searchRequest", searchRequest),
	)

	var err error
	var innerWg sync.WaitGroup

	// Perform HDMI sync
	// Call the sync function which will check if its enabled
	if config.IsHDMISyncEnabled() {
		if !strings.EqualFold(string(payload.Metadata.Type), string(models.MediaTypeMovie)) && config.IsSignalSourceTime() {
			log.Debug("skipping sync for non-movie type and time source")
		} else {
			// innerWg.Add(1) // TODO: enable
			go func() {
				if ctx.Err() != nil {
					log.Debug("mediaPlay was cancelled before hdmi sync")
					return // Exit early if context is cancelled
				}

				// optimistically try to hdmi sync. Will return if disabled
				// TODO: implement this ensure it calls innerWg.done

				// common.WaitForHDMISync(innerWg, skipActions, haClient, client)
			}()
		}
	}

	// dont need to set skipActions here because it will only send media.pause and media.resume. This is media.play

	if strings.EqualFold(string(payload.Metadata.Type), string(models.MediaTypeShow)) {
		if !config.IsBeqTVEnabled() {
			return nil
		}
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before loading BEQ profile")
		return nil
	}

	if beqClient != nil {
		err = beqClient.LoadBeqProfile(searchRequest)
		if err != nil {
			log.Error("Error loading BEQ profile")
			return err
		}
		log.Info("BEQ profile loaded")
		// send notification of it loaded
		if config.IsBeqNotifyOnLoadEnabled() && config.IsHomeAssistantEnabled() {
			err := homeAssistantClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, searchRequest.Codec))
			if err != nil {
				log.Error("Error sending notification to HA")
				return err
			}
			log.Debug("sent notification to HA")
		}
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled while waiting for goroutines")
		return nil
	}

	innerWg.Wait()

	return nil
}

func HandlePause(ctx context.Context, payload *models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) error {
	log := logger.GetLoggerFromContext(ctx)
	// Check if context is already cancelled before starting lets say you play but then stop, this should stop processing
	if ctx.Err() != nil {
		log.Debug("Context already cancelled, stopping processing",
			"error", ctx.Err(),
			"func", "HandlePause",
		)
		return nil
	}
	log.Debug("handlepause searchRequest",
		slog.Any("searchRequest", searchRequest),
	)

	var err error

	// dont need to set skipActions here because it will only send media.pause and media.resume. This is media.play

	if strings.EqualFold(string(payload.Metadata.Type), string(models.MediaTypeShow)) {
		if !config.IsBeqTVEnabled() {
			return nil
		}
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before unloading BEQ profile")
		return nil
	}

	if beqClient != nil {
		err = beqClient.UnloadBeqProfile(searchRequest)
		if err != nil {
			return err
		}
		log.Info("BEQ profile unloaded")

		// send notification of it unloaded only if a profile is currently loaded
		if config.IsBeqNotifyOnUnLoadEnabled() && config.IsHomeAssistantEnabled() && beqClient.IsProfileLoaded() {
			err := homeAssistantClient.SendNotification("BEQ Profile Unloaded")
			if err != nil {
				log.Error("Error sending unload notification to HA")
				return err
			}
		}
	}

	return nil
}

func HandleStop(ctx context.Context, payload *models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) error {
	// same thing as stop pretty much
	return HandlePause(ctx, payload, wg, beqClient, homeAssistantClient, searchRequest)
}

func HandleResume(ctx context.Context, payload *models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) error {
	// TODO: support skip search for faster resume
	return HandlePlay(ctx, payload, wg, beqClient, homeAssistantClient, searchRequest)
}

func HandleScrobble(ctx context.Context, payload *models.Event) error {
	// Implement Plex-specific scrobble event handling
	return nil
}
