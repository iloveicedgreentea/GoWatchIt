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

		log.Debugf("ProcessWebhook:  Media type is: %s", decodedPayload.Metadata.Type)
		log.Debugf("ProcessWebhook:  Media title is: %s", decodedPayload.Metadata.Title)

		// check filter for user if not blank
		userID := config.GetString("plex.ownerNameFilter")
		// only respond to events on a particular account if you share servers and only for movies and shows
		// TODO: decodedPayload.Account.Title seems to always map to server owner not player account
		if userID == "" || decodedPayload.Account.Title == userID {
			if decodedPayload.Metadata.Type == movieItemTitle || decodedPayload.Metadata.Type == showItemTitle {
				select {
				case plexChan <- decodedPayload:
					// send succeeded
					c.JSON(http.StatusOK, gin.H{"message": "Payload processed"})
				case <-time.After(time.Second * 3):
					log.Error("Send on plexChan timed out")
					c.JSON(http.StatusTooManyRequests, gin.H{"error": "Send on plexChan timed out"})
					return
				}
				log.Debugf("Added length of plexChan: %d", len(plexChan))
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

// play is both the "resume" button and play
func mediaPlay(client *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, avrClient avr.AVRClient, payload models.PlexWebhookPayload, m *models.SearchRequest, useAvrCodec bool, data models.MediaContainer, skipActions *bool, wg *sync.WaitGroup) {
	go common.ChangeLight("off")
	var err error
	// slower but more accurate
	// TODO: abstract library this for any AVR
	if useAvrCodec {
		// TODO: make below a function
		// wait for sync
		wg.Add(1)
		common.WaitForHDMISync(wg, skipActions, haClient, client)
		// denon needs time to show mutli ch in as atmos
		// TODO: test this
		time.Sleep(5 * time.Second)

		// get the codec from avr
		m.Codec, err = avrClient.GetCodec()
		if err != nil {
			log.Errorf("error getting codec from denon, can't continue: %s", err)
			return
		}

		// check if the expected codec is playing
		// TODO: test this
		expectedCodec, err := common.IsExpectedCodecPlayingAVR(avrClient, m.Codec)
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
				err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %v", m.Codec, expectedCodec))
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
	// TODO: check if beq is enabled
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
func mediaResume(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, m *models.SearchRequest, skipActions *bool) {
	if !*skipActions {
		// mediaType string, codec string, edition string
		// trigger lights
		err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
		if err != nil {
			log.Error(err)
		}
		go common.ChangeLight("off")
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
		// load beq with cache
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

func mediaScrobble() {
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

// TODO! make a generic eventRouter but route to implementation specific functions instead of making generic play functions
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
	switch payload.Event {
	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
	// play means a new file was started
	case "media.play":
		log.Debug("Event Router: media.play received")
		// TODO: add lights and stuff here to do async, not blocked by other functions
		wg := &sync.WaitGroup{}
		mediaPlay(plexClient, beqClient, haClient, avrClient, payload, model, useAvrCodec, data, skipActions, wg)
	case "media.stop":
		log.Debug("Event Router: media.stop received")
		mediaStop(beqClient, haClient, payload, model)
	case "media.pause":
		log.Debug("Event Router: media.pause received")
		mediaPause(beqClient, haClient, payload, model, skipActions)
	// Pressing the 'resume' button in plex is media.play
	case "media.resume":
		log.Debug("Event Router: media.resume received")
		mediaResume(beqClient, haClient, payload, model, skipActions)
	case "media.scrobble":
		log.Debug("Scrobble received")
		mediaScrobble()
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
	log.Info("PlexWorker started")

	var beqClient *ezbeq.BeqClient
	var haClient *homeassistant.HomeAssistantClient
	var err error
	var deviceNames []string
	var model *models.SearchRequest
	var avrClient avr.AVRClient
	var useAvrCodec bool

	// Server Info
	plexClient := plex.NewClient(config.GetString("plex.url"), config.GetString("plex.port"), config.GetString("plex.playerMachineIdentifier"), config.GetString("plex.playerIP"))

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
		SkipSearch: true,
		// TODO: make this a whitelist
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
