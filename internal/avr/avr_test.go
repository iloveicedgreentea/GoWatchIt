package avr

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func setupAvrTest() AVRClient {
	return GetAVRClient("denon", "192.168.88.40", "23")
}
func TestAvrGetAudioMode(t *testing.T) {
	c := setupAvrTest()

	mode, err := c.GetCodec()
	assert.NoError(t, err)
	t.Log(mode)
	assert.NotEmpty(t, mode)

}