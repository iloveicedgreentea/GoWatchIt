package jellyfin

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetCodec(t *testing.T) {
	// make client
	c := NewClient(
		config.GetString("jellyfin.url"),
		config.GetString("jellyfin.port"),
		"",
		"",
	)
	resp, err := c.GetCodec(config.GetString("jellyfin.userID"), "0a329b45-1faa-b210-c7b2-3aacd4775b1a")
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}