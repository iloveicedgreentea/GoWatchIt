package homeassistant

import (
	// "strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// make sure script can trigger
func TestScriptTrigger(t *testing.T) {

	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}

	// trigger an empty script to verify client
	haClient := NewClient(v.GetString("homeAssistant.url"), v.GetString("homeAssistant.port"), v.GetString("homeAssistant.token"))
	err = haClient.TriggerScript("test")
	assert.NoError(t, err)

}

// this tests an actual light
func TestLightTrigger(t *testing.T) {
	t.Skip()
	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}

	// trigger light and switch
	haClient := NewClient(v.GetString("homeAssistant.url"), v.GetString("homeAssistant.port"), v.GetString("homeAssistant.token"))

	err = haClient.SwitchLight("light", "caseta_r_wireless_in_wall_dimmer", "off")
	assert.NoError(t, err)

	err = haClient.SwitchLight("switch", "caseta_r_wireless_in_wall_neutral_switch", "off")
	assert.NoError(t, err)

}

// test sending a real notification
func TestNotification(t *testing.T) {
	
	v := viper.New()
	v.SetConfigFile("../config.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}

	// trigger light and switch
	haClient := NewClient(v.GetString("homeAssistant.url"), v.GetString("homeAssistant.port"), v.GetString("homeAssistant.token"))

	err = haClient.SendNotification("test from go-plex", v.GetString("ezbeq.notifyEndpointName"))
	assert.NoError(t, err)
}