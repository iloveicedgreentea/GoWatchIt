package models

type HomeAssistantScriptReq struct {
	EntityID string `json:"entity_id"`
}

type HomeAssistantNotificationReq struct {
	Message string `json:"message"`
}

type HomeAssistantWebhookPayload struct {
	
}

type HAEnvyResponse struct {
	EntityID    string     `json:"entity_id"`
	State       string     `json:"state"`
	Attributes  Attributes `json:"attributes"`
}
type Attributes struct {
	NoSignal bool `json:"no_signal"`
}