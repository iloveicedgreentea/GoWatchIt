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

func isExpectedCodecPlaying(c denon.DenonClient) (bool, error) {
	denonPlaying, err := c.GetAudioMode()
	if err != nil {
		return false, err
	}

	plexPlaying, err := getPlexAudioStream()
	if err != nil {
		return false, err
	}

	// compare the two
	if denonPlaying != plexPlaying {
		return false, nil
	}

	return true, nil
}