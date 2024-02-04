package handlers

import (
	"strings"

	// "github.com/iloveicedgreentea/go-plex/internal/avr"
	// "github.com/iloveicedgreentea/go-plex/internal/plex"
)

// functions to ensure plex is not being stupid and transcoding atmos for no reason
// I notice it tends to do it RANDOMLY and it is annoying as hell
// so I want to get notified when it happens


// TODO: map denon to beq
// TODO: test
func mapDenonToBeq(denonCodec string) string {
	// if False and false, then check others
	switch {
	// There are very few truehd 7.1 titles and many atmos titles have wrong metadata. This will get confirmed later
	case strings.Contains(denonCodec, "dolby hd"):
		return "AtmosMaybe"
	case strings.Contains(denonCodec, "DTS:X"):
		return "DTS-X"
	// DTS MA 7.1 containers but not DTS:X codecs
	case strings.Contains(denonCodec, "DTS-HD MA 7.1") && !strings.Contains(denonCodec, "DTS:X") && !strings.Contains(denonCodec, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case strings.Contains(denonCodec, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case strings.Contains(denonCodec, "DTS 5.1"):
		return "DTS 5.1"
	// TrueHD 5.1
	case strings.Contains(denonCodec, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case strings.Contains(denonCodec, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// DTS HRA
	case strings.Contains(denonCodec, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case strings.Contains(denonCodec, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case strings.Contains(denonCodec, "LPCM 5.1"):
		return "LPCM 5.1"
	case strings.Contains(denonCodec, "LPCM 7.1"):
		return "LPCM 7.1"
	case strings.Contains(denonCodec, "LPCM 2.0"):
		return "LPCM 2.0"
	case strings.Contains(denonCodec, "AAC Stereo"):
		return "AAC 2.0"
	case strings.Contains(denonCodec, "AC3 5.1") || strings.Contains(denonCodec, "EAC3 5.1"):
		return "AC3 5.1"
	default:
		return "Empty"
	}

}
