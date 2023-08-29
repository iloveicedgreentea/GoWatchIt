package mqtt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/spf13/viper"
	"fmt"
)

func TestPublish(t *testing.T) {
	v := viper.New()
	v.Set("mqtt.url", "tcp://192.168.88.57:1883")
	v.Set("mqtt.username", "mqtt")
	v.Set("mqtt.password", "mqtt")
	err := Publish([]byte(fmt.Sprintf("{\"aspect\":%f}", 2.4)), "theater/jvc/aspectratio")
	assert.NoError(t, err)
	err = Publish([]byte(fmt.Sprintf("{\"type\":\"%s\"}", "episode")), "theater/denon/volume")
	assert.NoError(t, err)
}