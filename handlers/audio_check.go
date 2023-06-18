package handlers

import (
	"github.com/iloveicedgreentea/go-plex/denon"
)

// functions to ensure plex is not being stupid and transcoding atmos for no reason
// I notice it tends to do it RANDOMLY and it is annoying as hell
// so I want to get notified when it happens

// func which gets current plex audio stream
// TODO: normalize this with denon/denon.go
func getPlexAudioStream() (string, error) {
	// get file playing
	// determine the codec like truehd, atmos, etc

	return "", nil
}

func isExpectedCodecPlaying(c *denon.DenonClient) bool {
	denonPlaying, err := c.GetAudioMode()
	if err != nil {
		log.Errorf("Error getting denon audio mode: %s", err)
		return false
	}

	plexPlaying, err := getPlexAudioStream()
	if err != nil {
		log.Errorf("Error getting plex audio stream: %s", err)
		return false
	}

	// compare the two
	if denonPlaying != plexPlaying {
		return false
	}

	return true
}
