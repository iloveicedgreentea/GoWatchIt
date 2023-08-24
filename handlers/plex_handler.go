package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/iloveicedgreentea/go-plex/denon"
	"github.com/iloveicedgreentea/go-plex/ezbeq"
	"github.com/iloveicedgreentea/go-plex/homeassistant"
	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/iloveicedgreentea/go-plex/mqtt"
	"github.com/iloveicedgreentea/go-plex/plex"
	"github.com/spf13/viper"
)

const showItemTitle = "episode"
const movieItemTitle = "movie"

var log = logger.GetLogger()

func decodeWebhook(payload []string) (models.PlexWebhookPayload, int, error) {
	var pwhPayload models.PlexWebhookPayload

	err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// unmarshall error
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return pwhPayload, http.StatusBadRequest, errors.New(msg)

		default:
			return pwhPayload, http.StatusInternalServerError, err
		}
	}

	log.Debugf("decodeWebhook: Received event: %s", pwhPayload.Event)
	return pwhPayload, 0, nil
}

// Sends the payload to the channel for background processing
func ProcessWebhook(plexChan chan<- models.PlexWebhookPayload, vip *viper.Viper) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(0); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("ProcessWebhook: Received valid multipart form")
		payload, hasPayload := r.MultipartForm.Value["payload"]
		if hasPayload {
			log.Debug("decoding payload")
			decodedPayload, statusCode, err := decodeWebhook(payload)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), statusCode)
				return
			}
			clientUUID := decodedPayload.Player.UUID
			log.Infof("Got a request from UUID: %s", clientUUID)

			log.Debugf("ProcessWebhook:  Media type is: %s", decodedPayload.Metadata.Type)
			log.Debugf("ProcessWebhook:  Media title is: %s", decodedPayload.Metadata.Title)

			// check filter for user if not blank
			userID := vip.GetString("plex.ownerNameFilter")
			// only respond to events on a particular account if you share servers and only for movies and shows
			if userID == "" || decodedPayload.Account.Title == userID {
				if decodedPayload.Metadata.Type == movieItemTitle || decodedPayload.Metadata.Type == showItemTitle {
					log.Debugf("Current length of plexChan: %d", len(plexChan))
					select {
					case plexChan <- decodedPayload:
						// send succeeded
					case <-time.After(time.Second * 3):
						log.Error("Send on plexChan timed out")
						http.Error(w, err.Error(), http.StatusTooManyRequests)
						return
					}
					log.Debugf("Added length of plexChan: %d", len(plexChan))
				}
			} else {
				log.Debugf("userID '%s' does not match filter", userID)
			}
		} else {
			log.Error("No payload found in request")
			http.Error(w, "No payload found in request", http.StatusBadRequest)
			return
		}
	}

	return http.HandlerFunc(fn)
}

// does plex send stop if you exit with back button? - Yes, with X for mobile player as well
func mediaStop(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go changeLight(vip, "on", wg)

	if vip.GetBool("ezbeq.enabled") {
		err := beqClient.UnloadBeqProfile(m)
		if err != nil {
			log.Error(err)
			if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err), vip.GetString("ezbeq.notifyEndpointName"))
				if err != nil {
					log.Error()
				}
			}
		}
		log.Info("BEQ profile unloaded")
	}
}

// pause only happens with literally pausing
func mediaPause(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest, skipActions *bool) {
	if !*skipActions {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go changeLight(vip, "on", wg)

		if vip.GetBool("ezbeq.enabled") {
			err := beqClient.UnloadBeqProfile(m)
			if err != nil {
				log.Error(err)
				if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
					err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err), vip.GetString("ezbeq.notifyEndpointName"))
					if err != nil {
						log.Error()
					}
				}
			}
			log.Info("BEQ profile unloaded")
		}
	}
}

// waitForHDMISync will wait until the envy reports a signal to assume hdmi sync. No API to do this with denon afaik
func waitForHDMISync(wg *sync.WaitGroup, skipActions *bool, haClient *homeassistant.HomeAssistantClient, PlexClient *plex.PlexClient, vip *viper.Viper) {
	if !vip.GetBool("ezbeq.waitforHDMIsync") {
		wg.Done()
		return
	}

	log.Debug("Running HDMI sync wait")
	defer func() {
		// play item no matter what
		err := PlexClient.PlayPlex()
		if err != nil {
			log.Errorf("Error playing plex: %v", err)
			return
		}

		// continue processing webhooks when done
		*skipActions = false
		wg.Done()
	}()

	// TODO: Explore Denon SD? -> SDNO, SDHMDI

	var err error
	var isSignal bool

	// pause plex
	log.Debug("pausing plex")
	err = PlexClient.PausePlex()
	if err != nil {
		log.Errorf("Error pausing plex: %v", err)
		return
	}

	// read envy attributes until its not nosignal
	// wait up to 30s to readEnvyAttributes
	for i := 0; i < 30; i++ {
		isSignal, err = haClient.ReadEnvyAttributes()
		if isSignal {
			log.Debug("HDMI sync complete")
			return
		}
		if err != nil {
			log.Errorf("Error reading envy attributes: %v", err)
			return
		}
		// otherwise continue
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Errorf("Error reading envy attributes: %v", err)
		return
	}

}

// play is both the "resume" button and play
func mediaPlay(client *plex.PlexClient, vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, denonClient *denon.DenonClient, payload models.PlexWebhookPayload, m *models.SearchRequest, useDenonCodec bool, data models.MediaContainer, skipActions *bool) {
	wg := &sync.WaitGroup{}

	// stop processing webhooks
	*skipActions = true

	wg.Add(3)
	go changeLight(vip, "off", wg)
	go changeAspect(client, payload, vip, wg)
	go changeMasterVolume(vip, m.MediaType, wg)

	// if not using denoncodec, do this in background
	if !useDenonCodec {
		wg.Add(1)
		// TODO: get this working
		// sets skipActions to false on completion
		go waitForHDMISync(wg, skipActions, haClient, client, vip)
	}

	if vip.GetBool("ezbeq.enabled") {
		// always unload in case something is loaded from movie for tv
		err := beqClient.UnloadBeqProfile(m)
		if err != nil {
			log.Errorf("Error unloading beq on startup!! : %v", err)
			return
		}

		// slower but more accurate
		if useDenonCodec {
			// wait for sync
			wg.Add(1)
			waitForHDMISync(wg, skipActions, haClient, client, vip)
			// denon needs time to show mutli ch in as atmos
			// TODO: test this
			time.Sleep(5 * time.Second)

			// get the codec from avr
			m.Codec, err = denonClient.GetCodec()
			if err != nil {
				log.Errorf("error getting codec from denon, can't continue: %s", err)
				return
			}

			// check if the expected codec is playing
			expectedCodec, isExpectedPlaying := isExpectedCodecPlaying(denonClient, client, payload.Player.UUID, m.Codec)
			if !isExpectedPlaying {
				// if enabled, stop playing
				if vip.GetBool("ezbeq.stopPlexIfMismatch") {
					log.Debug("Stopping plex because codec is not playing")
					err := client.StopPlex()
					if err != nil {
						log.Errorf("Error stopping plex: %v", err)
					}
				}

				log.Error("Expected codec is not playing! Please check your AVR and Plex settings!")
				if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
					err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %s", m.Codec, expectedCodec), vip.GetString("ezbeq.notifyEndpointName"))
					if err != nil {
						log.Error(err)
					}
				}
			}

		} else {
			m.Codec, err = client.GetAudioCodec(data)
			if err != nil {
				log.Errorf("error getting codec from plex, can't continue: %s", err)
				return
			}
		}

		log.Debugf("Found codec: %s", m.Codec)

		// if its a show and you dont want beq enabled, exit
		if payload.Metadata.Type == showItemTitle {
			if !vip.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}

		m.TMDB = getPlexMovieDb(payload)
		err = beqClient.LoadBeqProfile(m)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("BEQ profile loaded")

		// send notification of it loaded
		if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
			err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, m.Codec), vip.GetString("ezbeq.notifyEndpointName"))
			if err != nil {
				log.Error()
			}
		}
	}
	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("goroutines complete")
}

// resume is only after pausing as long as the media item is still active
func mediaResume(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest, skipActions *bool) {
	if !*skipActions {
		wg := &sync.WaitGroup{}
		// mediaType string, codec string, edition string
		// trigger lights
		wg.Add(1)
		go changeLight(vip, "off", wg)
		// Changing on resume is disabled because its annoying if you changed it since playing
		// go changeMasterVolume(vip, mediaType)

		// allow skipping search to save time
		if vip.GetBool("ezbeq.enabled") {
			// always unload in case something is loaded from movie for tv
			err := beqClient.UnloadBeqProfile(m)
			if err != nil {
				log.Errorf("Error on startup - unloading beq %v", err)
			}
			if payload.Metadata.Type == showItemTitle {
				if !vip.GetBool("ezbeq.enableTvBeq") {
					return
				}
			}
			// get the tmdb id to match with ezbeq catalog
			m.TMDB = getPlexMovieDb(payload)
			// load beq with cache
			err = beqClient.LoadBeqProfile(m)
			if err != nil {
				log.Error(err)
				return
			}
			log.Info("BEQ profile loaded")

			// send notification of it loaded
			if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, m.Codec), vip.GetString("ezbeq.notifyEndpointName"))
				if err != nil {
					log.Error()
				}
			}
		}
	}
}

func mediaScrobble(vip *viper.Viper) {
	// trigger lights
	log.Debug("Scrobble received. Not doing anything")
	// go changeLight(vip, "on")
}

// getEditionName tries to extract the edition from plex or file name. Assumes you have well named files
// Returned types, Unrated, Ultimate, Theatrical, Extended, Director, Criterion
func getEditionName(data models.MediaContainer) string {
	edition := data.Video.EditionTitle
	fileName := strings.ToLower(data.Video.Media.Part.File)

	// if there is an edition from plex metadata, use it
	if edition != "" {
		return edition
	}

	switch {
	case strings.Contains(fileName, "extended"):
		return "Extended"
	case strings.Contains(fileName, "unrated"):
		return "Unrated"
	case strings.Contains(fileName, "theatrical"):
		return "Theatrical"
	case strings.Contains(fileName, "ultimate"):
		return "Ultimate"
	case strings.Contains(fileName, "director"):
		return "Director"
	case strings.Contains(fileName, "criterion"):
		return "Criterion"
	default:
		return ""
	}
}

// based on event type, determine what to do
func eventRouter(plexClient *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, denonClient *denon.DenonClient, useDenonCodec bool, payload models.PlexWebhookPayload, vip *viper.Viper, model *models.SearchRequest, skipActions *bool) {
	// perform function via worker

	clientUUID := payload.Player.UUID
	// ensure the client matches so it doesnt trigger from unwanted clients

	if vip.GetString("plex.deviceUUIDFilter") != clientUUID || vip.GetString("plex.deviceUUIDFilter") == "" {
		log.Debug("Client UUID does not match enabled filter")
		return
	}

	plexClient.MachineID = clientUUID
	plexClient.MediaType = payload.Metadata.Type
	log.Infof("Processing media type: %s", payload.Metadata.Type)

	var err error
	var data models.MediaContainer
	var editionName string

	if vip.GetBool("ezbeq.enabled") {
		// make a call to plex to get the data based on key
		data, err = plexClient.GetMediaData(payload.Metadata.Key)
		if err != nil {
			if strings.Contains(err.Error(), "but have <html>") {
				log.Error("Error authenticating with plex. Please check your IP whitelist")
			} else {
				log.Errorf("Error getting media data from plex: %s", err)
			}
			return
		} else {
			// get the edition name
			editionName = getEditionName(data)
			log.Debugf("Event Router: Found edition: %s", editionName)
		}
	}

	log.Debugf("Event Router: Got media type of: %s ", payload.Metadata.Type)

	model.Year = payload.Metadata.Year
	model.Edition = editionName
	// this should be updated with every event
	model.EntryID = beqClient.CurrentProfile
	model.MVAdjust = beqClient.MasterVolume

	log.Debugf("Event Router: Using search model: %#v", model)
	switch payload.Event {
	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
	// play means a new file was started
	case "media.play":
		log.Debug("Event Router: media.play received")
		mediaPlay(plexClient, vip, beqClient, haClient, denonClient, payload, model, useDenonCodec, data, skipActions)
	case "media.stop":
		log.Debug("Event Router: media.stop received")
		mediaStop(vip, beqClient, haClient, payload, model)
	case "media.pause":
		log.Debug("Event Router: media.pause received")
		mediaPause(vip, beqClient, haClient, payload, model, skipActions)
	// Pressing the 'resume' button in plex is media.play
	case "media.resume":
		log.Debug("Event Router: media.resume received")
		mediaResume(vip, beqClient, haClient, payload, model, skipActions)
	case "media.scrobble":
		log.Debug("Scrobble received")
		mediaScrobble(vip)
	default:
		log.Debugf("Received unsupported event: %s", payload.Event)
	}
}

// get the imdb ID from plex metadata
func getPlexImdbID(payload models.PlexWebhookPayload) string {
	// try to get IMDB title from plex to save time
	for _, model := range payload.Metadata.GUID0 {
		if strings.Contains(model.ID, "imdb") {
			log.Debugf("Got title ID from plex - %s", model.ID)
			return strings.Split(model.ID, "imdb://")[1]
		}
	}

	return ""
}

// get the tmdb ID from plex metadata
func getPlexMovieDb(payload models.PlexWebhookPayload) string {
	// try to get IMDB title from plex to save time
	for _, model := range payload.Metadata.GUID0 {
		if strings.Contains(model.ID, "tmdb") {
			log.Debugf("getPlexMovieDb: Got tmdb ID from plex - %s", model.ID)
			return strings.Split(model.ID, "tmdb://")[1]
		}
	}
	log.Error("TMDB id not found in Plex. ezBEQ will not work. Please check your metadata for this title!")
	return ""
}

// will change aspect ratio
func changeAspect(client *plex.PlexClient, payload models.PlexWebhookPayload, vip *viper.Viper, wg *sync.WaitGroup) {
	defer wg.Done()
	if vip.GetBool("homeAssistant.triggerAspectRatioChangeOnEvent") && vip.GetBool("homeAssistant.enabled") {

		// if madvr enabled, only send a trigger via mqtt
		// This needs to be triggered by MQTT in automation. This sends a pulse. The automation reads envy aspect ratio 5 sec later
		if vip.GetBool("madvr.enabled") {
			log.Debug("Using madvr for aspect")
			topic := vip.GetString("mqtt.topicAspectratioMadVrOnly")

			// trigger automation
			err := mqtt.Publish(vip, []byte(""), topic)
			if err != nil {
				log.Error()
			}

			return
		} else {
			log.Debug("changeAspect: Changing aspect ratio")

			// get the imdb title ID
			imdbID := getPlexImdbID(payload)

			// lookup aspect based on imdb technical info
			aspect, err := client.GetAspectRatio(payload.Metadata.Title, payload.Metadata.Year, imdbID)
			if err != nil {
				log.Error(err)
				return
			}

			// handle logic for aspect ratios
			topic := vip.GetString("mqtt.topicAspectratio")

			// better to have you just decide what to do in HA, I'm not your mom
			err = mqtt.Publish(vip, []byte(fmt.Sprintf("{\"aspect\":%f}", aspect)), topic)
			if err != nil {
				log.Error()
			}
		}
	}

}

// trigger HA for MV change per type
func changeMasterVolume(vip *viper.Viper, mediaType string, wg *sync.WaitGroup) {
	defer wg.Done()
	if vip.GetBool("homeAssistant.triggerAvrMasterVolumeChangeOnEvent") && vip.GetBool("homeAssistant.enabled") {
		log.Debug("changeMasterVolume: Changing volume")
		err := mqtt.Publish(vip, []byte(fmt.Sprintf("{\"type\":\"%s\"}", mediaType)), vip.GetString("mqtt.topicVolume"))
		if err != nil {
			log.Error()
		}
	}
}

// trigger HA for light change given entity and desired state
func changeLight(vip *viper.Viper, state string, wg *sync.WaitGroup) {
	defer wg.Done()
	if vip.GetBool("homeAssistant.triggerLightsOnEvent") && vip.GetBool("homeAssistant.enabled") {
		log.Debug("changeLight: Changing light")
		err := mqtt.Publish(vip, []byte(fmt.Sprintf("{\"state\":\"%s\"}", state)), vip.GetString("mqtt.topicLights"))
		if err != nil {
			log.Error()
		}
	}
}

// entry point for background tasks
func PlexWorker(plexChan <-chan models.PlexWebhookPayload, vip *viper.Viper) {
	log.Info("PlexWorker started")

	var beqClient *ezbeq.BeqClient
	var haClient *homeassistant.HomeAssistantClient
	var err error
	var deviceNames []string
	var model *models.SearchRequest
	var denonClient *denon.DenonClient
	var useDenonCodec bool

	// Server Info
	plexClient := plex.NewClient(vip.GetString("plex.url"), vip.GetString("plex.port"))

	if vip.GetBool("ezbeq.enabled") {
		log.Info("Started with ezbeq enabled")
		beqClient, err = ezbeq.NewClient(vip.GetString("ezbeq.url"), vip.GetString("ezbeq.port"))
		if err != nil {
			log.Error(err)
		}
		log.Debugf("Discovered devices: %v", beqClient.DeviceInfo)
		if len(beqClient.DeviceInfo) == 0 {
			log.Error("No devices found. Please check your ezbeq settings!")
			return
		}

		// get the device names from the API call
		for _, k := range beqClient.DeviceInfo {
			log.Debugf("adding device %s", k.Name)
			deviceNames = append(deviceNames, k.Name)
		}

		log.Debugf("Device names: %v", deviceNames)
		model = &models.SearchRequest{
			DryrunMode: vip.GetBool("ezbeq.dryRun"),
			Devices:    deviceNames,
			Slots:      vip.GetIntSlice("ezbeq.slots"),
			// try to skip by default
			SkipSearch:      true,
			PreferredAuthor: vip.GetString("ezbeq.preferredAuthor"),
		}

		// unload existing profile for safety
		err = beqClient.UnloadBeqProfile(model)
		if err != nil {
			log.Errorf("Error on startup - unloading beq %v", err)
		}
	}

	if vip.GetBool("homeAssistant.enabled") {
		log.Info("Started with HA enabled")
		haClient = homeassistant.NewClient(vip.GetString("homeAssistant.url"), vip.GetString("homeAssistant.port"), vip.GetString("homeAssistant.token"), vip.GetString("homeAssistant.envyName"))
	}
	if vip.GetBool("ezbeq.useAVRCodecSearch") {
		log.Info("Started with AVR codec search enabled")
		denonClient = denon.NewClient(vip.GetString("ezbeq.DenonIP"), vip.GetString("ezbeq.DenonPort"))
		useDenonCodec = true
	}

	// pointer so it can be modified by mediaPlay at will and be shared
	skipActions := new(bool)

	// block forever until closed so it will wait in background for work
	for i := range plexChan {
		log.Debugf("Current length of plexChan in PlexWorker: %d", len(plexChan))
		// determine what to do
		log.Debug("Sending new payload to eventRouter")
		eventRouter(plexClient, beqClient, haClient, denonClient, useDenonCodec, i, vip, model, skipActions)
		log.Debug("eventRouter done processing payload")
	}

	log.Debug("Plex worker stopped")
}
