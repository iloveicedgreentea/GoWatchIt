package jellyfin

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/common"
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

func getMetadata(itemID string) (models.JellyfinMetadata, error) {
	c := testSetup()
	return c.GetMetadata(config.GetString("jellyfin.userID"), itemID)
}
func TestGetCodec(t *testing.T) {
	// make client
	c := testSetup()
	m, err := getMetadata("0a329b45-1faa-b210-c7b2-3aacd4775b1a")
	assert.NoError(t, err)
	codec, displayTitle, _, err := c.GetCodec(m)
	assert.NoError(t, err)
	t.Log(codec, displayTitle)
}
func TestGetEdition(t *testing.T) {
	// TODO: test this
	// make client
	c := testSetup()
	m, err := getMetadata("0a329b45-1faa-b210-c7b2-3aacd4775b1a")
	assert.NoError(t, err)
	edition := c.GetEdition(m)
	assert.Equal(t, "", edition)
}

type codecTest struct {
	codec     string
	fullcodec string
	expected  string
}

func TestInsensitiveContains(t *testing.T) {
	assert := assert.New(t)
	assert.True(common.InsensitiveContains("DTS-HD MA 5.1 - English - Default", "DTS-HD MA 5.1"))

}
func TestMapCodecs(t *testing.T) {
	assert := assert.New(t)
	tests := []codecTest{
		{
			codec:     "EAC3",
			fullcodec: "EAC3",
			expected:  "DD+",
		},
		{
			codec:     "AC3",
			fullcodec: "English - Dolby Digital - 5.1 - Default",
			expected:  "AC3 5.1",
		},
		// TODO: add other cases from JF
		{
			codec:     "EAC3 5.1",
			fullcodec: "German (German EAC3 5.1)",
			expected:  "DD+ Atmos",
		},
		{
			codec:     "DTS",
			fullcodec: "DTS-HD MA 5.1 - English - Default",
			expected:  "DTS-HD MA 5.1",
		},
		{
			codec:     "DTS",
			fullcodec: "Surround 5.1 - English - DTS-HD MA - Default",
			expected:  "DTS-HD MA 5.1",
		},
		{ // TODO: when DTS, look at profile
			codec:     "DDP 5.1 Atmos",
			fullcodec: "DDP 5.1 Atmos (Engelsk EAC3)",
			expected:  "DD+ Atmos",
		},
		{
			codec:     "English (TRUEHD 7.1)",
			fullcodec: "Surround 7.1 (English TRUEHD)",
			expected:  "AtmosMaybe",
		},
		{
			codec:     "English (TRUEHD 5.1)",
			fullcodec: "Dolby TrueHD Audio / 5.1 / 48 kHz / 1541 kbps / 16-bit (English)",
			expected:  "TrueHD 5.1",
		},
		{
			codec:     "English (DTS-HD MA 5.1)",
			fullcodec: "DTS-HD Master Audio / 5.1 / 48 kHz / 3887 kbps / 24-bit (English)",
			expected:  "DTS-HD MA 5.1",
		},
		{
			codec:     "English (TRUEHD 7.1)",
			fullcodec: "TrueHD Atmos 7.1 (English)",
			expected:  "Atmos",
		},
		{
			codec:     "English (DTS-HD MA 7.1)",
			fullcodec: "DTS:X / 7.1 / 48 kHz / 4213 kbps / 24-bit (English DTS-HD MA)",
			expected:  "DTS-X",
		},
		// TODO: verify other codecs without using extended display title
	}
	// execute each test
	for _, test := range tests {
		s := MapJFToBeqAudioCodec(test.codec, test.fullcodec)
		assert.Equal(test.expected, s)
	}
}