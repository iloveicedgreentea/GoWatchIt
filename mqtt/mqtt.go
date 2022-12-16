package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

func connect(vip *viper.Viper) (mqtt.Client, error) {
	broker := vip.GetString("mqtt.url")
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID("plexAutomatorClient")
	opts.SetUsername(vip.GetString("mqtt.username"))
	opts.SetPassword(vip.GetString("mqtt.password"))

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

// creates a connection to broker and sends the payload
func Publish(vip *viper.Viper, payload []byte, topic string) error {
	c, err := connect(vip)
	
	if err != nil {
		return err
	}
	defer c.Disconnect(10)
	
	token := c.Publish(topic, 1, false, payload)
	err = token.Error()
	if err != nil {
		return err
	}

	return nil
}
