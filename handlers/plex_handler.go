package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

// for trailers
const clipTitle = "clip"

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
		// fyi, sometimes media.play is not *SENT* by plex, I am investigating
		if err := r.ParseMultipartForm(0); err != nil {
			log.Error(err)
			return
		}

		payload, hasPayload := r.MultipartForm.Value["payload"]
		if hasPayload {
			decodedPayload, statusCode, err := decodeWebhook(payload)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), statusCode)
				return
			}
			clientUUID := decodedPayload.Player.UUID
			log.Debugf("!!! Your Player UUID is %s !!!!!", clientUUID)

			log.Debugf("ProcessWebhook:  Media type is: %s", decodedPayload.Metadata.Type)
			log.Debugf("ProcessWebhook:  Media title is: %s", decodedPayload.Metadata.Title)

			// check filter for user if not blank
			userID := vip.GetString("plex.ownerNameFilter")
			// only respond to events on a particular account if you share servers and only for movies and shows
			if userID == "" || decodedPayload.Account.Title == userID {
				if decodedPayload.Metadata.Type == movieItemTitle || decodedPayload.Metadata.Type == showItemTitle || decodedPayload.Metadata.Type == clipTitle {
					plexChan <- decodedPayload
				}
			} else {
				log.Debugf("userID '%s' does not match filter", userID)
			}
		}
	}

	return http.HandlerFunc(fn)
}

// does plex send stop if you exit with back button? - Yes, with X for mobile player as well
func mediaStop(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload) {
	go changeLight(vip, "on")

	if vip.GetBool("ezbeq.enabled") {
		err := beqClient.UnloadBeqProfile(vip.GetBool("ezbeq.dryRun"))
		if err != nil {
			log.Error(err)
			if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err), vip.GetString("ezbeq.notifyEndpointName"))
				if err != nil {
					log.Error()
				}
			}
		}
	}
}

// pause only happens with literally pausing
func mediaPause(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload) {
	go changeLight(vip, "on")

	if vip.GetBool("ezbeq.enabled") {
		err := beqClient.UnloadBeqProfile(vip.GetBool("ezbeq.dryRun"))
		if err != nil {
			log.Error(err)
			if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
				err := haClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err), vip.GetString("ezbeq.notifyEndpointName"))
				if err != nil {
					log.Error()
				}
			}
		}
	}
}

// play is both the "resume" button and play
func mediaPlay(client *plex.PlexClient, vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, mediaType string, codec string, edition string) {
	go changeLight(vip, "off")
	go changeAspect(client, payload, vip)
	go changeMasterVolume(vip, mediaType)

	// TODO: function to check expected codec, poll avr directly

	if vip.GetBool("ezbeq.enabled") {
		// always unload in case something is loaded from movie for tv
		err := beqClient.UnloadBeqProfile(false)
		if err != nil {
			log.Errorf("Error unloading beq on startup!! : %v", err)
		}

		// if its a show and you dont want beq enabled, exit
		if payload.Metadata.Type == showItemTitle {
			if !vip.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}

		tmdb := getPlexMovieDb(payload)
		err = beqClient.LoadBeqProfile(tmdb, payload.Metadata.Year, codec, false, "", 0, vip.GetBool("ezbeq.dryRun"), vip.GetString("ezbeq.preferredAuthor"), edition, mediaType)
		if err != nil {
			log.Error(err)
			return
		}
		// send notification of it loaded
		if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
			err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, codec), vip.GetString("ezbeq.notifyEndpointName"))
			if err != nil {
				log.Error()
			}
		}
	}
}

// resume is only after pausing as long as the media item is still active
func mediaResume(vip *viper.Viper, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, mediaType string, codec string, edition string) {

	// trigger lights
	go changeLight(vip, "off")
	// Changing on resume is disabled because its annoying if you changed it since playing
	// go changeMasterVolume(vip, mediaType)

	// allow skipping search to save time
	if vip.GetBool("ezbeq.enabled") {
		// always unload in case something is loaded from movie for tv
		err := beqClient.UnloadBeqProfile(false)
		if err != nil {
			log.Errorf("Error on startup - unloading beq %v", err)
		}
		if payload.Metadata.Type == showItemTitle {
			if !vip.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}
		// get the tmdb id to match with ezbeq catalog
		tmdb := getPlexMovieDb(payload)
		// load beq with cache
		err = beqClient.LoadBeqProfile(models.SearchRequest{
			TMDB: tmdb,
			Year: payload.Metadata.Year,
			Codec: codec,
			SkipSearch: true,
			EntryID: beqClient.CurrentProfile,
			MVAdjust: beqClient.MasterVolume,
			DryrunMode: vip.GetBool("ezbeq.dryRun"),
			PreferredAuthor: vip.GetString("ezbeq.preferredAuthor"),
			Edition: edition,
			MediaType: mediaType,
			Devices: vip.GetStringSlice("ezbeq.devices"),
			Slots: vip.GetIntSlice("ezbeq.slots"),
		})
		if err != nil {
			log.Error(err)
			return
		}
		// send notification of it loaded
		if vip.GetBool("ezbeq.notifyOnLoad") && vip.GetBool("homeAssistant.enabled") {
			err := haClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, codec), vip.GetString("ezbeq.notifyEndpointName"))
			if err != nil {
				log.Error()
			}
		}
	}
}

func mediaScrobble(vip *viper.Viper) {
	// trigger lights
	go changeLight(vip, "on")
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
func eventRouter(client *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, vip *viper.Viper) {
	// perform function via worker

	clientUUID := payload.Player.UUID
	// ensure the client matches so it doesnt trigger from unwanted clients

	if vip.GetString("plex.deviceUUIDFilter") != clientUUID || vip.GetString("plex.deviceUUIDFilter") == "" {
		log.Debug("Client UUID does not match enabled filter")
		return
	}

	// if its a clip and you didnt enable support for it, return
	if payload.Metadata.Type == clipTitle {
		if !vip.GetBool("plex.enableTrailerSupport") {
			log.Debug("Clip received but support not enabled")
			return
		}
	}

	log.Infof("Processing media type: %s", payload.Metadata.Type)

	var codec string
	var err error
	var data models.MediaContainer
	var editionName string

	if vip.GetBool("ezbeq.enabled") {
		// make a call to plex to get the data based on key
		data, err = client.GetMediaData(payload.Metadata.Key)
		if err != nil {
			log.Error(err)
			return
		} else {
			// get the edition name
			editionName = getEditionName(data)
			log.Debugf("Event Router: Found edition: %s", editionName)

			log.Debug("Event Router: Getting codec from data")
			codec, err = client.GetAudioCodec(data)
			if err != nil {
				log.Errorf("Event Router: error getting codec, can't continue: %s", err)

				return
			}
		}
	}

	log.Debugf("Event Router: Received codec: %s", codec)
	log.Debugf("Event Router: Got media type of: %s ", payload.Metadata.Type)

	switch payload.Event {
	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
	// play means a new file was started
	case "media.play":
		log.Debug("Event Router: media.play received")
		mediaPlay(client, vip, beqClient, haClient, payload, payload.Metadata.Type, codec, editionName)
	case "media.stop":
		log.Debug("Event Router: media.stop received")
		mediaStop(vip, beqClient, haClient, payload)
	case "media.pause":
		log.Debug("Event Router: media.pause received")
		mediaPause(vip, beqClient, haClient, payload)
	// Pressing the 'resume' button actually is media.play thankfully
	case "media.resume":
		log.Debug("Event Router: media.resume received")
		mediaResume(vip, beqClient, haClient, payload, payload.Metadata.Type, codec, editionName)
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
func changeAspect(client *plex.PlexClient, payload models.PlexWebhookPayload, vip *viper.Viper) {
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
func changeMasterVolume(vip *viper.Viper, mediaType string) {
	if vip.GetBool("homeAssistant.triggerAvrMasterVolumeChangeOnEvent") && vip.GetBool("homeAssistant.enabled") {
		log.Debug("changeMasterVolume: Changing volume")
		err := mqtt.Publish(vip, []byte(fmt.Sprintf("{\"type\":\"%s\"}", mediaType)), vip.GetString("mqtt.topicVolume"))
		if err != nil {
			log.Error()
		}
	}
}

// trigger HA for light change given entity and desired state
func changeLight(vip *viper.Viper, state string) {
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

	// Server Info
	plexClient := plex.NewClient(vip.GetString("plex.url"), vip.GetString("plex.port"))

	if vip.GetBool("ezbeq.enabled") {
		log.Info("Started with ezbeq enabled")
		beqClient, err = ezbeq.NewClient(vip.GetString("ezbeq.url"), vip.GetString("ezbeq.port"))
		if err != nil {
			log.Error(err)
		}
		// unload existing profile for safety
		err = beqClient.UnloadBeqProfile(false)
		if err != nil {
			log.Errorf("Error on startup - unloading beq %v", err)
		}
	}
	if vip.GetBool("homeAssistant.enabled") {
		log.Info("Started with HA enabled")
		haClient = homeassistant.NewClient(vip.GetString("homeAssistant.url"), vip.GetString("homeAssistant.port"), vip.GetString("homeAssistant.token"))
	}
	// block forever until closed so it will wait in background for work
	for i := range plexChan {
		// determine what to do
		eventRouter(plexClient, beqClient, haClient, i, vip)
	}
}
