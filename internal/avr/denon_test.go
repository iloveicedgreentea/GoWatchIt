package avr

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupTest() *DenonClient {
	return &DenonClient{
		ServerURL: "192.168.88.40",
		Port: "23",
	}
}
func TestMakeReq(t *testing.T) {
	c := setupTest()

	res, err := c.makeReq("PW?")
	assert.NoError(t, err)
	assert.Equal(t, "PWSTANDBY\r", res)

}
func TestGetAudioMode(t *testing.T) {
	c := setupTest()

	mode, err := c.GetCodec()
	assert.NoError(t, err)
	t.Logf("avr codec: %s", mode)
	assert.NotEqual(t, "Empty",mode)


}

func TestNewClient(t *testing.T) {
	c := GetAVRClient(config.GetString("ezbeq.avrUrl"))
	assert.NotNil(t, c)
}