package plex

import (
	"context"

	// "strings"
	"testing"

	// "github.com/StalkR/imdb"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getNewClient() (*PlexClient, error) {
	return NewClient("https", config.GetString("plex.url"), config.GetString("plex.port"))
}

// test to ensure server is white listed
func TestGetPlexReq(t *testing.T) {
	c, err := getNewClient()
	ctx := context.Background()
	require.NoError(t, err)
	d, err := c.makePlexReq(ctx, "/library/metadata/70390")
	assert.NoError(t, err)
	res := string(d)

	if !assert.NotContains(t, res, "Unauthorized", "Client is not authorized in plex server") {
		t.Fatal(err)
	}
}

func TestGetMediaData(t *testing.T) {
	c, err := getNewClient()
	ctx := context.Background()
	require.NoError(t, err)

	// no time to die
	event := models.Event{
		Metadata: models.Metadata{
			Key: "/library/metadata/70390",
		},
	}
	med, err := c.getMediaData(ctx, event)
	assert.NoError(t, err)

	assert.Equal(t, "Atmos", med.Video.Media.AudioCodec)

}
func TestGetCodecFromSession(t *testing.T) {
	t.SkipNow()
	c, err := getNewClient()
	ctx := context.Background()
	require.NoError(t, err)

	codec, err := c.GetCodecFromSession(ctx, config.GetString("plex.deviceuuidfilter"))
	assert.NoError(t, err)

	t.Log(codec)
}

type codecTest struct {
	codec     string
	fullcodec string
	expected  string
}

func TestMapCodecs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	tests := []codecTest{
		{
			codec:     "EAC3",
			fullcodec: "EAC3",
			expected:  "DD+",
		},
		{
			codec:     "EAC3 5.1",
			fullcodec: "German (German EAC3 5.1)",
			expected:  "DD+Atmos5.1Maybe",
		},
		{
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
		s := MapPlexToBeqAudioCodec(ctx, test.codec, test.fullcodec)
		assert.Equal(test.expected, s)
	}

}

// lets me print out every codec I have in a given library
// func getCodecTemp(c *PlexClient, libraryKey string) string {
// 	data, err := c.GetMediaData(ctx, libraryKey)
// 	if err != nil {
// 		return "fail"
// 	}
// 	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
// 	// loop instead of index because of edge case with two or more video streams
// 	for _, val := range data.Video.Media.Part.Stream {
// 		if val.StreamType == "2" {
// 			return fmt.Sprintf("%s --- %s \n", val.DisplayTitle, val.ExtendedDisplayTitle)
// 		}
// 	}

// 	return "fail"
// }

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
