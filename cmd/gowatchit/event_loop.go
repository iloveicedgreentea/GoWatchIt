package main

import (
	"context"
	"log/slog"
	"sync"

	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/models"
)

func eventHandler(ctx context.Context, c <-chan models.Event, beqClient *ezbeq.BeqClient, client mediaplayer.MediaAPIClient, homeAssistantClient *homeassistant.HomeAssistantClient) {
	// get payload from chan in a loop
	for payload := range c {
		ctx, cancel := context.WithCancel(ctx)
		var wg *sync.WaitGroup
		var searchRequest *models.BeqSearchRequest
		log := logger.GetLoggerFromContext(ctx)

		// get edition
		// TODO: config if editions matter
		edition, err := client.GetEdition(ctx, payload)
		if err != nil {
			// TODO: handle error more? is it fatal?
			log.Error("Error getting edition",
				slog.Any("error", err),
			)
		}
		// get codec
		// TODO: avr or session confic
		codec, err := client.GetAudioCodec(ctx, payload)
		if err != nil {
			// TODO: handle error more? is it fatal?
			log.Error("Error getting codec",
				slog.Any("error", err),
			)
		}
		// TODO: get TMDB
		// TODO: config for BEQ enabled
		searchRequest = beqClient.NewRequest(ctx, false, payload.Metadata.Year, payload.Metadata.Type, edition, payload.Metadata.TMDB, codec)
		
		// TODO: do this async first, wait to load until this is done
		
		if err = beqClient.UnloadBeqProfile(searchRequest); err != nil {
			log.Error("Error unloading BEQ during play",
				slog.Any("error", err),
			)
		}

		wg.Wait()

		// TODO: route event
		// TODO: handle cancelations
		go eventRouter(ctx, cancel, payload, wg, beqClient, client, homeAssistantClient, searchRequest)
	}
}

func eventRouter(ctx context.Context, cancel context.CancelFunc, event models.Event, wg *sync.WaitGroup, beqClient *ezbeq.BeqClient, client mediaplayer.MediaAPIClient, homeAssistantClient *homeassistant.HomeAssistantClient, searchRequest *models.BeqSearchRequest) {
	switch event.Action {
	case models.ActionPlay:
		// TODO: remove functions in main loop from this func like codec
		mediaplayer.HandlePlay(ctx, cancel, event, wg, beqClient, client, homeAssistantClient, searchRequest)
	}
}
