package homeassistant

import (
	// "strings"
	"strings"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func testSetup() (*viper.Viper, *HomeAssistantClient, error) {
	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	haClient := NewClient(v.GetString("homeAssistant.url"), v.GetString("homeAssistant.port"), v.GetString("homeAssistant.token"), v.GetString("homeAssistant.envyRemoteName"), v.GetString("homeAssistant.jvcRemoteName"), v.GetString("homeAssistant.binarySensorName"))

	return v, haClient, err
}

// make sure script can trigger
func TestScriptTrigger(t *testing.T) {

	_, haClient, err := testSetup()
	assert.NoError(t, err)

	// trigger an empty script to verify client
	err = haClient.TriggerScript("test")
	assert.NoError(t, err)

}

// this tests an actual light
func TestLightTrigger(t *testing.T) {
	t.Skip()
	_, haClient, err := testSetup()
	assert.NoError(t, err)

	// trigger light and switch

	err = haClient.SwitchLight("light", "caseta_r_wireless_in_wall_dimmer", "off")
	assert.NoError(t, err)

	err = haClient.SwitchLight("switch", "caseta_r_wireless_in_wall_neutral_switch", "off")
	assert.NoError(t, err)

}

// test sending a real notification
func TestNotification(t *testing.T) {

	v, haClient, err := testSetup()
	assert.NoError(t, err)

	// trigger light and switch
	err = haClient.SendNotification("test from go-plex", v.GetString("ezbeq.notifyEndpointName"))
	assert.NoError(t, err)
}
func TestReadAttributes(t *testing.T) {
	_, haClient, err := testSetup()
	assert.NoError(t, err)

	type testStruct struct {
		entName string
		test    HAAttributeResponse
		entType string
	}
	tt := []testStruct{
		{
			entName: haClient.JVCEntityName,
			test:    &models.HAjvcResponse{},
			entType: "remote",
		},
		{
			entName: haClient.EnvyEntityName,
			test:    &models.HAEnvyResponse{},
			entType: "remote",
		},
		{
			entName: haClient.BinaryName,
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
