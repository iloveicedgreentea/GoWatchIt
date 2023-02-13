package plex

import (
	"fmt"
	"net/http"
	"os"

	// "strings"
	"testing"
	"time"

	"github.com/StalkR/imdb"
	"github.com/stretchr/testify/assert"
)

type aspectTest struct {
	Data          testData
	ExpectedValue float64
}
type testData struct {
	Name    string
	TitleID string
	Year    int
	ID      string
}

// test to ensure server is white listed
func TestGetPlexReq(t *testing.T) {
	serverUrl := "http://192.168.88.56"
	serverPrt := "32400"
	c := NewClient(serverUrl, serverPrt)
	d, err := c.getPlexReq("/library/metadata/6262")
	assert.NoError(t, err)
	res := string(d)

	assert.NotContains(t, res, "Unauthorized", "Client is not authorized in plex server")
}

func TestGetMediaData(t *testing.T) {
	serverUrl := "http://192.168.88.56"
	serverPrt := "32400"
	c := NewClient(serverUrl, serverPrt)

	// no time to die
	med, err := c.GetMediaData("/library/metadata/7278")
	assert.NoError(t, err)

	code, err := c.GetAudioCodec(med)
	assert.NoError(t, err)
	t.Log(code)

}
func TestGetCodecFromSession(t *testing.T) {
	serverUrl := "http://192.168.88.56"
	serverPrt := "32400"
	c := NewClient(serverUrl, serverPrt)

	// no time to die
	sess, err := c.getRunningSession()
	assert.NoError(t, err)
	codec, err := c.getCodecFromSession(sess)
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
	}
	// execute each test
	for _, test := range tests {
		s := mapPlexToBeqAudioCodec(test.codec, test.fullcodec)
		assert.Equal(test.expected, s)
	}

}

func TestImdbClient(t *testing.T) {

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &customTransport{http.DefaultTransport},
	}

	title := "Lord of the rings"
	r, err := imdb.SearchTitle(client, title)
	if err != nil {
		t.Fatalf("SearchTitle(%s) error: %v", title, err)
	}
	if len(r) < 10 {
		t.Fatalf("SearchTitle(%s) len < 50: %d", title, len(r))
	}
	if accepted := map[string]bool{
		"tt7631058": true, // The Lord of the Rings (TV Series)
		"tt0120737": true, // The Lord of the Rings: The Fellowship of the Ring (2001)
	}; !accepted[r[0].ID] {
		t.Errorf("SearchTitle(%s)[0].ID = %v; want any of %v", title, r[0].ID, accepted)
	}

}

// test parsing of aspect ratio given a title
func TestImdbTechInfo(t *testing.T) {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &customTransport{http.DefaultTransport},
	}
	assert := assert.New(t)
	tests := []aspectTest{
		// test each kind of aspect + variable aspect movies until Nolan gets with the times
		{
			Data:          testData{Name: "tenet", TitleID: "tt6723592"},
			ExpectedValue: 2.39,
		},
		{
			Data:          testData{Name: "matrix", TitleID: "tt0133093"},
			ExpectedValue: 2.39,
		},
		{
			Data:          testData{Name: "21jumpst", TitleID: "tt1232829"},
			ExpectedValue: 2.35,
		},
		{
			Data:          testData{Name: "superbad", TitleID: "tt0829482"},
			ExpectedValue: 1.85,
		},
		{
			Data:          testData{Name: "theoffice", TitleID: "tt0386676"},
			ExpectedValue: 1.78,
		},
		{
			Data:          testData{Name: "ZSjusticleague", TitleID: "tt12361974"},
			ExpectedValue: 1.33,
		},
	}
	// execute each test
	for _, test := range tests {
		res, err := parseImdbTechnicalInfo(test.Data.TitleID, client)
		if err != nil {
			t.Fatalf("Test failed for %s: %v", test.Data.TitleID, err)
		}
		assert.Equal(test.ExpectedValue, res, fmt.Sprintf("%s Aspect ratio does not match", test.Data.TitleID))
	}
}

// test that it can find the correct title and return the aspect
func TestGetImdbInfoAspect(t *testing.T) {
	serverUrl := os.Getenv("PLEX_URL")
	serverPrt := os.Getenv("PLEX_PORT")
	c := NewClient(serverUrl, serverPrt)
	assert := assert.New(t)

	tests := []aspectTest{
		// test each kind of aspect + variable aspect movies until Nolan gets with the times
		{
			Data:          testData{Name: "the matrix", Year: 1999, ID: "tt0133093"},
			ExpectedValue: 2.39,
		},
		{
			Data:          testData{Name: "justice league", Year: 2021, ID: "tt12361974"},
			ExpectedValue: 1.33,
		},
		{
			Data:          testData{Name: "superbad", Year: 2007, ID: "tt0829482"},
			ExpectedValue: 1.85,
		},
	}
	for _, test := range tests {
		aspect, err := c.GetAspectRatio(test.Data.Name, test.Data.Year, test.Data.ID)
		if err != nil {
			t.Fatalf("failed for %s - %v", test.Data.Name, err)
		}
		assert.Equal(test.ExpectedValue, aspect, fmt.Sprintf("%s Aspect ratio does not match", test.Data.ID))
	}
}

// for dev only - get the entire table, ensure it can parse titles
func TestGetImdbTechInfo(t *testing.T) {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &customTransport{http.DefaultTransport},
	}
	// assert := assert.New(t)
	// superbad
	res, err := getImdbTechInfo("tt0829482", client)
	// multple aspects tenet
	// res, err := getImdbTechInfo("tt6723592", client)
	if err != nil {
		t.Fatal(err)
	}
	// Test that pulled technical info titles match
	expectedTitles := []string{"Runtime", "Sound mix", "Color", "Aspect ratio", "Camera", "Laboratory", "Film Length", "Negative Format", "Cinematographic Process", "Printed Film Format"}
	for index, title := range expectedTitles {
		assert.Equal(t, title, parseImdbTableSchema(res[index]))
	}
}

// For dev only - gets a list of every audio codec present in library
func TestGetPlexMovies(t *testing.T) {
	t.Skip()
	// edit if your movie lib is different
	librarySectionID := "1"

	serverUrl := os.Getenv("PLEX_URL")
	serverPrt := os.Getenv("PLEX_PORT")
	c := NewClient(serverUrl, serverPrt)

	data, err := c.getPlexReq(fmt.Sprintf("/library/sections/%s/all", librarySectionID))
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
