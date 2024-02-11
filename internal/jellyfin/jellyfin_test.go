package jellyfin

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

func testSetup() *JellyfinClient {
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
	m, err := getMetadata("76b466edcad9642a707201ecf1dfdf96")
	assert.NoError(t, err)
	codec, err := c.GetAudioCodec(m)
	assert.NoError(t, err)
	assert.Equal(t, "Atmos", codec)
}

func TestGetEdition(t *testing.T) {
	// make client
	c := testSetup()
	m, err := getMetadata("0a329b45-1faa-b210-c7b2-3aacd4775b1a")
	assert.NoError(t, err)
	edition := c.GetEdition(m)
	// TODO: test with known edition
	assert.Equal(t, "", edition)
	t.Logf("%#v", m)
}

// used for testing
func TestPrintCodec(t *testing.T) {
	// make client
	c := testSetup()
	m, err := getMetadata("1a1831da4f875cc5df09507fb49d2877")
	assert.NoError(t, err)
	codec, err := c.GetAudioCodec(m)
	assert.NoError(t, err)
	t.Log(codec)
}
func TestPrintMetadata(t *testing.T) {
	// make client
	m, err := getMetadata("1a1831da4f875cc5df09507fb49d2877")
	assert.NoError(t, err)
	t.Log(m.OriginalTitle)
}

type codecTest struct {
	codec        string
	displayTitle string
	profile      string
	layout       string
	expected     string
}

func TestMapCodecs(t *testing.T) {
	assert := assert.New(t)
	tests := []codecTest{
		// dd+
		{
			codec:        "EAC3",
			displayTitle: "English - Dolby Digital+ - Stereo - Default",
			profile:      "",
			expected:     "DD+",
		},
		// dd+ //TODO: revisit this find a dd5.1 title only
		// {
		// 	codec:        "EAC3",
		// 	displayTitle: "English - Dolby Digital+ - 5.1 - Default",
		// 	profile:      "",
		// 	expected:     "DD+",
		// },
		// ac3/DD
		{
			codec:        "AC3",
			displayTitle: "English - Dolby Digital - 5.1 - Default",
			profile:      "",
			layout:       "5.1",
			expected:     "AC3 5.1",
		},
		// DD+ Atmos old guard
		{
			codec:        "eac3",
			displayTitle: "English - Dolby Digital+ - 5.1 - Default",
			layout:       "5.1",
			expected:     "DD+ Atmos",
		},
		// dts-hd ma 5.1
		{
			codec:        "dts",
			displayTitle: "DTS-HD MA 5.1 - English - Default",
			profile:      "DTS-HD MA",
			expected:     "DTS-HD MA 5.1",
		},
		//dts-hd ma 5.1
		{
			codec:        "DTS",
			displayTitle: "Surround 5.1 - English - DTS-HD MA - Default",
			layout:       "5.1",
			expected:     "DTS-HD MA 5.1",
		},
		//dts-hd ma 5.1
		{
			codec:        "DTS",
			displayTitle: "DTS-HD Master Audio / 5.1 / 48 kHz / 2928 kbps / 24-bit - English - Default",
			layout:       "5.1",
			profile:      "DTS-HD MA",
			expected:     "DTS-HD MA 5.1",
		},
		//dts-hd ma 7.1
		{
			codec:        "dts",
			displayTitle: "English - DTS-HD MA - 7.1 - Default",
			profile:      "DTS-HD MA",
			layout:       "7.1",
			expected:     "DTS-HD MA 7.1",
		},
		// dts-x
		{
			codec:        "DTS",
			displayTitle: "DTS-X 7.1 - English - DTS-HD MA - Default",
			profile:      "DTS-HD MA",
			layout:       "7.1",
			expected:     "DTS-X",
		},
		// dts hd HRA fast 6
		{
			codec:        "DTS",
			displayTitle: "Surround 7.1 - English - DTS-HD HRA - Default",
			profile:      "DTS-HD HRA",
			layout:       "7.1",
			expected:     "DTS-HD HR 7.1",
		},
		// : dts 5.1
		{
			codec:        "DTS",
			displayTitle: "Surround 5.1 - English - Default",
			profile:      "",
			layout:       "5.1",
			expected:     "DTS 5.1",
		},
		// truehd 5.1
		{
			codec:        "truehd",
			displayTitle: " Dolby TrueHD Audio / 5.1 / 48 kHz / 16-bit (AC3 Embedded: 5.1 / 48 kHz / 640 kbps) - English - Default ",
			profile:      "",
			layout:       "5.1",
			expected:     "TrueHD 5.1",
		},
		// truehd 7.1 ghost protocl
		{
			codec:        "truehd",
			displayTitle: " Dolby TrueHD Audio / 7.1 / 48 kHz / 16-bit (AC3 Embedded: 5.1 / 48 kHz / 640 kbps) - English - Default ",
			profile:      "",
			layout:       "7.1",
			expected:     "AtmosMaybe",
		},
		// house of dragon -  truehd 7.1 from JF but its atmos
		{
			codec:        "truehd",
			displayTitle: "Surround - English - TRUEHD - 7.1 - Default",
			profile:      "",
			layout:       "7.1",
			expected:     "AtmosMaybe",
		},
		{
			codec:        "DTS",
			displayTitle: "DTS-HD Master Audio / 5.1 / 48 kHz / 3186 kbps / 24-bit - English - Default",
			profile:      "DTS-HD MA",
			expected:     "DTS-HD MA 5.1",
		},
		// loki
		{
			codec:        "TRUEHD",
			displayTitle: "TrueHD Atmos 7.1 - English - Default",
			profile:      "",
			expected:     "Atmos",
		},
	}
	// execute each test
	for _, test := range tests {
		test := test
		s := MapJFToBeqAudioCodec(test.codec, test.displayTitle, test.profile, test.layout)
		assert.Equal(test.expected, s, "Test failed for %v, got %s", test, s)
	}
}

func TestGetTMDB(t *testing.T) {
	c := testSetup()
	metadata, err := c.GetMetadata(config.GetString("jellyfin.userID"), "1efb0048dd138e93771ab59ab85c03f1")
	assert.NoError(t, err)
	tmdb, err := c.GetJfTMDB(metadata)
	assert.NoError(t, err)
	assert.Equal(t, "56292", tmdb)
}