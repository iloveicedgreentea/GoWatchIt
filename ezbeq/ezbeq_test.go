package ezbeq

import (
	// "strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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

	// fast five extended 
	res, err := c.searchCatalog("51497", 2011, "DTS-X", "none", "Extended")
	assert.NoError(err)
	assert.Equal("Extended", res.Edition)
	assert.Equal("cd630eb58b05beb95ca47355c1d5014ea84e00ae8c8133573b77ee604cf7119c", res.Digest)
	
	// 12 strong has multiple codecs AND authors, so good for testing
	// return 7.1 version of aron7awol
	res, err = c.searchCatalog("429351", 2018, "DTS-HD MA 7.1", "none", "")
	assert.NoError(err)
	assert.Equal("c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d", res.Digest, "12 strong does not match entry ID. Should match aron7awol version")
	assert.Equal(-3.5, res.MvAdjust, "12 strong does not match MV. Should match aron7awol version")
	
	// return 7.1 version of mobe1969
	res, err = c.searchCatalog("429351", 2018, "DTS-HD MA 7.1", "mobe1969", "")
	assert.NoError(err)
	assert.Equal("d4ffd507ac9a6597c5039a67f587141ca866013787ed2c06fe9ef6a86f3e5534", res.Digest, "12 strong does not match entry ID. Should match mobe1969 version")
	assert.Equal(0.0, res.MvAdjust, "12 strong does not match MV. Should match mobe1969 version")
	
	// return 7.1 version of aron7awol
	res, err = c.searchCatalog("429351", 2018, "DTS-HD MA 7.1", "aron7awol", "")
	assert.NoError(err)
	assert.Equal("c694bb4c1f67903aebc51998cd1aae417983368e784ed04bf92d873ee1ca213d", res.Digest, "12 strong does not match entry ID. Should match aron7awol version")
	assert.Equal(-3.5, res.MvAdjust, "12 strong does not match MV. Should match aron7awol version")

	// return 5.1 version of aron7awol
	res, err = c.searchCatalog("429351", 2018, "DTS-HD MA 5.1", "none", "")
	assert.NoError(err)
	assert.Equal("8788e00d86868bb894fbed2f73a41e9c1d1cd277815262b7fd8ae37524c0b8a5", res.Digest, "12 strong does not match entry ID. Should match aron7awol version")
	assert.Equal(-1.5, res.MvAdjust, "12 strong does not match MV. Should match aron7awol version")

	// return DD+ Atmos of the old guard
	res, err = c.searchCatalog("547016", 2020, "DD+ Atmos", "none", "")
	assert.NoError(err)
	assert.Equal("f9bb40bed45c6e7bb2e2cdacd31e6aed3837ee23ffdfaef4c045113beec44c5d", res.Digest, "The old guard does not match entry ID. Should match aron7awol version")
	assert.Equal(0.0, res.MvAdjust, "Old guard does not match MV. Should match aron7awol version")

	// should be blank
	res, err = c.searchCatalog("some random movie", 2018, "DTS-HD MA 5.1", "none", "")
	assert.Error(err)
	assert.Equal("", res.Digest)
	assert.Equal(0.0, res.MvAdjust)
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

	// fast five dts-x extended edition
	err = c.LoadBeqProfile("51497", 2011, "DTS-X", false, "bd4577c143e73851d6db0697e0940a8f34633eec\n_416", -1.5, false, "none", "Extended", "movie")
	assert.NoError(err)

	err = c.UnloadBeqProfile(false)
	assert.NoError(err)

}