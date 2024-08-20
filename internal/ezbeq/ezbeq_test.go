package ezbeq

import (
	"fmt"
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

// TestMuteCmds send commands to minidsp
func TestMuteCmds(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient()
	assert.NoError(err)

	// send mute commands
	assert.NoError(c.MuteCommand(true))
	assert.NoError(c.MuteCommand(false))
}
func TestGetStatus(t *testing.T) {
	c := &BeqClient{
		ServerURL: config.GetEZBeqUrl(),
		Port:      config.GetEZBeqPort(),
	}

	// send mute commands
	assert.NotEmpty(t, c.Port)
	assert.NotEmpty(t, c.ServerURL)

	err := c.GetStatus()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.DeviceInfo)

}

func TestSingleDevice(t *testing.T) {
	rawJson := `{
		"master": {
		  "type": "minidsp",
		  "name": "master",
		  "masterVolume": 0,
		  "mute": false,
		  "slots": [
			{
			  "id": "1",
			  "last": "10 Minutes Gone",
			  "active": true,
			  "gains": [
				{
				  "id": "1",
				  "value": 5
				},
				{
				  "id": "2",
				  "value": 5
				}
			  ],
			  "mutes": [
				{
				  "id": "1",
				  "value": false
				},
				{
				  "id": "2",
				  "value": false
				}
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "2",
			  "last": "Empty",
			  "active": false,
			  "gains": [
				{
				  "id": "1",
				  "value": 1.5
				},
				{
				  "id": "2",
				  "value": 1.5
				}
			  ],
			  "mutes": [
				{
				  "id": "1",
				  "value": false
				},
				{
				  "id": "2",
				  "value": false
				}
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "3",
			  "last": "10 Minutes Gone",
			  "active": false,
			  "gains": [
				{
				  "id": "1",
				  "value": 0
				},
				{
				  "id": "2",
				  "value": 0
				}
			  ],
			  "mutes": [
				{
				  "id": "1",
				  "value": false
				},
				{
				  "id": "2",
				  "value": false
				}
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "4",
			  "last": "Ready Player One",
			  "active": false,
			  "gains": [
				{
				  "id": "1",
				  "value": 2
				},
				{
				  "id": "2",
				  "value": 2
				}
			  ],
			  "mutes": [
				{
				  "id": "1",
				  "value": false
				},
				{
				  "id": "2",
				  "value": false
				}
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			}
		  ]
		}
	  }`

	payload, err := mapToBeqDevice([]byte(rawJson))
	assert.NoError(t, err)
	var deviceInfo []models.BeqDevices
	for _, v := range payload {
		deviceInfo = append(deviceInfo, v)
	}

	assert.NotEmpty(t, deviceInfo)
	assert.Len(t, deviceInfo, 1)

	for _, v := range deviceInfo {
		t.Log(v.Name)
	}
}
func TestDualDevice(t *testing.T) {
	rawJson := `{
		"master": {
		  "name": "master",
		  "masterVolume": 0,
		  "mute": true,
		  "slots": [
			{
			  "id": "1",
			  "last": "Empty",
			  "active": true,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "2",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "3",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "4",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			}
		  ]
		},
		"master2": {
		  "name": "master2",
		  "masterVolume": 0,
		  "mute": true,
		  "slots": [
			{
			  "id": "1",
			  "last": "Empty",
			  "active": true,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "2",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "3",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			},
			{
			  "id": "4",
			  "last": "Empty",
			  "active": false,
			  "gain1": 0,
			  "gain2": 0,
			  "mute1": false,
			  "mute2": false,
			  "gains": [
				0,
				0
			  ],
			  "mutes": [
				false,
				false
			  ],
			  "canActivate": true,
			  "inputs": 2,
			  "outputs": 4
			}
		  ]
		}
	  }`

	payload, err := mapToBeqDevice([]byte(rawJson))
	assert.NoError(t, err)
	var deviceInfo []models.BeqDevices
	for _, v := range payload {
		deviceInfo = append(deviceInfo, v)
	}

	assert.NotEmpty(t, deviceInfo)
	assert.Len(t, deviceInfo, 2)
	for _, v := range deviceInfo {
		t.Log(v.Name)
	}
}
func TestUrlEncode(t *testing.T) {
	s := urlEncode("DTS-HD MA 7.1")
	assert.Equal(t, "DTS-HD+MA+7.1", s)
}

func TestHasAuthor(t *testing.T) {
	assert := assert.New(t)
	type testStruct struct {
		author   string
		expected bool
	}
	tt := []testStruct{
		{
			author:   "aron7awol",
			expected: true,
		},
		{
			author:   "None",
			expected: false,
		},
		{
			author:   " ",
			expected: false,
		},
		{
			author:   "",
			expected: false,
		},
		{
			author:   "none",
			expected: false,
		},
		{
			author:   "aron7awol, mobe1969",
			expected: true,
		},
	}
	for _, tc := range tt {
		tc := tc
		s := hasAuthor(tc.author)
		assert.Equal(tc.expected, s)
	}
}

func TestBuildAuthorWhitelist(t *testing.T) {

	s := buildAuthorWhitelist("aron7awol, mobe1969", "/api/1/search?audiotypes=dts-x&years=2011&tmdbid=12345")
	assert.Equal(t, "/api/1/search?audiotypes=dts-x&years=2011&tmdbid=12345&authors=aron7awol&authors=mobe1969", s)
}

func TestSearchCatalog(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient()
	assert.NoError(err)

	// list of testing structs
	type testStruct struct {
		m                               models.BeqSearchRequest
		expectedEdition, expectedDigest string
		expectedMvAdjust                float64
	}
	tt := []testStruct{
		{
			// fast five extended
			m: models.BeqSearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "none",
				Edition:         "Extended",
			},
			expectedEdition:  "Extended",
			expectedDigest:   "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
			expectedMvAdjust: -1.5,
		},
		{
			// Jung E
			m: models.BeqSearchRequest{
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
			m: models.BeqSearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "",
				Edition:         "Extended",
			},
			expectedEdition:  "Extended",
			expectedDigest:   "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
			expectedMvAdjust: -1.5,
		},
		{
			m: models.BeqSearchRequest{
				TMDB:            "51497",
				Year:            2011,
				Codec:           "DTS-X",
				PreferredAuthor: "None",
				Edition:         "Extended",
			},
			expectedEdition:  "Extended",
			expectedDigest:   "cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c",
			expectedMvAdjust: -1.5,
		},
		{
			// 12 strong has multiple codecs AND authors, so good for testing
			// return 7.1 version of aron7awol
			m: models.BeqSearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d",
			expectedMvAdjust: -3.5,
		},
		{
			// return 7.1 version of mobe1969
			m: models.BeqSearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "mobe1969",
				Edition:         "",
			},
			expectedEdition: "",
			expectedDigest:  "73a1eef9ce33abba7df0a9d2b4cec41254f6a521d521e104fa3cd2e7297c26d9",
		},
		{
			// return 7.1 version with mutliple authors
			m: models.BeqSearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "mobe1969, aron7awol",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d",
			expectedMvAdjust: -3.5,
		},
		{
			// return 7.1 version with mutliple authors
			m: models.BeqSearchRequest{
				TMDB:            "429351",
				Year:            2018,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "aron7awol,mobe1969",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d",
			expectedMvAdjust: -3.5,
		},
		{
			// 12 strong has multiple codecs AND authors, so good for testing
			// return 7.1 version of aron7awol
			m: models.BeqSearchRequest{
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
			m: models.BeqSearchRequest{
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
			m: models.BeqSearchRequest{
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
			// should be TrueHD 7.1
			m: models.BeqSearchRequest{
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
		{
			//  spiderman universe
			m: models.BeqSearchRequest{
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
		res, err := c.searchCatalog(&tc.m)
		assert.NoError(err)
		assert.Equal(tc.expectedDigest, res.Digest, fmt.Sprintf("digest did not match %s", res.Digest))
		assert.Equal(tc.expectedEdition, res.Edition, fmt.Sprintf("edition did not match %s", res.Digest))
		assert.Equal(tc.expectedMvAdjust, res.MvAdjust, fmt.Sprintf("MV did not match %s", res.Digest))
	}

	_, err = c.searchCatalog(&models.BeqSearchRequest{
		TMDB:            "ojdsfojnekfw",
		Year:            2018,
		Codec:           "DTS-HD MA 5.1",
		PreferredAuthor: "none",
		Edition:         "",
	})
	assert.Error(err)
}

// load and unload a profile. Watch ezbeq UI to confirm, but if it doesnt error it probably loaded fine
// ezbeq doesnt expose a failure if the entry_id is wrong, so need to look at UI for now
// I could write a scraper to find instance of fast five in slot one, thats a lot of work for a small test
func TestLoadProfile(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient()
	assert.NoError(err)

	tt := []models.BeqSearchRequest{
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
		// {
		// 	TMDB:            "56292",
		// 	Year:            2011,
		// 	Codec:           "AtmosMaybe",
		// 	SkipSearch:      false,
		// 	EntryID:         "",
		// 	MVAdjust:        0.0,
		// 	DryrunMode:      false,
		// 	PreferredAuthor: "none",
		// 	Edition:         "",
		// 	MediaType:       "movie",
		// 	Devices:         []string{"master", "master2"},
		// 	Slots:           []int{1},
		// },
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
		// DD+Atmos5.1Maybe //underwater
		{
			TMDB:            "443791",
			Year:            2020,
			Codec:           "DD+Atmos5.1Maybe",
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
			TMDB:            "804095",
			Year:            2022,
			Codec:           "DD+Atmos7.1Maybe",
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
		err = c.LoadBeqProfile(&tc)
		assert.NoError(err)

		err = c.UnloadBeqProfile(&tc)
		assert.NoError(err)
	}

}
