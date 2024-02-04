package common

// common actions
import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"

	// "github.com/iloveicedgreentea/go-plex/internal/plex"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

// TODO: add optional thing for it to tell plex to stop playing and start the stream again
// TODO: make test
// func isExpectedCodecPlayingold(c *avr.DenonClient, p *plex.PlexClient, uuid string, denonCodec string) (string, bool) {
// 	plexPlaying, err := p.GetCodecFromSession(uuid)
// 	if err != nil {
// 		log.Errorf("Error getting plex audio stream: %s", err)
// 		return "", false
// 	}

// 	// compare the two
// 	return plexPlaying, mapAvrToBeq(denonCodec) != plex.MapPlexToBeqAudioCodec(plexPlaying, "")
// 	// if enabled, stop playing

// 	log.Error("Expected codec is not playing! Please check your AVR and Plex settings!")
// 	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
// 		err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %s", m.Codec, expectedCodec), config.GetString("ezbeq.notifyEndpointName"))
// 		if err != nil {
// 			log.Error(err)
// 		}
// 	}
// }

// isExpectedCodecPlaying checks if avr is playing expectedCodec (mapped and normalized string)
func isExpectedCodecPlaying(avrClient avr.AVRClient, expectedCodec string) (bool, error) {
	// get the codec from avr
	codec, err := avrClient.GetCodec()
	if err != nil {
		log.Errorf("error getting codec from denon, can't continue: %s", err)
		return false, err
	}

	if codec != expectedCodec {
		return false, nil
	}

	return true, nil
}

// common function for all supported players
// TODO: add generic plex/jf client
func commonPlay(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, mediaClient Client, avrClient avr.AVRClient, payload interface{}, m *models.SearchRequest, skipActions *bool, wg *sync.WaitGroup) {
	if payload == nil {
		log.Error("Payload is nil")
		return
	}
	if mediaClient == nil {
		log.Error("Media client is nil")
		return
	}
	// stop processing webhooks
	*skipActions = true
	var err error
	err = mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
	if err != nil {
		log.Error(err)
	}
	go changeLight("off")
	// go changeAspect(client, payload, wg)
	go changeMasterVolume(m.MediaType)

	// if not using denoncodec, do this in background because we need to pause it anyway
	// TODO: verify config key
	if !config.GetBool("ezbeq.useAvrCodec") {
		wg.Add(1)
		// sets skipActions to false on completion
		go waitForHDMISync(wg, skipActions, haClient, mediaClient)
	}

	// always unload in case something is loaded from movie for tv
	err = beqClient.UnloadBeqProfile(m)
	if err != nil {
		log.Errorf("Error unloading beq on startup!! : %v", err)
		return
	}
	var year int
	var tmdb string
	var itemType string
	var edition string
	var title string

	// TODO: make vars which are generic containers for things like year, codec, etc
	// have to use any because go does not allow switch on generics but lets me just use an interface
	switch p := payload.(type) {
	case models.JellyfinWebhook:
		// TODO: move actions into here
		yearInt, err := strconv.Atoi(p.Year)
		if err != nil {
			log.Errorf("Error converting year to integer: %v", err)
			return
		}
		year = yearInt
		// TODO: JF title

		// TODO: make jellyfin client

	case models.PlexWebhookPayload:
		year = p.Metadata.Year
		// if its a show and you dont want beq enabled, exit
		if p.Metadata.Type == "episode" {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}
		tmdb = mediaClient.GetPlexMovieDb(payload)
		title = p.Metadata.Title
	}

	m.Year = year
	m.TMDB = tmdb
	m.MediaType = itemType
	m.Edition = edition

	// get the codec
	if config.GetBool("ezbeq.useAvrCodec") {
		// TODO: map codec to map
		isexpected, err := isExpectedCodecPlaying(avrClient, m.Codec)
		if err != nil {
			log.Errorf("Error checking if expected codec is playing: %v", err)
			return
		}
		if config.GetBool("ezbeq.stopPlexIfMismatch") {
			if !isexpected {
				log.Debug("Stopping plex because codec is not playing")
				err := PlaybackInterface("stop", mediaClient)
				if err != nil {
					log.Errorf("Error stopping plex: %v", err)
				}
			}
		}
	} else {
		
		m.Codec, err = mediaClient.GetAudioCodec(payload)
		if err != nil {
			log.Errorf("error getting codec from plex, can't continue: %s", err)
			return
		}
	}

	log.Debugf("Found codec: %s", m.Codec)
	// TODO: check if beq is enabled

	err = beqClient.LoadBeqProfile(m)
	if err != nil {
		log.Errorf("Error loading beq profile: %v", err)
		return
	}
	log.Info("BEQ profile loaded")

	// send notification of it loaded
	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
		err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", title, year, m.Codec), config.GetString("ezbeq.notifyEndpointName"))
		if err != nil {
			log.Error()
		}
	}

	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("goroutines complete")
}

// trigger HA for MV change per type
func changeMasterVolume(mediaType string) {
	if config.GetBool("homeAssistant.triggerAvrMasterVolumeChangeOnEvent") && config.GetBool("homeAssistant.enabled") {
		log.Debug("changeMasterVolume: Changing volume")
		err := mqtt.Publish([]byte(fmt.Sprintf("{\"type\":\"%s\"}", mediaType)), config.GetString("mqtt.topicVolume"))
		if err != nil {
			log.Error()
		}
	}
}

// trigger HA for light change given entity and desired state
func changeLight(state string) {
	if config.GetBool("homeAssistant.triggerLightsOnEvent") && config.GetBool("homeAssistant.enabled") {
		log.Debug("changeLight: Changing light")
		err := mqtt.Publish([]byte(fmt.Sprintf("{\"state\":\"%s\"}", state)), config.GetString("mqtt.topicLights"))
		if err != nil {
			log.Errorf("Error changing light: %v", err)
		}
	}
}

// TODO: test this
// waitForHDMISync will wait until the envy reports a signal to assume hdmi sync. No API to do this with denon afaik
func waitForHDMISync(wg *sync.WaitGroup, skipActions *bool, haClient *homeassistant.HomeAssistantClient, mediaClient Client) {
	if !config.GetBool("signal.enabled") {
		*skipActions = false
		wg.Done()
		return
	}

	log.Debug("Running HDMI sync wait")
	defer func() {
		// play item no matter what
		err := PlaybackInterface("play", mediaClient)
		if err != nil {
			log.Errorf("Error playing plex: %v", err)
			return
		}

		// continue processing webhooks when done
		*skipActions = false
		wg.Done()
	}()

	signalSource := config.GetString("signal.source")
	var err error
	var signal bool

	// pause plex
	log.Debug("pausing plex")
	err = PlaybackInterface("pause", mediaClient)
	if err != nil {
		log.Errorf("Error pausing plex: %v", err)
		return
	}

	switch signalSource {
	case "envy":
		// read envy attributes until its not nosignal
		signal, err = readAttrAndWait(30, "remote", &models.HAEnvyResponse{}, haClient)
	case "jvc":
		// read jvc attributes until its not nosignal
		signal, err = readAttrAndWait(30, "remote", &models.HAjvcResponse{}, haClient)
	case "sensor":
		signal, err = readAttrAndWait(30, "binary_sensor", &models.HABinaryResponse{}, haClient)
	default:
		// TODO: maybe use a 15 sec delay?
		log.Debug("using seconds for hdmi sync")
		sec, err := strconv.Atoi(signalSource)
		if err != nil {
			log.Errorf("waitforHDMIsync enabled but no valid source provided: %v -- %v", signalSource, err)
			return
		}
		time.Sleep(time.Duration(sec) * time.Second)

	}

	log.Debugf("HDMI Signal value is %v", signal)
	if err != nil {
		log.Errorf("error getting HDMI signal: %v", err)
	}

}

// readAttrAndWait is a generic func to read attr from HA
func readAttrAndWait(waitTime int, entType string, attrResp homeassistant.HAAttributeResponse, haClient *homeassistant.HomeAssistantClient) (bool, error) {
	var err error
	var isSignal bool

	for i := 0; i < waitTime; i++ {
		isSignal, err = haClient.ReadAttributes(haClient.EntityName, attrResp, entType)
		if isSignal {
			log.Debug("HDMI sync complete")
			return isSignal, nil
		}
		if err != nil {
			log.Errorf("Error reading envy attributes: %v", err)
			return false, err
		}
		// otherwise continue
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Errorf("Error reading envy attributes: %v", err)
		return false, err
	}

	return false, err

}
