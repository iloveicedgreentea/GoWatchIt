package common

import (
	// "github.com/iloveicedgreentea/go-plex/models"
)

// A Client is an interface for interacting with a player like Plex or Jellyfin
type Client interface {
	DoPlaybackAction(action string) error
	GetAudioCodec(payload interface{}) (string, error)
	// TODO: generalize this
	GetPlexMovieDb(payload interface{}) string
}