package common

// all common actions
import (
	"context"
	"log/slog"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
)

// IsAtmosodecPlaying checks if Atmos (mapped and normalized from the player -> eg plex codec name into BEQ name) is being decoded instead of multi ch in (plex bug I believe)
func IsAtmosCodecPlaying(codec, expectedCodec string) (bool, error) {
	if codec == expectedCodec {
		return true, nil
	}

	return false, nil
}

// readAttrAndWait is a generic func to read attr from HA
func readAttrAndWait(ctx context.Context, waitTime int, entType string, entName string, attrResp homeassistant.HAAttributeResponse, haClient *homeassistant.HomeAssistantClient) (bool, error) {
	var err error
	var isSignal bool
	log := logger.GetLoggerFromContext(ctx)

	// read attributes until its not nosignal
	for i := 0; i < waitTime; i++ {
		isSignal, err = haClient.ReadAttributes(entName, attrResp, entType)
		if err != nil {
			log.Error("Error reading attributes",
				slog.String("entity", entName),
				slog.String("error", err.Error()),
			)
			return false, err
		}
		log.Debug("Signal value",
			slog.String("entity", entName),
			slog.Bool("isSignal", isSignal),
		)
		if isSignal {
			log.Debug("HDMI sync complete")
			return isSignal, nil
		}

		// otherwise continue
		time.Sleep(200 * time.Millisecond)
	}

	return false, err
}

// common function for all supported players
// TODO: add generic plex/jf client
// func commonPlay(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, mediaClient Client, avrClient avr.AVRClient, payload interface{}, m *models.BeqSearchRequest, skipActions *bool, wg *sync.WaitGroup) {
// 	if payload == nil {
// 		log.Error("Payload is nil")
// 		return
// 	}
// 	if mediaClient == nil {
// 		log.Error("Media client is nil")
// 		return
// 	}
// 	// stop processing webhooks
// 	*skipActions = true
// 	var err error
// 	err = mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	go changeLight("off")
// 	// go changeAspect(client, payload, wg)
// 	go changeMasterVolume(m.MediaType)

// 	// if not using denoncodec, do this in background because we need to pause it anyway
// 	// TODO: verify config key
// 	if !config.GetBool("ezbeq.useAvrCodec") {
// 		wg.Add(1)
// 		// sets skipActions to false on completion
// 		go waitForHDMISync(wg, skipActions, haClient, mediaClient)
// 	}

// 	// always unload in case something is loaded from movie for tv
// 	err = beqClient.UnloadBeqProfile(m)
// 	if err != nil {
// 		log.Errorf("Error unloading beq on startup!! : %v", err)
// 		return
// 	}
// 	var year int
// 	var tmdb string
// 	var itemType string
// 	var edition string
// 	var title string

// 	// TODO: make vars which are generic containers for things like year, codec, etc
// 	// have to use any because go does not allow switch on generics but lets me just use an interface
// 	switch p := payload.(type) {
// 	case models.JellyfinWebhook:
// 		// TODO: move actions into here
// 		yearInt, err := strconv.Atoi(p.Year)
// 		if err != nil {
// 			log.Errorf("Error converting year to integer: %v", err)
// 			return
// 		}
// 		year = yearInt
// 		// TODO: JF title

// 		// TODO: make jellyfin client

// 	case models.PlexWebhookPayload:
// 		year = p.Metadata.Year
// 		// if its a show and you dont want beq enabled, exit
// 		if p.Metadata.Type == "episode" {
// 			if !config.GetBool("ezbeq.enableTvBeq") {
// 				return
// 			}
// 		}
// 		tmdb = mediaClient.GetPlexMovieDb(payload)
// 		title = p.Metadata.Title
// 	}

// 	m.Year = year
// 	m.TMDB = tmdb
// 	m.MediaType = itemType
// 	m.Edition = edition

// 	// get the codec
// 	if config.GetBool("ezbeq.useAvrCodec") {
// 		// TODO: map codec to map
// 		isexpected, err := isExpectedCodecPlaying(avrClient, m.Codec)
// 		if err != nil {
// 			log.Errorf("Error checking if expected codec is playing: %v", err)
// 			return
// 		}
// 		// TODO: rename key
// 		if config.GetBool("ezbeq.stopPlexIfMismatch") {
// 			if !isexpected {
// 				log.Debug("Stopping client because correct codec is not playing")
// 				err := PlaybackInterface("stop", mediaClient)
// 				if err != nil {
// 					log.Errorf("Error stopping client: %v", err)
// 				}
// 			}
// 		}
// 	} else {

// 		m.Codec, err = mediaClient.GetAudioCodec(payload)
// 		if err != nil {
// 			log.Errorf("error getting codec frin client, can't continue: %s", err)
// 			return
// 		}
// 	}

// 	log.Debugf("Found codec: %s", m.Codec)
// 	// TODO: check if beq is enabled

// 	err = beqClient.LoadBeqProfile(m)
// 	if err != nil {
// 		log.Errorf("Error loading beq profile: %v", err)
// 		return
// 	}
// 	log.Info("BEQ profile loaded")

// 	// send notification of it loaded
// 	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
// 		err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", title, year, m.Codec), config.GetString("ezbeq.notifyEndpointName"))
// 		if err != nil {
// 			log.Error()
// 		}
// 	}

// 	log.Debug("Waiting for goroutines")
// 	wg.Wait()
// 	log.Debug("goroutines complete")
// }
