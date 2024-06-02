package plex

import (
	"fmt"

	// "strings"
	"testing"

	// "github.com/StalkR/imdb"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/stretchr/testify/assert"
)

// test to ensure server is white listed
func TestGetPlexReq(t *testing.T) {
	c := NewClient(config.GetString("plex.url"), config.GetString("plex.port"))
	d, err := c.makePlexReq("/library/metadata/70390")
	assert.NoError(t, err)
	res := string(d)

	if !assert.NotContains(t, res, "Unauthorized", "Client is not authorized in plex server") {
		t.Fatal(err)
	}
}

func TestGetMediaData(t *testing.T) {
	c := NewClient(config.GetString("plex.url"), config.GetString("plex.port"))

	// no time to die
	med, err := c.GetMediaData("/library/metadata/70390")
	assert.NoError(t, err)

	code, err := c.GetAudioCodec(med)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	assert.Equal(t, "Atmos", code)

}
func TestGetCodecFromSession(t *testing.T) {
	t.SkipNow()
	c := NewClient(config.GetString("plex.url"), config.GetString("plex.port"))

	codec, err := c.GetCodecFromSession(config.GetString("plex.deviceuuidfilter"))
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
		s := MapPlexToBeqAudioCodec(test.codec, test.fullcodec)
		assert.Equal(test.expected, s)
	}

}


// For dev only - gets a list of every audio codec present in library
func TestGetPlexMovies(t *testing.T) {
	t.Skip()
	// edit if your movie lib is different
	librarySectionID := "1"
	c := NewClient(config.GetString("plex.url"), config.GetString("plex.port"))

	data, err := c.makePlexReq(fmt.Sprintf("/library/sections/%s/all", librarySectionID))
	if err != nil {
		t.Fatal(err)
	}

	model, err := parseAllMediaContainer(data)
	if err != nil {
		t.Fatal(err)
	}

	var allMovies []string
	// get all movies by their key
	for _, val := range model.Video {
		allMovies = append(allMovies, val.Key)
	}

	var movieCodecs []string
	// parse each key and get the codec
	for _, movie := range allMovies {
		codec := getCodecTemp(c, movie)
		movieCodecs = append(movieCodecs, codec)
	}

	// remove duplicates
	finalList := removeDuplicateStr(movieCodecs)
	t.Log(finalList)
}

// lets me print out every codec I have in a given library
func getCodecTemp(c *PlexClient, libraryKey string) string {
	data, err := c.GetMediaData(libraryKey)
	if err != nil {
		return "fail"
	}
	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
	// loop instead of index because of edge case with two or more video streams
	for _, val := range data.Video.Media.Part.Stream {
		if val.StreamType == "2" {
			return fmt.Sprintf("%s --- %s \n", val.DisplayTitle, val.ExtendedDisplayTitle)
		}
	}

	return "fail"
}

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
