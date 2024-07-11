// File: plex/plex_test.go

package plex

import (
	"context"
	"testing"
	"time"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const https = "https://"

func newClient() (*PlexPlayer, error) {
	serverURL := viper.GetString("plex.url")
	port := viper.GetString("plex.port")
	c := make(chan<- models.PlexWebhookPayload)

	return NewPlexPlayer(https, serverURL, port, nil, nil, c)
}

func TestNewPlexPlayer(t *testing.T) {
	player, err := newClient()
	require.NoError(t, err)

	assert.NotNil(t, player)
	assert.NotNil(t, player.PlexClient)
}

// func TestPlexPlayer_Play(t *testing.T) {

// 	serverURL := viper.GetString("plex.url")
// 	port := viper.GetString("plex.port")

// 	player := NewPlexPlayer(serverURL, port)
// 	err := player.Play(context.Background())
// 	require.NoError(t, err)
// }

func TestPlexPlayerAction(t *testing.T) {

	player, err := newClient()
	require.NoError(t, err)

	ctx := context.Background()
	err = player.Pause(ctx)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	err = player.Play(ctx)
	require.NoError(t, err)
}

// func TestPlexPlayer_RouteEvent(t *testing.T) {

// 	serverURL := viper.GetString("plex.url")
// 	port := viper.GetString("plex.port")

// 	player := NewPlexPlayer(serverURL, port)

// 	tests := []struct {
// 		name      string
// 		eventType string
// 		wantErr   bool
// 	}{
// 		{"Play Event", "media.play", false},
// 		{"Pause Event", "media.pause", false},
// 		{"Unknown Event", "unknown.event", true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := player.RouteEvent(context.Background(), tt.eventType, models.MediaPayload{})
// 			if tt.wantErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestPlexPlayer_GetMediaData(t *testing.T) {

// 	serverURL := viper.GetString("plex.url")
// 	port := viper.GetString("plex.port")

// 	player := NewPlexPlayer(serverURL, port)

// 	// You'll need to provide a valid key for an actual item in your Plex server
// 	key := "/library/metadata/1234"

// 	data, err := player.GetMediaData(context.Background(), key)
// 	require.NoError(t, err)
// 	require.NotNil(t, data)
// 	// Add more specific assertions based on what you expect in the MediaContainer
// }
