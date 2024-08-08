package mediaplayer

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
	"github.com/iloveicedgreentea/go-plex/models"
)

// Implement MediaEventHandler methods
func HandlePlay(ctx context.Context, cancel context.CancelFunc, payload models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, client MediaAPIClient, homeAssistantClient *homeassistant.HomeAssistantClient) error {
	log := logger.GetLoggerFromContext(ctx)
	// Check if context is already cancelled before starting lets say you play but then stop, this should stop processing
	if ctx.Err() != nil {
		log.Debug("Context already cancelled, stopping processing",
			"error", ctx.Err(),
			"func", "HandlePlay",
		)
		return nil
	}

	var err error
	go func() {
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before lights and volume change")
			return
		}
		common.ChangeLight(models.ActionOff)
		common.ChangeMasterVolume(payload.Metadata.Type)
	}()

	// TODO: check this
	// Perform HDMI sync
	if !strings.EqualFold(string(payload.Metadata.Type), string(models.MediaTypeMovie)) && config.GetString("signal.source") == "time" {
		log.Debug("skipping sync for non-movie type and time source")
	} else {
		// wg.Add(1)
		go func() {
			if ctx.Err() != nil {
				log.Debug("mediaPlay was cancelled before hdmi sync")
				return // Exit early if context is cancelled
			}

			// optimistically try to hdmi sync. Will return if disabled
			// TODO: implement this
			// common.WaitForHDMISync(wg, skipActions, haClient, client)
		}()
	}

	// dont need to set skipActions here because it will only send media.pause and media.resume. This is media.play

	go func() {
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before publishing playing status")
			return
		}
		// TODO: make a send playing topic function isntead of passing in topic
		if err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true"); err != nil {
			log.Error("Error publishing playing status: ", err)
		}
	}()

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before unloading BEQ profile")
		return nil
	}
	// TODO: create search request
	searchReq := beqClient.NewRequest()
	if err = beqClient.UnloadBeqProfile(searchReq); err != nil {
		log.Error("Error unloading BEQ during play",
			slog.Any("error", err),
		)
	}

	select {
	case <-ctx.Done():
		log.Error("mediaPlay cancelled before unloading BEQ profile")
		return nil
	default:
		log.Debug("Using plex to get codec")
		// TODO: try session data then fallback to lookup
		searchReq.Codec, err = client.GetAudioCodec(ctx, models.DataMediaContainer{PlexPayload: &data})
		if err != nil {
			return fmt.Errorf("error getting codec from plex, can't continue: %s", err)
		}
		// slower but more accurate especially with atmos
		// TODO: implement avr stuff
		// if useAvrCodec {
		// 	p.SearchRequest.Codec, err = checkAvrCodec(client, haClient, avrClient, payload, data)
		// 	// if it failed, get codec data from client
		// 	if err != nil {
		// 		log.Warnf("error getting codec from AVR, falling back to client: %s", err)
		// 		m.Codec, err = client.GetAudioCodec(data)
		// 		if err != nil {
		// 			log.Errorf("error getting codec from plex, can't continue: %s", err)
		// 			return
		// 		}
		// 	}
		// } else {
		// 	log.Debug("Using plex to get codec")
		// 	// TODO: try session data then fallback to lookup
		// 	m.Codec, err = client.GetAudioCodec(data)
		// 	if err != nil {
		// 		log.Errorf("error getting codec from plex, can't continue: %s", err)
		// 		return
		// 	}
		// }
		log.Debug("Found codec: %s", searchReq.Codec)
		// if its a show and you dont want beq enabled, exit
		if strings.EqualFold(string(payload.Metadata.Type), string(models.MediaTypeShow)) {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return nil
			}
		}

		// TODO: get TMDB from client
		searchReq.TMDB = getPlexMovieDb(payload)
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before loading BEQ profile")
		return nil
	}
	// TODO: pass in object
	err = beqClient.LoadBeqProfile(searchReq)
	if err != nil {
		return err
	}
	log.Info("BEQ profile loaded")
	// send notification of it loaded
	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
		err := homeAssistantClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, searchReq.Codec))
		if err != nil {
			return err
		}
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled while waiting for goroutines")
		return nil
	}

	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("Goroutines complete")

	return nil
}

func HandlePause(ctx context.Context, cancel context.CancelFunc, payload models.Event) error {
	// Implement Plex-specific pause event handling
	return nil
}

func HandleStop(ctx context.Context, cancel context.CancelFunc, payload models.Event) error {
	// Implement Plex-specific stop event handling
	return nil
}

func HandleResume(ctx context.Context, cancel context.CancelFunc, payload models.Event) error {
	// Implement Plex-specific resume event handling
	return nil
}

func HandleScrobble(ctx context.Context, payload models.Event) error {
	// Implement Plex-specific scrobble event handling
	return nil
}
