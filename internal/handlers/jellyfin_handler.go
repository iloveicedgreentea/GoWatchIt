package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/jellyfin"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
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
		return
	}

	var payload models.JellyfinWebhook
	err = json.Unmarshal(read, &payload)
	if err != nil {
		log.Errorf("Error decoding payload: %v", err)
		log.Debugf("ProcessJfWebhook Request: %v", string(read))
		return
	}
	log.Debugf("Payload: %#v", payload)
	// respond to request with 200
	c.JSON(200, gin.H{"status": "ok"})
	// send payload to worker
	jfChan <- payload
}

func jfEventRouter(jfClient *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, model *models.SearchRequest, skipActions *bool) {
	// perform function via worker

	clientUUID := payload.ClientName
	// ensure the client matches so it doesnt trigger from unwanted clients

	if !checkUUID(clientUUID, config.GetString("jellyfin.deviceUUIDFilter")) {
		log.Infof("Got a webhook but Client UUID '%s' does not match enabled filter", clientUUID)
		return
	}

	var err error
	var data models.JellyfinMetadata
	var editionName string
	var codec string

	metadata, err := jfClient.GetMetadata(payload.UserID, payload.ItemID)
	if err != nil {
		log.Errorf("Error getting metadata from jellyfin API: %v", err)
		return
	}

	log.Debugf("Processing media type: %s", metadata.Type)

	// get the edition name
	editionName = jfClient.GetEdition(metadata)
	log.Debugf("Event Router: Found edition: %s", editionName)

	// mutate with data from JF
	year, err := strconv.Atoi(payload.Year)
	if err != nil {
		log.Errorf("Error converting year to int: %v", err)
		return
	}
	model.Year = year
	model.MediaType = metadata.Type
	model.Edition = editionName
	// this should be updated with every event
	model.EntryID = beqClient.CurrentProfile
	model.MVAdjust = beqClient.MasterVolume

	log.Debugf("Event Router: Using search model: %#v", model)
	log.Debugf("Got notification type %s", payload.NotificationType)
	if config.GetBool("ezbeq.useAVRCodecSearch") {
		// TODO: rewrite this
		c := avr.GetAVRClient(config.GetString("ezbeq.DenonIP"))
		if c != nil {
			codec, err = c.GetCodec()
			if err != nil {
				log.Errorf("Error getting codec from AVR: %v", err)
			}
			log.Debugf("Got codec from AVR: %s", codec)
			// TODO: make generic function that looks at which AVR and maps correctly
			codec = mapDenonToBeq(codec)
		} else {
			log.Error("Error getting AVR client. Trying to poll jellyfin")
			codec, err = jfClient.GetAudioCodec(metadata)
			if err != nil {
				log.Errorf("Error getting codec from jellyfin: %v", err)
				return
			}
		}
	} else {
		// return the normalized codec
		codec, err = jfClient.GetAudioCodec(metadata)
		if err != nil {
			log.Errorf("Error getting codec from jellyfin: %v", err)
		}
	}
	// add codec
	model.Codec = codec

	switch payload.NotificationType {
	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
	case "PlaybackStart":
		// TODO: test start/resume/pause
		jfMediaPlay(jfClient, beqClient, haClient, payload, model, false, data, skipActions)
	case "PlaybackStop":
		jfMediaStop(jfClient, beqClient, haClient, payload, model, false, data, skipActions)
	// really annoyingly jellyfin doesnt send a pause or resume event only progress every X seconds with a isPaused flag
	// TODO: support pause resume without running resume on every playbackprogress
	case "PlaybackProgress":
		if payload.IsPaused == "true" {
			jfMediaPause(beqClient, haClient, payload, model, skipActions)
		} else {
			jfMediaResume(jfClient, beqClient, haClient, payload, model, false, data, skipActions)
		}
	default:
		log.Warnf("Received unsupported webhook event. Nothing to do: %s", payload.NotificationType)
	}
}

func jfMediaPlay(client *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, m *models.SearchRequest, useDenonCodec bool, data models.JellyfinMetadata, skipActions *bool) {
	log.Debug("Processing media play event")
	wg := &sync.WaitGroup{}

	// stop processing webhooks
	*skipActions = true
	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
	if err != nil {
		log.Error(err)
	}
	go common.ChangeLight("off")
	// go changeAspect(client, payload, wg)
	go common.ChangeMasterVolume(m.MediaType)

	// if not using denoncodec, do this in background
	if !useDenonCodec {
		wg.Add(1)
		// sets skipActions to false on completion
		go common.WaitForHDMISync(wg, skipActions, haClient, client)
	}

	// always unload in case something is loaded from movie for tv
	err = beqClient.UnloadBeqProfile(m)
	if err != nil {
		log.Errorf("Error unloading beq on startup!! : %v", err)
		return
	}

	// TODO: check if beq is enabled
	// if its a show and you dont want beq enabled, exit
	if data.Type == showItemTitle {
		if !config.GetBool("ezbeq.enableTvBeq") {
			return
		}
	}

	m.TMDB, err = client.GetJfTMDB(data)
	if err != nil {
		if config.GetBool("jellyfin.skiptmdb") {
			log.Warn("TMDB data not found. TMDB is allowed to be skipped")
		} else {
			log.Errorf("Error getting TMDB data from metadata: %v", err)
			return
		}
	}
	err = beqClient.LoadBeqProfile(m)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("BEQ profile loaded")

	// send notification of it loaded
	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
		err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%s) // Codec %s", data.OriginalTitle, payload.Year, m.Codec))
		if err != nil {
			log.Error()
		}
	}

	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("goroutines complete")
}

func jfMediaStop(client *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, m *models.SearchRequest, useDenonCodec bool, data models.JellyfinMetadata, skipActions *bool) {
	log.Debug("Processing media stop event")
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

func jfMediaPause(beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, m *models.SearchRequest, skipActions *bool) {
	log.Debug("Processing media pause event")
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
func jfMediaResume(client *jellyfin.JellyfinClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.JellyfinWebhook, m *models.SearchRequest, useDenonCodec bool, data models.JellyfinMetadata, skipActions *bool) {
	log.Debug("Processing media resume event")
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
		if data.Type == showItemTitle {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}
		// get the tmdb id to match with ezbeq catalog
		m.TMDB, err = client.GetJfTMDB(data)
		if err != nil {
			log.Errorf("Error getting TMDB data from metadata: %v", err)
			return
		}
		// if the server was restarted, cached data is lost
		if m.Codec == "" {
			log.Warn("No codec found in cache on resume. Was server restarted? Getting new codec")
			log.Debug("Using jellyfin to get codec because its not cached")
			m.Codec, err = client.GetAudioCodec(data)
			if err != nil {
				log.Errorf("error getting codec from jellyfin, can't continue: %s", err)
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
			err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%s) // Codec %s", data.OriginalTitle, payload.Year, m.Codec))
			if err != nil {
				log.Error()
			}
		}
	}
}

// // entry point for background tasks
func JellyfinWorker(jfChan <-chan models.JellyfinWebhook, readyChan chan<- bool) {
	if !config.GetBool("jellyfin.enabled") {
		log.Debug("Jellyfin is disabled")
		readyChan <- true
		return
	}

	// Server Info
	jellyfinClient := jellyfin.NewClient(config.GetString("jellyfin.url"), config.GetString("jellyfin.port"), config.GetString("jellyfin.playerMachineIdentifier"), config.GetString("jellyfin.playerIP"))

	var beqClient *ezbeq.BeqClient
	var haClient *homeassistant.HomeAssistantClient
	var err error
	var deviceNames []string
	var model *models.SearchRequest
	// var denonClient *denon.DenonClient
	// var useDenonCodec bool

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
	// if config.GetBool("ezbeq.useAVRCodecSearch") {
	// 	log.Info("Started with AVR codec search enabled")
	// 	denonClient = denon.NewClient(config.GetString("ezbeq.DenonIP"), config.GetString("ezbeq.DenonPort"))
	// 	useDenonCodec = true
	// }

	// pointer so it can be modified by mediaPlay at will and be shared
	skipActions := new(bool)
	readyChan <- true
	log.Info("JellyfinWorker is ready")
	// block forever until closed so it will wait in background for work
	for i := range jfChan {
		log.Debugf("Sending new payload to eventRouter - %#v", i)
		// if its not an empty struct
		if i != (models.JellyfinWebhook{}) {
			// get metadata
			jfEventRouter(jellyfinClient, beqClient, haClient, i, model, skipActions)
		} else {
			log.Warning("Received empty payload, skipping")
		}
		log.Debug("eventRouter done processing payload")
	}

	log.Info("JellyfinWorker worker stopped")
}
