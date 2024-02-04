package avr

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/iloveicedgreentea/go-plex/internal/config"
)

func setupAvrTest() AVRClient {
	
	return GetAVRClient("192.168.88.40")
}
func TestAvrGetAudioMode(t *testing.T) {
	config.Set("ezbeq.avrbrand", "denon")
	c := setupAvrTest()

	mode, err := c.GetCodec()
	assert.NoError(t, err)
	t.Log(mode)
	assert.NotEmpty(t, mode)

}
func TestAvrGetAudioModeFail(t *testing.T) {
	config.Set("ezbeq.avrbrand", "")
	c := setupAvrTest()

	assert.Nil(t, c)

}