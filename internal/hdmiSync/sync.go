package hdmisync

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/models"
)

// waitForHDMISync will pause until the source reports HDMI sync is complete
func WaitForHDMISync(ctx context.Context, wg *sync.WaitGroup, skipActions *bool, haClient *homeassistant.HomeAssistantClient, mediaClient mediaplayer.MediaAPIClient) {
	// if called and disabled, skip
	// stop processing webhooks because if we call pause, that will fire another one and then we get into a loop
	*skipActions = true
	log := logger.GetLoggerFromContext(ctx)
	if !config.IsHDMISyncEnabled() {
		*skipActions = false
		wg.Done()
		return
	}
	log.Debug("Running HDMI sync wait")

	defer func() {
		// play item no matter what happens
		err := mediaClient.DoPlaybackAction(ctx, models.ActionPlay)
		if err != nil {
			log.Error("Error playing client", slog.String("error", err.Error()))
			return
		}

		// continue processing webhooks when done/
		// if webhook is delayed, resume will get processed so wait
		time.Sleep(10 * time.Second)
		*skipActions = false
		wg.Done()
	}()

	signalSource := config.GetHDMISyncSource()
	var err error
	var signal bool

	// pause client
	log.Debug("Pausing client")
	err = mediaClient.DoPlaybackAction(ctx, models.ActionPause)
	if err != nil {
		log.Error("Error pausing client", slog.String("error", err.Error()))
		return
	}

	// check signal source
	switch signalSource {
	case "envy":
		log.Debug("Using envy for hdmi sync")
		// read envy attributes until its not nosignal
		envyName := config.GetHDMISyncEnvyName() // TODO: this should support any device
		// remove remote. if present
		if strings.Contains(envyName, "remote") {
			envyName = strings.ReplaceAll(envyName, "remote.", "")
		}
		signal, err = haClient.ReadAttrAndWait(ctx, 60, "remote", envyName) // TODO: implement
		// will break out here
	case "time":
		seconds := config.GetHDMISyncSeconds()
		log.Debug("Using time for hdmi sync", slog.String("seconds", seconds))
		sec, err := strconv.Atoi(seconds)
		if err != nil {
			log.Error("waitforHDMIsync enabled but no valid source provided. Make sure you have 'time' set as a plain number",
				slog.String("source", signalSource),
				slog.String("error", err.Error()),
			)
			return
		}
		time.Sleep(time.Duration(sec) * time.Second)
		return
	case "jvc":
		// read jvc attributes until its not nosignal
		log.Warn("jvc HDMI sync is not implemented")
	case "sensor":
		log.Warn("sensor HDMI sync is not implemented")
	default:
		log.Warn("No valid source provided for hdmi sync")
	}

	log.Debug("HDMI sync complete", slog.Bool("signal", signal))

	if err != nil {
		log.Error("Error getting HDMI signal", slog.String("error", err.Error()))
	}
}
