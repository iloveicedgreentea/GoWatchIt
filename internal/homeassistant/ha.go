package homeassistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

type HomeAssistantClient struct {
	ServerURL  string
	Port       string
	Token      string
	HTTPClient http.Client
	EntityName string
}

// // A client to interface with home assistant
func NewClient() (*HomeAssistantClient, error) {
	if !config.IsHomeAssistantEnabled() {
		return nil, nil
	}

	url := config.GetHomeAssistantUrl()
	port := config.GetHomeAssistantPort()
	token := config.GetHomeAssistantToken()
	entityName := config.GetHomeAssistantRemoteEntityName()
	// TODO: use scheme validation
	url = strings.ReplaceAll(url, "http://", "")
	return &HomeAssistantClient{
		ServerURL:  url,
		Port:       port,
		Token:      token,
		EntityName: entityName,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (c *HomeAssistantClient) doRequest(endpoint string, payload []byte, methodType string) ([]byte, error) {
	var req *http.Request
	var err error
	// TODO: use schema from db/client
	url := fmt.Sprintf("http://%s:%s%s", c.ServerURL, c.Port, endpoint)
	if len(payload) == 0 {
		req, err = http.NewRequest(methodType, url, http.NoBody)
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
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.GetLogger().Warn("error closing response body: %v")
		}
	}()
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
func (c *HomeAssistantClient) SwitchLight(entityType, entityName, state string) error {
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

func (c *HomeAssistantClient) SendNotification(msg string) error {
	// trigger script
	scriptData := models.HomeAssistantNotificationReq{
		Message: msg,
	}

	jsonPayload, err := json.Marshal(scriptData)
	if err != nil {
		return err
	}
	//  remove notify. if present
	name := strings.ReplaceAll(config.GetHomeAssistantNotifyEndpointName(), "notify.", "")
	endpoint := fmt.Sprintf("/api/services/notify/%s", name)
	_, err = c.doRequest(endpoint, jsonPayload, http.MethodPost)
	return err
}

// HAAttributeResponse is an interface for anything that implements these functions
type HAAttributeResponse interface {
	GetState() models.HomeAssistantMediaPlayerState
	GetSignalStatus() bool
	GetAttributes() models.Attributes
}

func (c *HomeAssistantClient) ReadAttrAndWait(ctx context.Context, waitTime int, entType models.HomeAssistantEntity, entName string, attrResp HAAttributeResponse) (bool, error) {
	var err error
	var isSignal bool
	var attributes models.Attributes
	log := logger.GetLoggerFromContext(ctx)

	// read attributes until its not nosignal
	for i := 0; i < waitTime; i++ {
		attributes, err = c.ReadAttributes(entName, attrResp, entType)
		if err != nil {
			log.Error("Error reading attributes",
				slog.String("entity", entName),
				slog.String("error", err.Error()),
			)
			return false, err
		}
		isSignal = attributes.SignalStatus
		log.Debug("Signal value",
			slog.String("entity", entName),
			slog.Bool("isSignal", isSignal),
		)
		if isSignal {
			log.Debug("HDMI sync complete")
			return isSignal, nil
		}

		// otherwise continue
		time.Sleep(200 * time.Millisecond)
	}

	return false, err
}

// TODO: read this once and then read the fields from the object instead of two calls

// ReadAttributes generic function to read attribute. entType remote || binary_sensor
func (c *HomeAssistantClient) ReadState(entityName string, respObj HAAttributeResponse, entType models.HomeAssistantEntity) (models.HomeAssistantMediaPlayerState, error) {
	endpoint := fmt.Sprintf("/api/states/%s.%s", entType, entityName)
	resp, err := c.doRequest(endpoint, nil, http.MethodGet)
	if err != nil {
		return "", err
	}

	// unmarshal
	err = json.Unmarshal(resp, respObj)
	if err != nil {
		return "", err
	}
	return respObj.GetState(), nil
}

// ReadAttributes generic function to read attribute.
func (c *HomeAssistantClient) ReadAttributes(entityName string, respObj HAAttributeResponse, entType models.HomeAssistantEntity) (models.Attributes, error) {
	endpoint := fmt.Sprintf("/api/states/%s.%s", entType, entityName)
	resp, err := c.doRequest(endpoint, nil, http.MethodGet)
	if err != nil {
		return models.Attributes{}, err
	}
	log := logger.GetLogger()
	log.Debug("Response", slog.String("response", string(resp)))
	// unmarshal
	err = json.Unmarshal(resp, respObj)
	if err != nil {
		return models.Attributes{}, err
	}
	return respObj.GetAttributes(), nil
}

func (c *HomeAssistantClient) SendEvent(eventType string, eventData map[string]interface{}) error {
	url := fmt.Sprintf("http://%s:%s/api/events/%s", c.ServerURL, c.Port, eventType)

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("error marshaling event data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.GetLogger().Warn("error closing response body: %v")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
