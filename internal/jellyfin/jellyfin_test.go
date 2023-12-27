package jellyfin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCodec(t *testing.T) {
	// make client
	c := NewClient(
		"http://",
		"32400",
		"",
		"",
	)
	resp, err := c.GetCodec("abc", "123")
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

