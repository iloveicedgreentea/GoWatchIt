package homeassistant

import (
	// "strings"
	"strings"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/stretchr/testify/assert"
)

func testSetup() (*HomeAssistantClient) {

	haClient := NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), config.GetString("homeAssistant.remoteentityname"))

	return haClient
}

// make sure script can trigger
func TestScriptTrigger(t *testing.T) {

	haClient := testSetup()

	// trigger an empty script to verify client
	err := haClient.TriggerScript("test")
	assert.NoError(t, err)

}

// this tests an actual light
func TestLightTrigger(t *testing.T) {
	t.Skip()
	haClient := testSetup()

	// trigger light and switch

	err := haClient.SwitchLight("light", "caseta_r_wireless_in_wall_dimmer", "off")
	assert.NoError(t, err)

	err = haClient.SwitchLight("switch", "caseta_r_wireless_in_wall_neutral_switch", "off")
	assert.NoError(t, err)

}

// test sending a real notification
func TestNotification(t *testing.T) {

	haClient := testSetup()

	// trigger light and switch
	err := haClient.SendNotification("test from go-plex", config.GetString("ezbeq.notifyEndpointName"))
	assert.NoError(t, err)
}
func TestReadAttributes(t *testing.T) {
	haClient := testSetup()

	type testStruct struct {
		entName string
		test    HAAttributeResponse
		entType string
	}
	tt := []testStruct{
		{
			entName: "nz7",
			test:    &models.HAjvcResponse{},
			entType: "remote",
		},
		{
			entName: "envy",
			test:    &models.HAEnvyResponse{},
			entType: "remote",
		},
		{
			entName: "test_sensor",
			test:    &models.HABinaryResponse{},
			entType: "binary_sensor",
		},
	}

	for _, k := range tt {
		k := k

		signal, err := haClient.ReadAttributes(k.entName, k.test, k.entType)
		t.Log(k.entName, signal)

		// if its off, we expect an error
		if err != nil {
			if strings.Contains(err.Error(), "state is off") {
				t.Logf("%s is off", k.entName)
				assert.Equal(t, false, signal)
			} else {
				t.Error(err)
			}
		}

	}
}
