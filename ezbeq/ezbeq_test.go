package ezbeq

import (
	// "strings"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestMuteCmds send commands to minidsp
func TestMuteCmds(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("../config.json")

	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	c, err := NewClient(v.GetString("ezbeq.url"), v.GetString("ezbeq.port"))
	assert.NoError(err)

	// send mute commands
	assert.NoError(c.MuteCommand(true))
	assert.NoError(c.MuteCommand(false))
}
func TestUrlEncode(t *testing.T) {
	s := urlEncode("DTS-HD MA 7.1")
	assert.Equal(t, "DTS-HD+MA+7.1", s)
}

func TestSearchCatalog(t *testing.T) {

	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	c, err := NewClient(v.GetString("ezbeq.url"), v.GetString("ezbeq.port"))
	assert.NoError(err)

	// list of testing structs
	type testStruct struct {
		m                               models.SearchRequest
		expectedEdition, expectedDigest string
		expectedMvAdjust                float64
	}
	tt := []testStruct{
		{
			// fast five extended
			m: models.SearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "none",
				Edition:         "Extended",
			},
			expectedEdition: "Extended",
			expectedDigest:  "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
		},
		{
			// Jung E
			m: models.SearchRequest{
				TMDB:            "843794",
				Year:            2023,
				Codec:           "DD+ Atmos",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition: "",
			expectedDigest:  "1678d7860ead948132f70ba3d823d7493bb3bb79302f308d135176bf4ff6f7d0",
		},
		{
			m: models.SearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "",
				Edition:         "Extended",
			},
			expectedEdition: "Extended",
			expectedDigest:  "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
		},
		{
			m: models.SearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "None",
				Edition:         "Extended",
			},
			expectedEdition: "Extended",
			expectedDigest:  "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
		},
		{
			// 12 strong has multiple codecs AND authors, so good for testing
			// return 7.1 version of aron7awol
			m: models.SearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "8788e00d86868bb894fbed2f73a41e9c1d1cd277815262b7fd8ae37524c0b8a5",
			expectedMvAdjust: -1.5,
		},
		{
			// return 7.1 version of mobe1969
			m: models.SearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "mobe1969",
				Edition:         "",
			},
			expectedEdition: "",
			expectedDigest:  "d4ffd507ac9a6597c5039a67f587141ca866013787ed2c06fe9ef6a86f3e5534",
		},
		{
			// 12 strong has multiple codecs AND authors, so good for testing
			// return 7.1 version of aron7awol
			m: models.SearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "aron7awol",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d",
			expectedMvAdjust: -3.5,
		},
		{
			// return 5.1 version of aron7awol
			m: models.SearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 5.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "8788e00d86868bb894fbed2f73a41e9c1d1cd277815262b7fd8ae37524c0b8a5",
			expectedMvAdjust: -1.5,
		},
		{
			// return 5.1 version of aron7awol
			m: models.SearchRequest{
				TMDB:            "547016",
				Year:            2020,
				Codec:           "DD+ Atmos",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "f9bb40bed45c6e7bb2e2cdacd31e6aed3837ee23ffdfaef4c045113beec44c5d",
			expectedMvAdjust: 0.0,
		},
		{
			// should be blank
			m: models.SearchRequest{
				TMDB:            "ojdsfojnekfw",
				Year:            2018,
				Codec:           "DTS-HD MA 5.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "",
			expectedMvAdjust: 0.0,
		},
		{
			// should be TrueHD 7.1
			m: models.SearchRequest{
				TMDB:            "56292",
				Year:            2011,
				Codec:           "TrueHD 7.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "f7e8c32e58b372f1ea410165607bc1f6b3f589a832fda87edaa32a17715438f7",
			expectedMvAdjust: 0.0,
		},
	}

	for _, tc := range tt {
		tc := tc
		// should be TrueHD 7.1
		res, err := c.searchCatalog(tc.m)
		assert.NoError(err)
		assert.Equal(tc.expectedDigest, res.Digest)
		assert.Equal(tc.expectedEdition, res.Edition)
		assert.Equal(tc.expectedMvAdjust, res.MvAdjust)
	}
}

// load and unload a profile. Watch ezbeq UI to confirm, but if it doesnt error it probably loaded fine
// ezbeq doesnt expose a failure if the entry_id is wrong, so need to look at UI for now
// I could write a scraper to find instance of fast five in slot one, thats a lot of work for a small test
func TestLoadProfile(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	c, err := NewClient(v.GetString("ezbeq.url"), v.GetString("ezbeq.port"))
	assert.NoError(err)

	tt := []models.SearchRequest{
		{
			TMDB:            "51497",
			Year:            2011,
			Codec:           "DTS-X",
			SkipSearch:      false,
			EntryID:         "bd4577c143e73851d6db0697e0940a8f34633eec\n_416",
			MVAdjust:        -1.5,
			DryrunMode:      false,
			PreferredAuthor: "none",
			Edition:         "Extended",
			MediaType:       "movie",
			Devices:         []string{"master", "master2"},
			Slots:           []int{1},
		},
		{
			TMDB:            "56292",
			Year:            2011,
			Codec:           "AtmosMaybe",
			SkipSearch:      false,
			EntryID:         "",
			MVAdjust:        0.0,
			DryrunMode:      false,
			PreferredAuthor: "none",
			Edition:         "",
			MediaType:       "movie",
			Devices:         []string{"master", "master2"},
			Slots:           []int{1},
		},
		{
			TMDB:            "399579",
			Year:            2019,
			Codec:           "AtmosMaybe",
			SkipSearch:      false,
			EntryID:         "",
			MVAdjust:        0.0,
			DryrunMode:      false,
			PreferredAuthor: "none",
			Edition:         "",
			MediaType:       "movie",
			Devices:         []string{"master", "master2"},
			Slots:           []int{1},
		},
	}

	for _, tc := range tt {
		tc := tc
		err = c.LoadBeqProfile(tc)
		assert.NoError(err)

		err = c.UnloadBeqProfile(tc)
		assert.NoError(err)
	}

}
