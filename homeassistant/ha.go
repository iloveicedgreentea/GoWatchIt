package homeassistant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

type HomeAssistantClient struct {
	ServerURL  string
	Port       string
	Token      string
	HTTPClient http.Client
	EnvyEntityName string
}

// // A client to interface with home assistant
func NewClient(url, port string, token string, envyName string) *HomeAssistantClient {
	return &HomeAssistantClient{
		ServerURL: url,
		Port:      port,
		Token:     token,
		EnvyEntityName: envyName,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HomeAssistantClient) doRequest(endpoint string, payload []byte, methodType string) ([]byte, error) {
	var req *http.Request
	var err error

	log.Debugf("Using method %s", methodType)
	// bodyReader := bytes.NewReader(jsonBody)
	url := fmt.Sprintf("%s:%s%s", c.ServerURL, c.Port, endpoint)
	if len(payload) == 0 {
		req, err = http.NewRequest(methodType, url, nil)
	} else {
		req, err = http.NewRequest(methodType, url, bytes.NewBuffer(payload))
	}
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []byte{}, errors.New(res.Status)
	}

	return resp, nil
}

// run arbitrary scripts
func (c *HomeAssistantClient) TriggerScript(scriptName string) error {
	// trigger script
	scriptData := models.HomeAssistantScriptReq{
		EntityID: fmt.Sprintf("script.%s", scriptName),
	}

	jsonPayload, err := json.Marshal(scriptData)
	if err != nil {
		return err
	}
	_, err = c.doRequest("/api/services/script/turn_on", jsonPayload, http.MethodPost)
	return err
}

// switch a light on/off
func (c *HomeAssistantClient) SwitchLight(entityType string, entityName string, state string) error {
	// trigger script
	payload := models.HomeAssistantScriptReq{
		EntityID: fmt.Sprintf("%s.%s", entityType, entityName),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var action string
	switch state {
	case "on":
		action = "turn_on"
	case "off":
		action = "turn_off"
	}
	endpoint := fmt.Sprintf("/api/services/%s/%s", entityType, action)
	_, err = c.doRequest(endpoint, jsonPayload, http.MethodPost)
	return err
}

func (c *HomeAssistantClient) SendNotification(msg string, endpointName string) error {
	// trigger script
	scriptData := models.HomeAssistantNotificationReq{
		Message: msg,
	}

	jsonPayload, err := json.Marshal(scriptData)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/api/services/notify/%s", endpointName)
	_, err = c.doRequest(endpoint, jsonPayload, http.MethodPost)
	return err
}



// ReadEnvyAttributes returns true if there is a signal
func (c *HomeAssistantClient) ReadEnvyAttributes() (bool, error) {
	endpoint := fmt.Sprintf("/api/states/remote.%s", c.EnvyEntityName)
	resp, err := c.doRequest(endpoint, nil, http.MethodGet)
	if err != nil {
		return false, err
	}
	log.Debugf("Response: %s", resp)

	// unmarshal
	var envyResp models.HAEnvyResponse
	err = json.Unmarshal(resp, &envyResp)

	return envyResp.Attributes.NoSignal, err
}
