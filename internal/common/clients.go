package common

import (
	// "github.com/iloveicedgreentea/go-plex/models"
)

type Client interface {
	DoPlaybackAction(action string) error
	GetAudioCodec(payload interface{}) (string, error)
	GetPlexMovieDb(payload interface{}) string
}