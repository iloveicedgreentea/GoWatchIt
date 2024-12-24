package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/plex"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/models"
)

func initializeMediaClient(event *models.Event) (mediaplayer.MediaAPIClient, error) {
	// For Plex
	if event.EventType == models.EventTypePlex && config.IsPlexEnabled() {
		return plex.NewClient(
			config.GetPlexScheme(),
			config.GetPlexUrl(),
			config.GetPlexPort(),
		)
	}

	// For Jellyfin
	if event.EventType == models.EventTypeJellyfin && config.IsJellyfinEnabled() {
		// TODO: Implement Jellyfin client
		return nil, nil
	}

	// TODO: home assistant
	if event.EventType == models.EventTypeHomeAssistant && config.IsHomeAssistantEnabled() {
		return nil, nil
	}

	return nil, fmt.Errorf("unsupported server type for Event: %s", event.ServerUUID)
}

func getClient(payload *models.Event) (client mediaplayer.MediaAPIClient, err error) {
	if payload.Client == nil {
		var err error
		payload.Client, err = initializeMediaClient(payload)
		if err != nil {
			return nil, fmt.Errorf("error initializing media client: %v", err)
		}
	}
	client, ok := payload.Client.(mediaplayer.MediaAPIClient)
	if !ok {
		return nil, fmt.Errorf("error checking client is MediaAPIClient : %v - payload client is %v", payload, payload.Client)
	}

	return client, nil
}

func eventHandler(ctx context.Context, c <-chan models.Event, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient) {
	log := logger.GetLoggerFromContext(ctx)
	for payload := range c {
		log.Debug("Received event", slog.Any("event", payload))
		client, err := getClient(&payload)
		if err != nil {
			log.Error("Error getting client for payload", slog.Any("error", err))
			continue
		}

		// Create a new context for each event, but don't cancel it immediately
		eventCtx, eventCancel := context.WithCancel(ctx)

		// Use a WaitGroup to ensure all goroutines complete before moving to the next event
		var wg sync.WaitGroup

		// Process the event in a separate goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer eventCancel() // Ensure the context is canceled when this goroutine exits

			log := logger.GetLoggerFromContext(eventCtx)

			// Get edition
			// TODO: config isEditionEnabled to ignore edition matching
			edition, err := client.GetEdition(eventCtx, &payload)
			if err != nil {
				log.Error("Error getting edition", slog.Any("error", err))
				return
			}

			// get codec
			// TODO: avr or session confic
			// slower but more accurate especially with atmos
			// TODO: implement avr stuff
			// if useAvrCodec {
			// p.SearchRequest.Codec, err = checkAvrCodec(client, haClient, avrClient, payload, data)
			// // if it failed, get codec data from client
			// if err != nil {
			// log.Warnf("error getting codec from AVR, falling back to client: %s", err)
			// m.Codec, err = client.GetAudioCodec(data)
			// if err != nil {
			// log.Errorf("error getting codec from plex, can't continue: %s", err)
			// return
			// }
			// }
			// } else {
			// log.Debug("Using plex to get codec")
			// // TODO: try session data then fallback to lookup
			// m.Codec, err = client.GetAudioCodec(data)
			// if err != nil {
			// log.Errorf("error getting codec from plex, can't continue: %s", err)
			// return
			// }
			// }
			// TODO: use config to get which kind of codec but insice this func
			codec, err := client.GetAudioCodec(eventCtx, &payload)
			if err != nil {
				log.Error("Error getting codec", slog.Any("error", err))
				return
			}
			var searchRequest *models.BeqSearchRequest
			if beqClient != nil {
				// Create BEQ search request
				searchRequest = beqClient.NewRequest(eventCtx, false, payload.Metadata.Year, payload.Metadata.Type, edition, payload.Metadata.TMDB, codec)
				if searchRequest == nil {
					log.Error("Error creating BEQ search request. Unable to proceed with BEQ operations. Check your config.")
					return
				}
				log.Debug("created new searchRequest", slog.Any("searchRequest", searchRequest))

				// Unload BEQ profile
				if err := beqClient.UnloadBeqProfile(searchRequest); err != nil {
					log.Error("Error unloading BEQ during play", slog.Any("error", err))
				}
			}

			// Route event
			eventRouter(eventCtx, eventCancel, &payload, &wg, beqClient, homeAssistantClient, searchRequest)
		}()

		// Wait for all goroutines to complete before processing the next event
		wg.Wait()
	}
}

func eventRouter(ctx context.Context, cancel context.CancelFunc, event *models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) {
	if event.Action == models.ActionPlay {
		err := mediaplayer.HandlePlay(ctx, cancel, event, wg, beqClient, homeAssistantClient, searchRequest)
		if err != nil {
			logger.GetLogger().Error("Error handling play event", slog.Any("error", err))
		}
	}
	// TODO: add other actions like pause
}
