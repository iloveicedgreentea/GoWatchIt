package handlers

// import (

// 	// "github.com/iloveicedgreentea/go-plex/internal/avr"
// 	"github.com/iloveicedgreentea/go-plex/internal/common"
// 	"github.com/iloveicedgreentea/go-plex/internal/plex"
// 	// "github.com/iloveicedgreentea/go-plex/internal/config"
// )

// functions to ensure plex is not being stupid and transcoding atmos for no reason
// I notice it tends to do it RANDOMLY and it is annoying as hell
// so I want to get notified when it happens

// TODO: test and finish mapping
// Maps Denon codecs to BEQ codecs


// // TODO: finish this and generalize
// func isExpectedCodecPlayingPlex(p *plex.PlexClient, uuid string, denonCodec string) (string, bool) {
// 	plexPlaying, err := p.GetCodecFromSession(uuid)
// 	if err != nil {
// 		log.Errorf("Error getting plex audio stream: %s", err)
// 		return "", false
// 	}

// 	// compare the two
// 	log.Error("Expected codec is not playing! Please check your AVR and Client settings!")
// 	// TODO: use IsExpectedCodecPlaying and such
// 	// TODO: use correct AVR mapping
// 	return plexPlaying, mapDenonToBeq(denonCodec) != plex.MapPlexToBeqAudioCodec(plexPlaying, "")
// 	// if enabled, stop playing

// 	// if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
// 	// 	err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %s", m.Codec, expectedCodec), config.GetString("ezbeq.notifyEndpointName"))
// 	// 	if err != nil {
// 	// 		log.Error(err)
// 	// 	}
// 	// }
// }
