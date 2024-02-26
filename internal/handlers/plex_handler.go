package handlers

import (
	// "encoding/json"
	"errors"
	"fmt"
	"net/http"

	// "strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
	"github.com/iloveicedgreentea/go-plex/internal/plex"
	"github.com/iloveicedgreentea/go-plex/models"
	"golang.org/x/exp/slices"
)

const showItemTitle = "episode"
const movieItemTitle = "movie"

var log = logger.GetLogger()

// interfaceRemote sends the cmd to your desired script to stop or play
func interfaceRemote(cmd string, c *homeassistant.HomeAssistantClient) error {
	switch cmd {
	case "play":
		return c.TriggerScript(config.GetString("homeAssistant.playScriptName"))
	case "pause":
		return c.TriggerScript(config.GetString("homeAssistant.pauseScriptName"))
	case "stop":
		return c.TriggerScript(config.GetString("homeAssistant.stopScriptName"))
	default:
		return errors.New("unknown cmd given")
	}

}

// Sends the payload to the channel for background processing
func ProcessWebhook(plexChan chan<- models.PlexWebhookPayload, c *gin.Context) {
	if err := c.Request.ParseMultipartForm(0); err != nil {
		log.Errorf("invalid multipart form: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"invalid multipart form": err.Error()})
		return
	}

	payload, hasPayload := c.Request.MultipartForm.Value["payload"]
	if hasPayload {
		log.Debug("decoding payload")
		decodedPayload, statusCode, err := common.DecodeWebhook(payload)
		if err != nil {
			log.Error(err)
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		clientUUID := decodedPayload.Player.UUID
		log.Infof("Got a request from UUID: %s", clientUUID)

		t := strings.ToLower(decodedPayload.Metadata.Type)

		log.Debugf("ProcessWebhook:  Media type is: %s", t)
		log.Debugf("ProcessWebhook:  Media title is: %s", decodedPayload.Metadata.Title)

		// check filter for user if not blank
		userID := config.GetString("plex.ownerNameFilter")
		// only respond to events on a particular account if you share servers and only for movies and shows
		// TODO: decodedPayload.Account.Title seems to always map to server owner not player account
		if userID == "" || decodedPayload.Account.Title == userID {
			if t == movieItemTitle || t == showItemTitle {
				log.Debug("adding item to plexChan")
				select {
				case plexChan <- decodedPayload:
					// send succeeded
					log.Debugf("Added length of plexChan: %d", len(plexChan))
					c.JSON(http.StatusOK, gin.H{"message": "Payload processed"})
				case <-time.After(time.Second * 3):
					log.Error("Send on plexChan timed out")
					c.JSON(http.StatusTooManyRequests, gin.H{"error": "Send on plexChan timed out"})
					return
				}
			} else {
				log.Debugf("Media type of %s is not supported", t)
			}
		} else {
			// TODO: this seems to be hitting even when the filter matches
			log.Debugf("userID '%s' does not match filter of %s", decodedPayload.Account.Title, userID)
		}
	} else {
		log.Error("No payload found in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No payload found in request"})
		return
	}
}

// does plex send stop if you exit with back button? - Yes, with X for mobile player as well
func mediaStop(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest) {
	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
	if err != nil {
		log.Error(err)
	}
	go common.ChangeLight("on")

	err = beqClient.UnloadBeqProfile(m)
	if err != nil {
		log.Error(err)
		if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
			err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err))
			if err != nil {
				log.Error()
			}
		}
	}
	log.Info("BEQ profile unloaded")
}

// pause only happens with literally pausing
func mediaPause(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest, skipActions *bool) {
	if !*skipActions {
		err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
		if err != nil {
			log.Error(err)
		}

		go common.ChangeLight("on")

		err = beqClient.UnloadBeqProfile(m)
		if err != nil {
			log.Error(err)
			if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err))
				if err != nil {
					log.Error()
				}
			}
		}
		log.Info("BEQ profile unloaded")
	}
}

func checkAvrCodec(client *plex.PlexClient, haClient *homeassistant.HomeAssistantClient, avrClient avr.AVRClient, payload models.PlexWebhookPayload, data models.MediaContainer) (codec string, err error) {
	// AVRs need time to decode the stream
	time.Sleep(3 * time.Second)

	// get the codec from avr
	avrCodec, err := avrClient.GetCodec()
	if err != nil {
		log.Errorf("error getting codec from denon, can't continue: %s", err)
		return codec, err
	}

	// TODO: try session data then fallback to lookup
	clientCodec, err := client.GetAudioCodec(data)
	if err != nil {
		log.Errorf("error getting codec from plex, can't continue: %s", err)
		return codec, err
	}
	codec = avr.MapDenonToBeq(avrCodec, clientCodec)
	// check if the expected codec is playing
	// TODO: check the plex metadata anyway and see if it contains Atmos
	if strings.Contains(clientCodec, "Atmos") {
		expectedCodec, err := common.IsAtmosCodecPlaying(codec, "Atmos")
		// TODO: retry above 5 times becaue it takes time to decode

		if err != nil {
			log.Errorf("Error checking if expected codec is playing: %v", err)
		}
		if !expectedCodec {
			// if enabled, stop playing
			if config.GetBool("ezbeq.stopPlexIfMismatch") {
				log.Debug("Stopping plex because codec is not playing")
				err := common.PlaybackInterface("stop", client)
				if err != nil {
					log.Errorf("Error stopping plex: %v", err)
				}
			}

			log.Error("Expected codec is not playing! Please check your AVR and Plex settings!")
			if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %v", clientCodec, avrCodec))
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return codec, err
}

// play is both the "resume" UI button and play
func mediaPlay(client *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, avrClient avr.AVRClient, payload models.PlexWebhookPayload, m *models.SearchRequest, useAvrCodec bool, data models.MediaContainer, skipActions *bool, wg *sync.WaitGroup) {
	var err error

	err = mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
	if err != nil {
		log.Error(err)
	}
	// dont need to set skipActions here because it will only send media.pause and media.resume. This is media.play
	go common.ChangeLight("off")
	go common.ChangeMasterVolume(m.MediaType)

	// optimistically try to hdmi sync. Will return if disabled
	wg.Add(1)
	go common.WaitForHDMISync(wg, skipActions, haClient, client)

	err = beqClient.UnloadBeqProfile(m)
	if err != nil {
		log.Errorf("error unloading beq during play %v", err)
	}
	// slower but more accurate especially with atmos
	if useAvrCodec {
		m.Codec, err = checkAvrCodec(client, haClient, avrClient, payload, data)
		// if it failed, get codec data from client
		if err != nil {
			log.Warnf("error getting codec from AVR, falling back to client: %s", err)
			m.Codec, err = client.GetAudioCodec(data)
			if err != nil {
				log.Errorf("error getting codec from plex, can't continue: %s", err)
				return
			}
		}
	} else {
		log.Debug("Using plex to get codec")
		// TODO: try session data then fallback to lookup
		m.Codec, err = client.GetAudioCodec(data)
		if err != nil {
			log.Errorf("error getting codec from plex, can't continue: %s", err)
			return
		}
	}

	log.Debugf("Found codec: %s", m.Codec)
	// if its a show and you dont want beq enabled, exit
	if payload.Metadata.Type == showItemTitle {
		if !config.GetBool("ezbeq.enableTvBeq") {
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
	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
		err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, m.Codec))
		if err != nil {
			log.Error()
		}
	}

	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("goroutines complete")
}

// resume is only after pausing as long as the media item is still active
func mediaResume(client *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest, data models.MediaContainer, skipActions *bool) {
	if !*skipActions {
		// mediaType string, codec string, edition string
		// trigger lights
		go common.ChangeLight("off")
		err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
		if err != nil {
			log.Error(err)
		}
		// Changing on resume is disabled because its annoying if you changed it since playing
		// go changeMasterVolume(vip, mediaType)

		// allow skipping search to save time
		// always unload in case something is loaded from movie for tv
		err = beqClient.UnloadBeqProfile(m)
		if err != nil {
			log.Errorf("Error on startup - unloading beq %v", err)
		}
		if payload.Metadata.Type == showItemTitle {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}
		// get the tmdb id to match with ezbeq catalog
		m.TMDB = getPlexMovieDb(payload)
		// if the server was restarted, cached data is lost
		if m.Codec == "" {
			log.Warn("No codec found in cache on resume. Was server restarted? Getting new codec")
			log.Debug("Using plex to get codec because its not cached")
			m.Codec, err = client.GetAudioCodec(data)
			if err != nil {
				log.Errorf("error getting codec from plex, can't continue: %s", err)
				return
			}
		}
		if m.Codec == "" {
			log.Error("No codec found after trying to resume")
			return
		}

		err = beqClient.LoadBeqProfile(m)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("BEQ profile loaded")

		// send notification of it loaded
		if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
			err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, m.Codec))
			if err != nil {
				log.Error()
			}
		}
	}
}

func mediaScrobble(beqClient *ezbeq.BeqClient, m *models.SearchRequest) {
	// trigger lights
	// go changeLight(vip, "on")
	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
	if err != nil {
		log.Error(err)
	}
	log.Debug("Scrobble received. Unloading profile")
	// unload beq
	err = beqClient.UnloadBeqProfile(m)
	if err != nil {
		log.Errorf("Error on startup - unloading beq %v", err)
	}

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
	// otherwise try to extract from file name
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

// ensure the client matches so it doesnt trigger from unwanted clients
func checkUUID(clientUUID string, filterConfig string) bool {

	// trim all spaces from the string
	clientUUID = strings.ReplaceAll(clientUUID, " ", "")
	filter := strings.ReplaceAll(filterConfig, " ", "") // trim all spaces from the string
	if filter == "" {
		log.Debug("No filter set, allowing all clients")
		return true
	}
	// split the filter string by comma
	filterArr := strings.Split(filter, ",")
	// iterate over each uuid in filterArr and compare with clientUUID
	if !slices.Contains(filterArr, clientUUID) {
		log.Debugf("filter '%s' does not match uuid '%s'", filterArr, clientUUID)
		return false
	}

	return true
}

// based on event type, determine what to do
func eventRouter(plexClient *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, avrClient avr.AVRClient, useAvrCodec bool, payload models.PlexWebhookPayload, model *models.SearchRequest, skipActions *bool) {
	// perform function via worker

	clientUUID := payload.Player.UUID
	// ensure the client matches so it doesnt trigger from unwanted clients

	if !checkUUID(clientUUID, config.GetString("plex.deviceUUIDFilter")) {
		log.Infof("Got a webhook but Client UUID '%s' does not match enabled filter", clientUUID)
		return
	}

	plexClient.MachineID = clientUUID
	plexClient.MediaType = payload.Metadata.Type
	log.Infof("Processing media type: %s", payload.Metadata.Type)

	var err error
	var data models.MediaContainer
	var editionName string

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

	log.Debugf("Event Router: Got media type of: %s ", payload.Metadata.Type)

	// mutate with data from plex
	model.Year = payload.Metadata.Year
	model.MediaType = payload.Metadata.Type
	model.Edition = editionName
	// this should be updated with every event
	model.EntryID = beqClient.CurrentProfile
	model.MVAdjust = beqClient.MasterVolume

	log.Debugf("Event Router: Using search model: %#v", model)
	log.Debugf("skipActions is: %v", *skipActions)
	switch payload.Event {
	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
	// play means a new file was started
	case "media.play":
		log.Debug("Event Router: media.play received")
		wg := &sync.WaitGroup{}
		mediaPlay(plexClient, beqClient, haClient, avrClient, payload, model, useAvrCodec, data, skipActions, wg)
	case "media.stop":
		log.Debug("Event Router: media.stop received")
		mediaStop(beqClient, haClient, payload, model)
	case "media.pause":
		log.Debug("Event Router: media.pause received")
		mediaPause(beqClient, haClient, payload, model, skipActions)
	// Pressing the 'resume' button in plex UI is media.play
	case "media.resume":
		log.Debug("Event Router: media.resume received")
		mediaResume(plexClient, beqClient, haClient, payload, model, data, skipActions)
	case "media.scrobble":
		log.Debug("Scrobble received")
		mediaScrobble(beqClient, model)
	default:
		log.Debugf("Received unsupported event: %s", payload.Event)
	}
}

// get the imdb ID from plex metadata
// func getPlexImdbID(payload models.PlexWebhookPayload) string {
// 	// try to get IMDB title from plex to save time
// 	for _, model := range payload.Metadata.GUID0 {
// 		if strings.Contains(model.ID, "imdb") {
// 			log.Debugf("Got title ID from plex - %s", model.ID)
// 			return strings.Split(model.ID, "imdb://")[1]
// 		}
// 	}

// 	return ""
// }

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
// func changeAspect(client *plex.PlexClient, payload models.PlexWebhookPayload, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	if config.GetBool("homeAssistant.triggerAspectRatioChangeOnEvent") && config.GetBool("homeAssistant.enabled") {

// 		// if madvr enabled, only send a trigger via mqtt
// 		// This needs to be triggered by MQTT in automation. This sends a pulse. The automation reads envy aspect ratio 5 sec later
// 		if config.GetBool("madvr.enabled") {
// 			log.Debug("Using madvr for aspect")
// 			topic := config.GetString("mqtt.topicAspectratioMadVrOnly")

// 			// trigger automation
// 			err := mqtt.Publish([]byte(""), topic)
// 			if err != nil {
// 				log.Error()
// 			}

// 			return
// 		} else {
// 			log.Debug("changeAspect: Changing aspect ratio")

// 			// get the imdb title ID
// 			imdbID := getPlexImdbID(payload)

// 			// lookup aspect based on imdb technical info
// 			aspect, err := client.GetAspectRatio(payload.Metadata.Title, payload.Metadata.Year, imdbID)
// 			if err != nil {
// 				log.Error(err)
// 				return
// 			}

// 			// handle logic for aspect ratios
// 			topic := config.GetString("mqtt.topicAspectratio")

// 			// better to have you just decide what to do in HA, I'm not your mom
// 			err = mqtt.Publish([]byte(fmt.Sprintf("{\"aspect\":%f}", aspect)), topic)
// 			if err != nil {
// 				log.Error()
// 			}
// 		}
// 	}

// }

// entry point for background tasks
func PlexWorker(plexChan <-chan models.PlexWebhookPayload, readyChan chan<- bool) {
	if !config.GetBool("plex.enabled") {
		log.Debug("Plex is disabled")
		readyChan <- true
		return
	}
	log.Info("PlexWorker started")

	var beqClient *ezbeq.BeqClient
	var haClient *homeassistant.HomeAssistantClient
	var err error
	var deviceNames []string
	var model *models.SearchRequest
	var avrClient avr.AVRClient
	var useAvrCodec bool

	// Server Info
	plexClient := plex.NewClient(config.GetString("plex.url"), config.GetString("plex.port"))

	log.Info("Started with ezbeq enabled")
	beqClient, err = ezbeq.NewClient(config.GetString("ezbeq.url"), config.GetString("ezbeq.port"))
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Discovered devices: %v", beqClient.DeviceInfo)
	if len(beqClient.DeviceInfo) == 0 {
		log.Error("No devices found. Please check your ezbeq settings!")
	}

	// get the device names from the API call
	for _, k := range beqClient.DeviceInfo {
		log.Debugf("adding device %s", k.Name)
		deviceNames = append(deviceNames, k.Name)
	}

	log.Debugf("Device names: %v", deviceNames)
	model = &models.SearchRequest{
		DryrunMode: config.GetBool("ezbeq.dryRun"),
		Devices:    deviceNames,
		Slots:      config.GetIntSlice("ezbeq.slots"),
		// try to skip by default
		SkipSearch:      true,
		PreferredAuthor: config.GetString("ezbeq.preferredAuthor"),
	}

	// unload existing profile for safety
	err = beqClient.UnloadBeqProfile(model)
	if err != nil {
		log.Errorf("Error on startup - unloading beq %v", err)
	}

	if config.GetBool("homeAssistant.enabled") {
		log.Info("Started with HA enabled")
		haClient = homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), config.GetString("homeAssistant.remoteentityname"))
	}

	if config.GetBool("ezbeq.useAVRCodecSearch") {
		log.Info("Started with AVR codec search enabled")
		avrClient = avr.GetAVRClient(config.GetString("ezbeq.avrurl"))
		if avrClient != nil {
			useAvrCodec = true
		}
	}

	// pointer so it can be modified by mediaPlay at will and be shared
	skipActions := new(bool)
	log.Info("Plex worker is ready")
	readyChan <- true
	// block forever until closed so it will wait in background for work
	for i := range plexChan {
		log.Debugf("Current length of plexChan in PlexWorker: %d", len(plexChan))
		// determine what to do
		log.Debug("Sending new payload to eventRouter")
		eventRouter(plexClient, beqClient, haClient, avrClient, useAvrCodec, i, model, skipActions)
		log.Debug("eventRouter done processing payload")
	}

	log.Debug("Plex worker stopped")
}
