package handlers

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io"
	// "net/http"
	// "strings"

	// "github.com/iloveicedgreentea/go-plex/ezbeq"
	// "github.com/iloveicedgreentea/go-plex/homeassistant"
	// "github.com/iloveicedgreentea/go-plex/models"
	// "github.com/iloveicedgreentea/go-plex/mqtt"
	// "github.com/spf13/viper"
)

// const showItemTitle = "episode"
// const movieItemTitle = "movie"

// // for trailers
// const clipTitle = "clip"

func decodeJellyfinWebhook(payload []string) () {
	// var pwhPayload models.PlexWebhookPayload
	log.Debug(payload)
	return

	// err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	// if err != nil {
	// 	var unmarshalTypeError *json.UnmarshalTypeError

	// 	switch {
	// 	// unmarshall error
	// 	case errors.As(err, &unmarshalTypeError):
	// 		msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	// 		return pwhPayload, http.StatusBadRequest, errors.New(msg)

	// 	default:
	// 		return pwhPayload, http.StatusInternalServerError, err
	// 	}
	// }

	// log.Debugf("decodeWebhook: Received event: %s", pwhPayload.Event)
	// return pwhPayload, 0, nil
}

// Sends the payload to the channel for background processing
// func ProcessJellyfinWebhook(plexChan chan<- models.PlexWebhookPayload, vip *viper.Viper) http.Handler {
// 	// TODO: pause is not working
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// fyi, sometimes media.play is not *SENT* by plex, I am investigating
// 		b, _ := io.ReadAll(r.Body)
// 		log.Debug(string(b))
// 	}

// 	return http.HandlerFunc(fn)
// }

// getEditionName tries to extract the edition from plex or file name. Assumes you have well named files
// Returned types, Unrated, Ultimate, Theatrical, Extended, Director, Criterion
// func getJellyfinEditionName(data models.MediaContainer) string {
// 	edition := data.Video.EditionTitle
// 	fileName := strings.ToLower(data.Video.Media.Part.File)

// 	// if there is an edition from plex metadata, use it
// 	if edition != "" {
// 		return edition
// 	}

// 	switch {
// 	case strings.Contains(fileName, "extended"):
// 		return "Extended"
// 	case strings.Contains(fileName, "unrated"):
// 		return "Unrated"
// 	case strings.Contains(fileName, "theatrical"):
// 		return "Theatrical"
// 	case strings.Contains(fileName, "ultimate"):
// 		return "Ultimate"
// 	case strings.Contains(fileName, "director"):
// 		return "Director"
// 	case strings.Contains(fileName, "criterion"):
// 		return "Criterion"
// 	default:
// 		return ""
// 	}
// }

// based on event type, determine what to do
// func eventJellyfinRouter(client *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, payload models.PlexWebhookPayload, vip *viper.Viper) {
// 	// perform function via worker

// 	clientUUID := payload.Player.UUID
// 	// ensure the client matches so it doesnt trigger from unwanted clients

// 	if vip.GetString("plex.deviceUUIDFilter") != clientUUID || vip.GetString("plex.deviceUUIDFilter") == "" {
// 		log.Debug("Client UUID does not match enabled filter")
// 		return
// 	}

// 	// if its a clip and you didnt enable support for it, return
// 	if payload.Metadata.Type == clipTitle {
// 		if !vip.GetBool("plex.enableTrailerSupport") {
// 			log.Debug("Clip received but support not enabled")
// 			return
// 		}
// 	}

// 	log.Infof("Processing media type: %s", payload.Metadata.Type)

// 	var codec string
// 	var err error
// 	var data models.MediaContainer
// 	var editionName string

// 	if vip.GetBool("ezbeq.enabled") {
// 		// make a call to plex to get the data based on key
// 		data, err = client.GetMediaData(payload.Metadata.Key)
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		} else {
// 			// get the edition name
// 			editionName = getEditionName(data)
// 			log.Debugf("Event Router: Found edition: %s", editionName)

// 			log.Debug("Event Router: Getting codec from data")
// 			codec, err = client.GetAudioCodec(data)
// 			if err != nil {
// 				log.Errorf("Event Router: error getting codec, can't continue: %s", err)

// 				return
// 			}
// 		}
// 	}

// 	log.Debugf("Event Router: Received codec: %s", codec)
// 	log.Debugf("Event Router: Got media type of: %s ", payload.Metadata.Type)

// 	switch payload.Event {
// 	// unload BEQ on pause OR stop because I never press stop, just pause and then back.
// 	// play means a new file was started
// 	case "media.play":
// 		log.Debug("Event Router: media.play received")
// 		mediaPlay(client, vip, beqClient, haClient, payload, payload.Metadata.Type, codec, editionName)
// 	case "media.stop":
// 		log.Debug("Event Router: media.stop received")
// 		mediaStop(vip, beqClient, haClient, payload)
// 	case "media.pause":
// 		log.Debug("Event Router: media.pause received")
// 		mediaPause(vip, beqClient, haClient, payload)
// 	// Pressing the 'resume' button actually is media.play thankfully
// 	case "media.resume":
// 		log.Debug("Event Router: media.resume received")
// 		mediaResume(vip, beqClient, haClient, payload, payload.Metadata.Type, codec, editionName)
// 	case "media.scrobble":
// 		log.Debug("Scrobble received")
// 		mediaScrobble(vip)
// 	default:
// 		log.Debugf("Received unsupported event: %s", payload.Event)
// 	}
// }

// get the imdb ID from plex metadata
// func getJellyfin(payload models.PlexWebhookPayload) string {
// 	// try to get IMDB title from plex to save time
// 	for _, model := range payload.Metadata.GUID0 {
// 		if strings.Contains(model.ID, "imdb") {
// 			log.Debugf("Got title ID from plex - %s", model.ID)
// 			return strings.Split(model.ID, "imdb://")[1]
// 		}
// 	}

// 	return ""
// }

// // entry point for background tasks
// func JellyFinWorker(plexChan <-chan models.PlexWebhookPayload, vip *viper.Viper) {
// 	log.Info("JellyFinWorker started")

// 	// block forever until closed so it will wait in background for work
// 	for i := range plexChan {
// 		// determine what to do
// 		eventRouter(plexClient, beqClient, haClient, i, vip)
// 	}
// }
