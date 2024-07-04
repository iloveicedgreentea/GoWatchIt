// File: plex/plex_test.go

package plex

import (
	"context"
	"testing"

	// "github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlexPlayer(t *testing.T) {

	serverURL := viper.GetString("plex.url")
	port := viper.GetString("plex.port")

	player := NewPlexPlayer(serverURL, port)
	require.NotNil(t, player)
	require.NotNil(t, player.plexClient)
}

// func TestPlexPlayer_Play(t *testing.T) {

// 	serverURL := viper.GetString("plex.url")
// 	port := viper.GetString("plex.port")

// 	player := NewPlexPlayer(serverURL, port)
// 	err := player.Play(context.Background())
// 	require.NoError(t, err)
// }

func TestPlexPlayerAction(t *testing.T) {

	serverURL := viper.GetString("plex.url")
	port := viper.GetString("plex.port")
	t.Logf("serverURL: %s, port: %s", serverURL, port)

	player := NewPlexPlayer(serverURL, port)
	err := player.Play(context.Background())
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
