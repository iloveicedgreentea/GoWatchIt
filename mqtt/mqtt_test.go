package mqtt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/spf13/viper"
)

func TestPublish(t *testing.T) {
	v := viper.New()
	v.Set("mqtt.url", "tcp://192.168.88.57:1883")
	v.Set("mqtt.username", "mqtt")
	v.Set("mqtt.password", "mqtt")
	err := Publish(v, []byte("go test"), "test")
	assert.NoError(t, err)
}