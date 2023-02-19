package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/spf13/viper"
)

var log = logger.GetLogger()

func connect(vip *viper.Viper, clientID string) (mqtt.Client, error) {
	broker := vip.GetString("mqtt.url")
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
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
	// use the topic as clientID so each invocation
	// of Publish does not trip over each other
	c, err := connect(vip, topic)
	if err != nil {
		return err
	}

	defer c.Disconnect(5000)
	// max retry
	attempts := 3

	// if there is some error, retry up to attempts
	for i := 0; i < attempts; i++ {
		log.Debugf("Sending payload %v to topic %v", string(payload), topic)
		token := c.Publish(topic, 1, false, payload)
		err = token.Error()
		log.Debugf("Error with sending MQTT: %v. Attemps: %v", err, i)
		// sleep for 1 sec and try again
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// if this doesnt return true, it timed out
		if !token.WaitTimeout(10 * time.Second) {
			log.Debug("Timeout when waiting for mqtt token")
			continue
		}
	}

	return nil
}
