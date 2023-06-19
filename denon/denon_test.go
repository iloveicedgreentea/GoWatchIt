package denon

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func setupTest() *DenonClient {
	return NewClient("192.168.88.40", "23")
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
	t.Log(mode)
	assert.NotEmpty(t, mode)

}