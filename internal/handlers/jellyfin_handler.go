package handlers

import (
	"encoding/json"
	"io"
// 	"strconv"

// 	"sync"

	"github.com/gin-gonic/gin"
// 	"github.com/iloveicedgreentea/go-plex/internal/avr"
// 	"github.com/iloveicedgreentea/go-plex/internal/config"
// 	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
// 	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
// 	"github.com/iloveicedgreentea/go-plex/internal/jellyfin"
// 	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
// 	// "github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/models"
)

func ProcessJfWebhook(jfChan chan<- models.JellyfinWebhook, c *gin.Context) {
	// send payload to worker
	log.Debug("Sending payload to JellyfinWorker")
	r := c.Request.Body
	defer r.Close()
	read, err := io.ReadAll(r)
	if err != nil {
		log.Errorf("Error reading request body: %v", err)
	}

	// log.Debugf("ProcessJfWebhook Request: %v", string(read))
	var payload models.JellyfinWebhook
	err = json.Unmarshal(read, &payload)
	if err != nil {
		log.Errorf("Error decoding payload: %v", err)
	}
	log.Debugf("Payload: %#v", payload)
	// respond to request with 200
	c.JSON(200, gin.H{"status": "ok"})
	// send payload to worker
	jfChan <- payload
}

// func jfEventRouter(jfClient *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, model *models.SearchRequest, skipActions *bool) {
// 	// perform function via worker

// 	clientUUID := payload.ClientName
// 	// ensure the client matches so it doesnt trigger from unwanted clients

// 	if !checkUUID(clientUUID, config.GetString("plex.deviceUUIDFilter")) {
// 		log.Infof("Got a webhook but Client UUID '%s' does not match enabled filter", clientUUID)
// 		return
// 	}

// 	var err error
// 	var data models.MediaContainer
// 	var editionName string
// 	metadata, err := jfClient.GetMetadata(payload.UserID, payload.ItemID)
// 	if err != nil {
// 		log.Errorf("Error getting metadata from jellyfin API: %v", err)
// 	}

// 	log.Debugf("Processing media type: %s", metadata.Type)

// 	// get the edition name
// 	editionName = jfClient.GetEdition(metadata)
// 	log.Debugf("Event Router: Found edition: %s", editionName)

// 	// mutate with data from plex
// 	year, err := strconv.Atoi(payload.Year)
// 	if err != nil {
// 		log.Errorf("Error converting year to int: %v", err)
// 		return
// 	}
// 	model.Year = year
// 	model.MediaType = metadata.Type
// 	model.Edition = editionName
// 	// this should be updated with every event
// 	model.EntryID = beqClient.CurrentProfile
// 	model.MVAdjust = beqClient.MasterVolume

// 	log.Debugf("Event Router: Using search model: %#v", model)
// 	log.Debugf("Got notification type %s", payload.NotificationType)
// 	switch payload.NotificationType {
// 	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
// 	// play means a new file was started
// 	case "PlaybackStart":
// 		log.Debug("Event Router: media.play received")
// 		if config.GetBool("ezbeq.useAVRCodecSearch") {
// 			c := avr.GetAVRClient(config.GetString("ezbeq.DenonIP"))
// 			if c != nil {
// 				codec, err := c.GetCodec()
// 				if err != nil {
// 					log.Errorf("Error getting codec from AVR: %v", err)
// 				}
// 				log.Debugf("Got codec from AVR: %s", codec)
// 				model.Codec = mapDenonToBeq(codec)
// 			} else {
// 				log.Error("Error getting AVR client")
// 				model.Codec = ""
// 			}
// 		} else {
// 			codec, displayTitle, codecProfile, err := jfClient.GetCodec(metadata)
// 			if err != nil {
// 				log.Errorf("Error getting codec from jellyfin: %v", err)
// 			}
// 		}
// 		// TODO: normalize codec
// 		commonPlay(jfClient, beqClient, haClient, payload, model, false, data, skipActions)
// 	case "PlaybackStop":
// 		log.Debug("Event Router: media.stop received")
// 		jfMediaStop(beqClient, haClient, payload, model)
// 	// really annoyingly jellyfin doesnt send a pause or resume
// 	// TODO: support pause resume without running resume on every playbackprogress
// 	// case "PlaybackProgress":
// 	// 	log.Debug("Event Router: PlaybackProgress received")
// 	// 	if payload.IsPaused == "true" {
// 	// 		mediaPause(beqClient, haClient, payload, model, skipActions)
// 	// 	} else {
// 	// 		mediaResume()
// 	// 	}
// 	// Pressing the 'resume' button in plex is media.play
// 	default:
// 		log.Debugf("Received unsupported event: %s", payload.NotificationType)
// 	}
// }

// func jfMediaPlay(client *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, m *models.SearchRequest, useDenonCodec bool, data models.MediaContainer, skipActions *bool) {
// 	wg := &sync.WaitGroup{}

// 	// stop processing webhooks
// 	*skipActions = true
// 	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	go changeLight("off")
// 	// go changeAspect(client, payload, wg)
// 	go changeMasterVolume(m.MediaType)

// 	// if not using denoncodec, do this in background
// 	if !useDenonCodec {
// 		wg.Add(1)
// 		// sets skipActions to false on completion
// 		go waitForHDMISync(wg, skipActions, haClient, client)
// 	}

// 	// always unload in case something is loaded from movie for tv
// 	err = beqClient.UnloadBeqProfile(m)
// 	if err != nil {
// 		log.Errorf("Error unloading beq on startup!! : %v", err)
// 		return
// 	}

// 	// slower but more accurate
// 	// TODO: abstract library this for any AVR
// 	if useDenonCodec {
// 		// TODO: make below a function
// 		// wait for sync
// 		wg.Add(1)
// 		waitForHDMISync(wg, skipActions, haClient, client)
// 		// denon needs time to show mutli ch in as atmos
// 		// TODO: test this
// 		time.Sleep(5 * time.Second)

// 		// get the codec from avr
// 		m.Codec, err = denonClient.GetCodec()
// 		if err != nil {
// 			log.Errorf("error getting codec from denon, can't continue: %s", err)
// 			return
// 		}

// 		// check if the expected codec is playing
// 		expectedCodec, isExpectedPlaying := isExpectedCodecPlaying(denonClient, client, payload.Player.UUID, m.Codec)
// 		if !isExpectedPlaying {
// 			// if enabled, stop playing
// 			if config.GetBool("ezbeq.stopPlexIfMismatch") {
// 				log.Debug("Stopping plex because codec is not playing")
// 				err := playbackInteface("stop", haClient, client)
// 				if err != nil {
// 					log.Errorf("Error stopping plex: %v", err)
// 				}
// 			}

// 			log.Error("Expected codec is not playing! Please check your AVR and Plex settings!")
// 			if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
// 				err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %s", m.Codec, expectedCodec), config.GetString("ezbeq.notifyEndpointName"))
// 				if err != nil {
// 					log.Error(err)
// 				}
// 			}
// 		}

// 	} else {
// 		m.Codec, err = client.GetAudioCodec(data)
// 		if err != nil {
// 			log.Errorf("error getting codec from plex, can't continue: %s", err)
// 			return
// 		}
// 	}

// 	log.Debugf("Found codec: %s", m.Codec)
// 	// TODO: check if beq is enabled
// 	// if its a show and you dont want beq enabled, exit
// 	if payload.Metadata.Type == showItemTitle {
// 		if !config.GetBool("ezbeq.enableTvBeq") {
// 			return
// 		}
// 	}

// 	m.TMDB = getPlexMovieDb(payload)
// 	err = beqClient.LoadBeqProfile(m)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	log.Info("BEQ profile loaded")

// 	// send notification of it loaded
// 	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
// 		err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, m.Codec), config.GetString("ezbeq.notifyEndpointName"))
// 		if err != nil {
// 			log.Error()
// 		}
// 	}

// 	log.Debug("Waiting for goroutines")
// 	wg.Wait()
// 	log.Debug("goroutines complete")
// }

// func jfMediaStop()

// // entry point for background tasks
func JellyfinWorker(jfChan <-chan models.JellyfinWebhook, readyChan chan<- bool) {
	readyChan <- true
	
	// log.Info("JellyfinWorker started")

	// // Server Info
	// jellyfinClient := jellyfin.NewClient(config.GetString("jellyfin.url"), config.GetString("jellyfin.port"), config.GetString("jellyfin.playerMachineIdentifier"), config.GetString("jellyfin.playerIP"))

	// var beqClient *ezbeq.BeqClient
	// var haClient *homeassistant.HomeAssistantClient
	// var err error
	// var deviceNames []string
	// var model *models.SearchRequest
	// // var denonClient *denon.DenonClient
	// // var useDenonCodec bool

	// log.Info("Started with ezbeq enabled")
	// beqClient, err = ezbeq.NewClient(config.GetString("ezbeq.url"), config.GetString("ezbeq.port"))
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Debugf("Discovered devices: %v", beqClient.DeviceInfo)
	// if len(beqClient.DeviceInfo) == 0 {
	// 	log.Error("No devices found. Please check your ezbeq settings!")
	// }

	// // get the device names from the API call
	// for _, k := range beqClient.DeviceInfo {
	// 	log.Debugf("adding device %s", k.Name)
	// 	deviceNames = append(deviceNames, k.Name)
	// }

	// log.Debugf("Device names: %v", deviceNames)
	// model = &models.SearchRequest{
	// 	DryrunMode: config.GetBool("ezbeq.dryRun"),
	// 	Devices:    deviceNames,
	// 	Slots:      config.GetIntSlice("ezbeq.slots"),
	// 	// try to skip by default
	// 	SkipSearch: true,
	// 	// TODO: make this a whitelist
	// 	PreferredAuthor: config.GetString("ezbeq.preferredAuthor"),
	// }

	// // unload existing profile for safety
	// err = beqClient.UnloadBeqProfile(model)
	// if err != nil {
	// 	log.Errorf("Error on startup - unloading beq %v", err)
	// }

	// if config.GetBool("homeAssistant.enabled") {
	// 	log.Info("Started with HA enabled")
	// 	haClient = homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), config.GetString("homeAssistant.remoteentityname"))
	// }
	// // if config.GetBool("ezbeq.useAVRCodecSearch") {
	// // 	log.Info("Started with AVR codec search enabled")
	// // 	denonClient = denon.NewClient(config.GetString("ezbeq.DenonIP"), config.GetString("ezbeq.DenonPort"))
	// // 	useDenonCodec = true
	// // }

	// // pointer so it can be modified by mediaPlay at will and be shared
	// skipActions := new(bool)
	// readyChan <- true
	// log.Info("JellyfinWorker is ready")
	// // block forever until closed so it will wait in background for work
	// for i := range jfChan {
	// 	log.Debugf("Sending new payload to eventRouter - %#v", i)
	// 	// if its not an empty struct
	// 	if i != (models.JellyfinWebhook{}) {
	// 		// get metadata
	// 		jfEventRouter(jellyfinClient, beqClient, haClient, i, model, skipActions)
	// 	} else {
	// 		log.Warning("Received empty payload, skipping")
	// 	}
	// 	log.Debug("eventRouter done processing payload")
	// }

	// log.Info("JellyfinWorker worker stopped")
}
