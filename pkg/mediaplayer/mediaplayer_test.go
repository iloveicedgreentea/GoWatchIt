// File: mediaplayer/mediaplayer_test.go

package mediaplayer

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig() {
	viper.SetConfigFile("../config_test.yaml")
	viper.ReadInConfig()
}

func TestNewMediaPlayerFactory(t *testing.T) {
	factory := NewMediaPlayerFactory()
	assert.NotNil(t, factory)
	assert.Empty(t, factory.players)
}

func TestRegisterAndGetPlayer(t *testing.T) {
	factory := NewMediaPlayerFactory()
	mockPlayer := struct{}{}

	factory.RegisterPlayer("mock", mockPlayer)
	player, ok := factory.GetPlayer("mock")

	assert.True(t, ok)
	assert.Equal(t, mockPlayer, player)

	_, ok = factory.GetPlayer("nonexistent")
	assert.False(t, ok)
}
