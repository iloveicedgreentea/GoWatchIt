package ezbeq

import (
	"database/sql"
	"fmt"
	l "log"
	"os"
	"sync"
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

// TestMain sets up the database and runs the tests
func TestMain(m *testing.M) {
	var code int
	dbOnce.Do(func() {
		// Setup code before tests
		var err error

		// Open SQLite database connection
		db, err = database.GetDB(":memory:")
		if err != nil {
			l.Fatalf("Failed to open database: %v", err)
		}

		// run migrations
		err = database.RunMigrations(db)
		if err != nil {
			l.Fatalf("Failed to run migrations: %v", err)
		}

		// Initialize the config with the database
		err = config.InitConfig(db)
		if err != nil {
			l.Fatalf("Failed to initialize config: %v", err)
		}

		cf := config.GetConfig()

		// populate test data
		beqCfg := models.EZBEQConfig{
			Enabled:              true,
			DryRun:               true,
			URL:                  "ezbeq.local",
			Port:                 "8080",
			Scheme:               "http",
			LooseEditionMatching: true,
			SkipEditionMatching:  false,
		}
		err = cf.SaveConfig(&beqCfg)
		if err != nil {
			l.Fatalf("Failed to save ezbeq config: %v", err)
		}
		// Run the tests
		code = m.Run()

		// Cleanup code after tests
		err = db.Close()
		if err != nil {
			l.Printf("Error closing database: %v", err)
		}
	})
	// Exit with the test result code
	os.Exit(code)
}

// TestMuteCmds send commands to minidsp
func TestMuteCmds(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := NewClient()
	a.NoError(err)

	// send mute commands
	a.NoError(c.MuteCommand(true))
	a.NoError(c.MuteCommand(false))
}

func TestCheckEdition(t *testing.T) {
	t.Parallel()
	type test struct {
		beqEdition string
		edition    models.Edition
		expected   bool
	}
	tests := []test{
		{
			beqEdition: "Extended",
			edition:    models.EditionExtended,
			expected:   true,
		},
		{
			beqEdition: "ex",
			edition:    models.EditionExtended,
			expected:   true,
		},
		{
			beqEdition: "EX",
			edition:    models.EditionExtended,
			expected:   true,
		},
		{
			beqEdition: "DC",
			edition:    models.EditionDirectorsCut,
			expected:   true,
		},
		{
			beqEdition: "DC+SE+TC",
			edition:    models.EditionDirectorsCut,
			expected:   true,
		},
	}

	for i := range tests {
		test := tests[i] // Capture range variable
		t.Run(fmt.Sprintf("Edition_%s", test.beqEdition), func(t *testing.T) {
			t.Parallel()
			match := checkEdition(&models.BeqCatalog{Edition: test.beqEdition}, test.edition)
			assert.Equal(t, test.expected, match, "Expected: ", test.expected, "Got: ", match, "for ", test.beqEdition)
		})
	}
}

func TestGetStatus(t *testing.T) {
	t.Parallel()
	c, err := NewClient()
	assert.NoError(t, err)

	// send mute commands
	assert.NotEmpty(t, c.Port)
	assert.NotEmpty(t, c.ServerURL)

	err = c.GetStatus()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.DeviceInfo)
}

func TestSingleDevice(t *testing.T) {
	t.Parallel()
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

	deviceInfo := make([]models.BeqDevices, 0, len(payload))
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
	t.Parallel()
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
	// preallocate memory
	deviceInfo := make([]models.BeqDevices, 0, len(payload))
	for _, v := range payload {
		deviceInfo = append(deviceInfo, v)
	}

	assert.NotEmpty(t, deviceInfo)
	assert.Len(t, deviceInfo, 2)
	for _, v := range deviceInfo {
		t.Log(v.Name)
	}
}

func TestHasAuthor(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
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
	for i := range tt {
		tc := tt[i] // Capture range variable
		t.Run(fmt.Sprintf("Author_%s", tc.author), func(t *testing.T) {
			t.Parallel()
			s := hasAuthor(tc.author)
			a.Equal(tc.expected, s)
		})
	}
}

func TestSearchCatalog(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := NewClient()
	a.NotNil(c)
	a.NoError(err)

	// list of testing structs
	type testStruct struct {
		m                               models.BeqSearchRequest
		expectedEdition, expectedDigest string
		expectedMvAdjust                float64
	}
	tt := []testStruct{
		{
			// Stargate (1994) {edition-Extended Edition} Remux 1080p
			m: models.BeqSearchRequest{
				TMDB:            "2164",
				Year:            1994,
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "none",
				Edition:         "Extended",
			},
			expectedEdition:  "Extended Cut",
			expectedDigest:   "6d9cfaed8335a348491eebae27f7f5fb11752e32df64b46d24d6f995dd74d96d",
			expectedMvAdjust: 0,
		},
		{
			// Stargate (1994) {edition-Extended Edition} Remux 1080p no year
			m: models.BeqSearchRequest{
				TMDB:            "2164",
				Codec:           "DTS-HD MA 7.1",
				PreferredAuthor: "none",
				Edition:         "Extended",
			},
			expectedEdition:  "Extended Cut",
			expectedDigest:   "6d9cfaed8335a348491eebae27f7f5fb11752e32df64b46d24d6f995dd74d96d",
			expectedMvAdjust: 0,
		},
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
			// return 7.1 version with multiple authors
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
			// return 7.1 version with multiple authors
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
		{
			//  spiderman universe blank year
			m: models.BeqSearchRequest{
				TMDB:            "56292",
				Year:            0,
				Codec:           "TrueHD 7.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "f7e8c32e58b372f1ea410165607bc1f6b3f589a832fda87edaa32a17715438f7",
			expectedMvAdjust: 0.0,
		},
		{
			//  spiderman universe no year
			m: models.BeqSearchRequest{
				TMDB:            "56292",
				Codec:           "TrueHD 7.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "",
			expectedDigest:   "f7e8c32e58b372f1ea410165607bc1f6b3f589a832fda87edaa32a17715438f7",
			expectedMvAdjust: 0.0,
		},
		{
			//  Star Wars (1977) {edition-Project 4K77} Remux 2160p DTS-HD MA
			m: models.BeqSearchRequest{
				TMDB:            "11",
				Year:            1977,
				Codec:           "DTS-HD MA 5.1",
				PreferredAuthor: "none",
				Edition:         "",
			},
			expectedEdition:  "Project 4K77",
			expectedDigest:   "83954ea27172605f8bdd8c4731bfc5f164075ce05d436cd319ea13db9978110a",
			expectedMvAdjust: 1.0,
		},
		{
			//  Star Wars (1977) {edition-Project 4K77} Remux 2160p DTS-HD MA
			m: models.BeqSearchRequest{
				TMDB:            "11",
				Year:            1977,
				Codec:           "DTS-HD MA 5.1",
				PreferredAuthor: "none",
				Edition:         "Project 4K77",
			},
			expectedEdition:  "Project 4K77",
			expectedDigest:   "83954ea27172605f8bdd8c4731bfc5f164075ce05d436cd319ea13db9978110a",
			expectedMvAdjust: 1.0,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.m.Title, func(t *testing.T) {
			t.Parallel()
			res, err := c.searchCatalog(&tc.m)
			a.NoError(err)
			a.Equal(tc.expectedDigest, res.Digest, fmt.Sprintf("digest did not match %s", res.Digest))
			a.Equal(tc.expectedEdition, res.Edition, fmt.Sprintf("edition did not match %s", res.Digest))
			a.Equal(tc.expectedMvAdjust, res.MvAdjust, fmt.Sprintf("MV did not match %s", res.Digest))
		})
	}
	// should always fail
	_, err = c.searchCatalog(&models.BeqSearchRequest{
		TMDB:            "ojdsfojnekfw",
		Year:            2018,
		Codec:           "DTS-HD MA 5.1",
		PreferredAuthor: "none",
		Edition:         "",
	})
	a.Error(err)
}

// load and unload a profile. Watch ezbeq UI to confirm, but if it doesnt error it probably loaded fine
// ezbeq doesnt expose a failure if the entry_id is wrong, so need to look at UI for now
// I could write a scraper to find instance of fast five in slot one, thats a lot of work for a small test
func TestLoadProfile(t *testing.T) {
	t.Parallel()
	if os.Getenv("RUN_INTEGRATION") != "true" {
		t.Skip("skipping TestLoadProfile test")
	}
	a := assert.New(t)

	c, err := NewClient()
	a.NoError(err)

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

	// this should not be parallel
	for _, tc := range tt {
		err = c.LoadBeqProfile(&tc)
		a.NoError(err)

		err = c.UnloadBeqProfile(&tc)
		a.NoError(err)
	}
}
