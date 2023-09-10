package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
)

var log = logger.GetLogger()

func connect(clientID string) (mqtt.Client, error) {
	broker := config.GetString("mqtt.url")
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(config.GetString("mqtt.username"))
	opts.SetPassword(config.GetString("mqtt.password"))

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func PublishWrapper(topic string, msg string) error {
	// trigger automation
	return Publish([]byte(msg), config.GetString(fmt.Sprintf("mqtt.%s", topic)))
}

// creates a connection to broker and sends the payload
func Publish(payload []byte, topic string) error {
	// use the topic as clientID so each invocation
	// of Publish does not trip over each other
	c, err := connect(topic)
	if err != nil {
		return fmt.Errorf("error connecting to topic - %v", err)
	}

	defer c.Disconnect(5000)
	// max retry
	attempts := 3

	// if there is some error, retry up to attempts
	for i := 0; i < attempts; i++ {
		log.Debugf("Sending payload %v to topic %v", string(payload), topic)
		token := c.Publish(topic, 1, false, payload)
		err = token.Error()
		// sleep for 1 sec and try again
		if err != nil {
			log.Debugf("Error with sending MQTT: %v. Attemps: %v", err, i)
			time.Sleep(1 * time.Second)
			continue
		}

		// if this doesnt return true, it timed out
		if !token.WaitTimeout(10 * time.Second) {
			log.Debug("Timeout when waiting for mqtt token")
			continue
		}

		break
	}

	return nil
}
