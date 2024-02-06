package common

import (
	// "github.com/iloveicedgreentea/go-plex/models"
)

func PlaybackInterface(action string, c Client) error {
	return c.DoPlaybackAction(action)
}

