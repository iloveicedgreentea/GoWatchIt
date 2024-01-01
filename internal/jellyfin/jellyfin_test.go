package jellyfin

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

func testSetup() (*JellyfinClient) {
	c := NewClient(
		config.GetString("jellyfin.url"),
		config.GetString("jellyfin.port"),
		"",
		"",
	)
	return c
}

func getMetadata() (models.JellyfinMetadata, error) {
	c := testSetup()
	return c.GetMetadata(config.GetString("jellyfin.userID"), "0a329b45-1faa-b210-c7b2-3aacd4775b1a")
}
func TestGetCodec(t *testing.T) {
	// make client
	c := testSetup()
	m, err := getMetadata()
	assert.NoError(t, err)
	codec, displayTitle, err := c.GetCodec(m)
	assert.NoError(t, err)
	t.Log(codec, displayTitle)
}
func TestGetEdition(t *testing.T) {
	// make client
	c := testSetup()
	m, err := getMetadata()
	assert.NoError(t, err)
	edition := c.GetEdition(m)
	assert.Equal(t, "", edition)
}