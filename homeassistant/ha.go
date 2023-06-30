package homeassistant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iloveicedgreentea/go-plex/models"
)

type HomeAssistantClient struct {
	ServerURL  string
	Port       string
	Token      string
	HTTPClient http.Client
}

// // A client to interface with home assistant
func NewClient(url, port string, token string) *HomeAssistantClient {
	return &HomeAssistantClient{
		ServerURL: url,
		Port:      port,
		Token:     token,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HomeAssistantClient) doRequest(endpoint string, payload []byte) error {
	// bodyReader := bytes.NewReader(jsonBody)
	url := fmt.Sprintf("%s:%s%s", c.ServerURL, c.Port, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}

// run arbitrary scripts
func (c *HomeAssistantClient) TriggerScript(scriptName string) error {
	// trigger script
	scriptData :=  models.HomeAssistantScriptReq{
		EntityID: fmt.Sprintf("script.%s", scriptName),
	}

	jsonPayload, err := json.Marshal(scriptData)
	if err != nil {
		return err
	}
	return c.doRequest("/api/services/script/turn_on", jsonPayload)
}

// switch a light on/off
func (c *HomeAssistantClient) SwitchLight(entityType string, entityName string, state string) error {
	// trigger script
	payload :=  models.HomeAssistantScriptReq{
		EntityID: fmt.Sprintf("%s.%s", entityType, entityName),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var action string
	switch state{
	case "on":
		action = "turn_on"
	case "off":
		action = "turn_off"
	}
	endpoint := fmt.Sprintf("/api/services/%s/%s", entityType, action)
	return c.doRequest(endpoint, jsonPayload)
}

func (c *HomeAssistantClient) SendNotification(msg string, endpointName string) error {	
	// trigger script
	scriptData :=  models.HomeAssistantNotificationReq{
		Message: msg,
	}

	jsonPayload, err := json.Marshal(scriptData)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/api/services/notify/%s", endpointName)
	return c.doRequest(endpoint, jsonPayload)
}

func (c *HomeAssistantClient) ReadEnvyAttributes() (bool, error) {
	return true, nil
}
