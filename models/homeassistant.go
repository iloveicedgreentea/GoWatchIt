package models

import (
	"strconv"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
)

var log = logger.GetLogger()

type HomeAssistantScriptReq struct {
	EntityID string `json:"entity_id"`
}

type HomeAssistantNotificationReq struct {
	Message string `json:"message"`
}

type HomeAssistantWebhookPayload struct {
}

type HAEnvyResponse struct {
	EntityID   string         `json:"entity_id"`
	State      string         `json:"state"`
	Attributes EnvyAttributes `json:"attributes"`
}
type EnvyAttributes struct {
	SignalStatus bool `json:"is_signal"`
}
type HAjvcResponse struct {
	EntityID   string        `json:"entity_id"`
	State      string        `json:"state"`
	Attributes JVCAttributes `json:"attributes"`
}
type JVCAttributes struct {
	SignalStatus string `json:"signal_status"`
}

type HABinaryResponse struct {
	State string `json:"state"`
}

func (r *HABinaryResponse) GetState() string {
	return r.State
}
func (r *HABinaryResponse) GetSignalStatus() bool {
	return false
}

func (r *HAEnvyResponse) GetState() string {
	return r.State
}
func (r *HAEnvyResponse) GetSignalStatus() bool {
	return r.Attributes.SignalStatus
}
func (r *HAjvcResponse) GetState() string {
	return r.State
}
func (r *HAjvcResponse) GetSignalStatus() bool {
	s, err := strconv.ParseBool(r.Attributes.SignalStatus)
	if err != nil {
		log.Error("error parsing JVC signal attribute",
			"error", err,
		)
		return false
	}
	return s
}
