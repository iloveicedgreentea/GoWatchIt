package plex

import (
	"context"
	"fmt"
	"net/http"

	"strings"
	"sync"
	"time"

	"encoding/json"
	"errors"

	"github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
	"github.com/iloveicedgreentea/go-plex/internal/plex"
	"github.com/iloveicedgreentea/go-plex/models"
	"golang.org/x/exp/slices"
)

// decodes the multipart form from plex
// TODO test this
func parsePlexMultipartForm(payload []string) (models.PlexWebhookPayload, error) {
	var pwhPayload models.PlexWebhookPayload

	err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// unmarshall error
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return pwhPayload, errors.New(msg)

		default:
			return pwhPayload, err
		}
	}

	log.Debugf("decodeWebhook: Received event: %s", pwhPayload.Event)
	return pwhPayload, nil
}

// TODO: test this
// getMultipartPayload gets the payload from the multipart form and returns if ok
func getMultipartPayload(request *http.Request) ([]string, bool, error) {
	if err := request.ParseMultipartForm(0); err != nil {
		log.Errorf("invalid multipart form: %s", err)
		return []string{}, false, fmt.Errorf("invalid multipart form: %s", err)
	}

	payload, ok := request.MultipartForm.Value["payload"]
	return payload, ok, nil
}

// Sends the payload to the channel for background processing
// TODO: this takes an HTTP request and sends it to a channel
func (p *PlexPlayer) ProcessPlexWebhook(ctx context.Context, request *http.Request) error {
	payload, ok, err := getMultipartPayload(request)
	if err != nil {
		return fmt.Errorf("error getting payload: %s", err)
	}
	if ok {
		// parse the payload
		log.Debug("decoding payload")
		decodedPayload, err := parsePlexMultipartForm(payload)
		if err != nil {
			log.Error(err)
			return fmt.Errorf("error decoding payload: %s", err)
		}

		log.Debugf("Got a request from UUID: %s", decodedPayload.Player.UUID)

		mediaType := decodedPayload.Metadata.Type

		log.Debugf("ProcessWebhook:  Media type is: %s", mediaType)
		log.Debugf("ProcessWebhook:  Media title is: %s", decodedPayload.Metadata.Title)

		// check filter for user if not blank
		userID := config.GetString("plex.ownerNameFilter")
		// only respond to events on a particular account if you share servers and only for movies and shows
		// TODO: decodedPayload.Account.Title seems to always map to server owner not player account
		if len(userID) == 0 || strings.EqualFold(decodedPayload.Account.Title, userID) {
			if strings.EqualFold(mediaType, string(movieItemTitle)) || strings.EqualFold(mediaType, string(showItemTitle)) {
				log.Debug("adding item to plexChan")
				select {
				case p.plexChan <- decodedPayload:
					// send succeeded
					return nil
				// context was cancelled
				case <-ctx.Done():
					log.Error("Processing was cancelled")
					return ctx.Err()
				case <-time.After(time.Second * 3):
					log.Error("Send on plexChan timed out")
					return errors.New("Send on plexChan timed out")
				}
			} else {
				log.Debugf("Media type of %s is not supported", mediaType)
			}
		} else {
			// TODO: this seems to be hitting even when the filter matches
			log.Debugf("userID '%s' does not match filter of %s", decodedPayload.Account.Title, userID)
		}
	} else {
		log.Error("No payload found in request")
		return errors.New("no payload found in request")
	}
	return nil
}

// does plex send stop if you exit with back button? - Yes, with X for mobile player as well
func (p *PlexPlayer) mediaStop(cancel context.CancelFunc, payload models.PlexWebhookPayload) {
	// cancel mediaPlay and resume
	cancel()
	// TODO: move these functions somewhere
	go common.ChangeLight("on")
	// TODO: deprecate mqtt, send event
	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
	if err != nil {
		log.Error(err)
	}

	err = p.BeqClient.UnloadBeqProfile(p.SearchRequest)
	if err != nil {
		log.Error(err)
		if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
			err := p.HaClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err))
			if err != nil {
				log.Error()
			}
		}
	}
	log.Info("BEQ profile unloaded")
}

// pause only happens with literally pausing
func (p *PlexPlayer) mediaPause(cancel context.CancelFunc, payload models.PlexWebhookPayload) {
	// skip processing webhooks since HDMI sync will send pause and resume
	if !*p.skipActions {
		// cancel other running functions
		cancel()
		go common.ChangeLight("on")

		err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
		if err != nil {
			log.Error(err)
		}

		err = p.BeqClient.UnloadBeqProfile(m)
		if err != nil {
			log.Error(err)
			if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
				err := p.HaClient.SendNotification(fmt.Sprintf("Error UNLOADING profile: %v -- Unsafe to play movies!", err))
				if err != nil {
					log.Error()
				}
			}
		}
		log.Info("BEQ profile unloaded")
	}
}

// play is both the "resume" UI button and play
func (p *PlexPlayer) mediaPlay(ctx context.Context, avrClient avr.AVRClient, payload models.PlexWebhookPayload, useAvrCodec bool, data models.MediaContainer, wg *sync.WaitGroup) {
	// Check if context is already cancelled before starting lets say you play but then stop, this should stop processing
	if ctx.Err() != nil {
		log.Debug("mediaPlay was called with a cancelled context")
		return
	}

	var err error
	go func() {
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before lights and volume change")
			return
		}
		common.ChangeLight("off")
		common.ChangeMasterVolume(p.SearchRequest.MediaType)
	}()

	// if its not a movie and time is the source, skip hdmi sync
	if !strings.EqualFold(payload.Metadata.Type, string(movieItemTitle)) && config.GetString("signal.source") == "time" {
		log.Debug("skipping sync for non-movie type and time source")
	} else {
		wg.Add(1)
		go func() {
			if ctx.Err() != nil {
				log.Debug("mediaPlay was cancelled before hdmi sync")
				return // Exit early if context is cancelled
			}

			// optimistically try to hdmi sync. Will return if disabled
			// TODO: implement this
			// common.WaitForHDMISync(wg, skipActions, haClient, client)
		}()
	}

	// dont need to set skipActions here because it will only send media.pause and media.resume. This is media.play

	go func() {
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before publishing playing status")
			return
		}
		if err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true"); err != nil {
			log.Error("Error publishing playing status: ", err)
		}
	}()

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before unloading BEQ profile")
		return
	}
	if err = p.BeqClient.UnloadBeqProfile(p.SearchRequest); err != nil {
		log.Errorf("Error unloading BEQ during play: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Error("mediaPlay cancelled before unloading BEQ profile")
		return
	default:
		log.Debug("Using plex to get codec")
		// TODO: try session data then fallback to lookup
		p.SearchRequest.Codec, err = p.PlexClient.GetAudioCodec(data)
		if err != nil {
			log.Errorf("error getting codec from plex, can't continue: %s", err)
			return
		}
		// slower but more accurate especially with atmos
		// TODO: implement avr stuff
		// if useAvrCodec {
		// 	p.SearchRequest.Codec, err = checkAvrCodec(client, haClient, avrClient, payload, data)
		// 	// if it failed, get codec data from client
		// 	if err != nil {
		// 		log.Warnf("error getting codec from AVR, falling back to client: %s", err)
		// 		m.Codec, err = client.GetAudioCodec(data)
		// 		if err != nil {
		// 			log.Errorf("error getting codec from plex, can't continue: %s", err)
		// 			return
		// 		}
		// 	}
		// } else {
		// 	log.Debug("Using plex to get codec")
		// 	// TODO: try session data then fallback to lookup
		// 	m.Codec, err = client.GetAudioCodec(data)
		// 	if err != nil {
		// 		log.Errorf("error getting codec from plex, can't continue: %s", err)
		// 		return
		// 	}
		// }
		log.Debugf("Found codec: %s", m.Codec)
		// if its a show and you dont want beq enabled, exit
		if strings.EqualFold(payload.Metadata.Type, string(showItemTitle)) {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}

		p.SearchRequest.TMDB = getPlexMovieDb(payload)
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled before loading BEQ profile")
		return
	}
	err = p.BeqClient.LoadBeqProfile(m)
	if err != nil {
		if err.Error() == "beq profile was not found in catalog" {
			log.Warnf("BEQ profile was not found in the catalog. Either the metadata is wrong or this %s does not have a BEQ", payload.Metadata.Type)
			return
		} else {
			log.Error("Error loading BEQ profile: ", err)
			return
		}
	}
	log.Info("BEQ profile loaded")
	// send notification of it loaded
	if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
		err := p.HaClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, p.SearchRequest.Codec))
		if err != nil {
			log.Error()
		}
	}

	if ctx.Err() != nil {
		log.Debug("mediaPlay was cancelled at a later stage")
		return
	}

	log.Debug("Waiting for goroutines")
	wg.Wait()
	log.Debug("Goroutines complete")
}

// resume is only after pausing as long as the media item is still active
func (p *PlexPlayer) mediaResume(ctx context.Context, payload models.PlexWebhookPayload, data models.MediaContainer) {
	if !*p.skipActions {
		if ctx.Err() != nil {
			log.Debug("mediaResume was called with a cancelled context")
			return
		}
		// mediaType string, codec string, edition string
		// trigger lights
		go common.ChangeLight("off") // TODO: split stuff like this into functions
		err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "true")
		if err != nil {
			log.Error(err)
		}
		// Changing on resume is disabled because its annoying if you changed it since playing
		// go changeMasterVolume(vip, mediaType)

		// allow skipping search to save time
		// always unload in case something is loaded from movie for tv

		// TODO: make all of this a function
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before unloading BEQ profile")
			return
		}
		err = p.BeqClient.UnloadBeqProfile(p.SearchRequest)
		if err != nil {
			log.Errorf("Error on startup - unloading beq %v", err)
		}
		if strings.EqualFold(payload.Metadata.Type, string(showItemTitle)) {
			if !config.GetBool("ezbeq.enableTvBeq") {
				return
			}
		}
		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before getting plex data")
			return
		}
		// get the tmdb id to match with ezbeq catalog
		p.SearchRequest.TMDB = getPlexMovieDb(payload)
		// if the server was restarted, cached data is lost
		if len(m.Codec) == 0 {
			log.Warn("No codec found in cache on resume. Was server restarted? Getting new codec")
			log.Debug("Using plex to get codec because its not cached")
			p.SearchRequest.Codec, err = p.PlexClient.GetAudioCodec(data)
			if err != nil {
				log.Errorf("error getting codec from plex, can't continue: %s", err)
				return
			}
		}

		if ctx.Err() != nil {
			log.Debug("mediaPlay was cancelled before loading BEQ profile")
			return
		}

		err = p.BeqClient.LoadBeqProfile(p.SearchRequest)
		if err != nil {
			if err.Error() == "beq profile was not found in catalog" {
				log.Warnf("BEQ profile was not found in the catalog. Either the metadata is wrong or this %s does not have a BEQ", payload.Metadata.Type)
				return
			} else {
				log.Error("Error loading BEQ profile: ", err)
				return
			}
		}
		log.Info("BEQ profile loaded")

		// send notification of it loaded
		if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
			err := p.HaClient.SendNotification(fmt.Sprintf("BEQ Profile: Title - %s  (%d) // Codec %s", payload.Metadata.Title, payload.Metadata.Year, p.SearchRequest.Codec))
			if err != nil {
				log.Error()
			}
		}
	}
}

func (p *PlexPlayer) mediaScrobble() {
	// trigger lights
	// go changeLight(vip, "on")
	err := mqtt.PublishWrapper(config.GetString("mqtt.topicplayingstatus"), "false")
	if err != nil {
		log.Error(err)
	}
	log.Debug("Scrobble received. Unloading profile")
	// unload beq
	err = p.BeqClient.UnloadBeqProfile(p.SearchRequest)
	if err != nil {
		log.Errorf("Error on startup - unloading beq %v", err)
	}

}

func mapToEdition(s string) models.Edition {
	switch {
	case strings.Contains(s, "extended"):
		return models.EditionExtended
	case strings.Contains(s, "unrated"):
		return models.EditionUnrated
	case strings.Contains(s, "theatrical"):
		return models.EditionTheatrical
	case strings.Contains(s, "ultimate"):
		return models.EditionUltimate
	case strings.Contains(s, "director"):
		return models.EditionDirectorsCut
	case strings.Contains(s, "criterion"):
		return models.EditionCriterion
	default:
		return models.EditionUnknown
	}
}

// getEditionName tries to extract the edition from plex or file name. Assumes you have well named files
// Returned types, Unrated, Ultimate, Theatrical, Extended, Director, Criterion
// TODO: create union type for data and others
func (p *PlexPlayer) getEditionName(ctx context.Context ,container models.DataMediaContainer) (models.Edition, error) {
	data := container.PlexPayload
	if data == nil {
		return "", errors.New("data is nil")
	}
	edition := strings.ToLower(data.Video.EditionTitle)
	fileName := strings.ToLower(data.Video.Media.Part.File)

	// First, check the edition from Plex metadata
	if edition != "" {
		mappedEdition := mapToEdition(edition)
		if mappedEdition != "" {
			return mappedEdition, nil
		}
		// If we couldn't map it, return it unknown
		return models.EditionUnknown, errors.New("could not map edition")
	}

	// If no edition in metadata, try to extract from file name
	mappedEdition := mapToEdition(fileName)
	if mappedEdition != "" {
		return mappedEdition, nil
	}

	// no edition found, so its standard
	return models.EditionNone, nil
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
func eventRouter(ctx context.Context, cancel context.CancelFunc, plexClient *plex.PlexClient, beqClient *ezbeq.BeqClient, haClient *homeassistant.HomeAssistantClient, avrClient avr.AVRClient, useAvrCodec bool, payload models.PlexWebhookPayload, model *models.SearchRequest, skipActions *bool) {
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
		wg := &sync.WaitGroup{}
		mediaPlay(ctx, plexClient, beqClient, haClient, avrClient, payload, model, useAvrCodec, data, skipActions, wg)
	case "media.stop":
		log.Debug("Event Router: media.stop received")
		mediaStop(cancel, beqClient, haClient, payload, model)
	case "media.pause":
		log.Debug("Event Router: media.pause received")
		mediaPause(cancel, beqClient, haClient, payload, model, skipActions)
	// Pressing the 'resume' button in plex UI is media.play
	case "media.resume":
		log.Debug("Event Router: media.resume received")
		mediaResume(ctx, plexClient, beqClient, haClient, payload, model, data, skipActions)
	case "media.scrobble":
		log.Debug("Scrobble received")
		mediaScrobble(beqClient, model)
	default:
		log.Debugf("Received unsupported event: %s", payload.Event)
	}
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
	var ctx context.Context
	var cancel context.CancelFunc
	// block forever until closed so it will wait in background for work
	for i := range plexChan {
		// TODO: this means every request will have its only context and wont be cancellable
		ctx, cancel = context.WithCancel(context.Background())
		log.Debugf("Current length of plexChan in PlexWorker: %d", len(plexChan))
		// determine what to do
		log.Debug("Sending new payload to eventRouter")
		eventRouter(ctx, cancel, plexClient, beqClient, haClient, avrClient, useAvrCodec, i, model, skipActions)
		log.Debug("eventRouter done processing payload")
	}
	cancel()
	log.Debug("Plex worker stopped")
}
