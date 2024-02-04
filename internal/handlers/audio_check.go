package handlers

import (

	// "github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/plex"
	// "github.com/iloveicedgreentea/go-plex/internal/config"
)

// functions to ensure plex is not being stupid and transcoding atmos for no reason
// I notice it tends to do it RANDOMLY and it is annoying as hell
// so I want to get notified when it happens

// TODO: test and finish mapping
// Maps Denon codecs to BEQ codecs
func mapDenonToBeq(denonCodec string) string {
	// if False and false, then check others
	switch {
	case common.InsensitiveContains(denonCodec, "dolby atmos"):
		return "Atmos"
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later.
	// Most of the time, TrueHD 7.1 is Atmos
	// TODO: test this
	case common.InsensitiveContains(denonCodec, "dolby hd"):
		return "AtmosMaybe"
	case common.InsensitiveContains(denonCodec, "DOLBY DIGITAL +"):
		return "AtmosMaybe"
	case common.InsensitiveContains(denonCodec, "DTS:X"):
		return "DTS-X"
	// DTS MA 7.1 containers but not DTS:X codecs
	// DTS-HD MSTR
	case common.InsensitiveContains(denonCodec, "DTS-HD MA 7.1") && !common.InsensitiveContains(denonCodec, "DTS:X") && !common.InsensitiveContains(denonCodec, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case common.InsensitiveContains(denonCodec, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case common.InsensitiveContains(denonCodec, "DTS 5.1"):
		return "DTS 5.1"
	// TrueHD 5.1
	case common.InsensitiveContains(denonCodec, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case common.InsensitiveContains(denonCodec, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// DTS HRA
	case common.InsensitiveContains(denonCodec, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case common.InsensitiveContains(denonCodec, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case common.InsensitiveContains(denonCodec, "LPCM 5.1"):
		return "LPCM 5.1"
	case common.InsensitiveContains(denonCodec, "LPCM 7.1"):
		return "LPCM 7.1"
	case common.InsensitiveContains(denonCodec, "LPCM 2.0"):
		return "LPCM 2.0"
	case common.InsensitiveContains(denonCodec, "AAC Stereo"):
		return "AAC 2.0"
	case common.InsensitiveContains(denonCodec, "AC3 5.1") || common.InsensitiveContains(denonCodec, "EAC3 5.1"):
		return "AC3 5.1"
	default:
		return "Empty"
	}

}

// TODO: finish this and generalize
func isExpectedCodecPlayingPlex(p *plex.PlexClient, uuid string, denonCodec string) (string, bool) {
	plexPlaying, err := p.GetCodecFromSession(uuid)
	if err != nil {
		log.Errorf("Error getting plex audio stream: %s", err)
		return "", false
	}

	// compare the two
	log.Error("Expected codec is not playing! Please check your AVR and Client settings!")
	// TODO: use IsExpectedCodecPlaying and such
	// TODO: use correct AVR mapping
	return plexPlaying, mapDenonToBeq(denonCodec) != plex.MapPlexToBeqAudioCodec(plexPlaying, "")
	// if enabled, stop playing

	// if config.GetBool("ezbeq.notifyOnLoad") && config.GetBool("homeAssistant.enabled") {
	// 	err := haClient.SendNotification(fmt.Sprintf("Wrong codec is playing. Expected codec %s but got %s", m.Codec, expectedCodec), config.GetString("ezbeq.notifyEndpointName"))
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// }
}
