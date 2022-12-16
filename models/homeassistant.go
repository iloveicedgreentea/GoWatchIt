package models

type HomeAssistantScriptReq struct {
	EntityID string `json:"entity_id"`
}

type HomeAssistantNotificationReq struct {
	Message string `json:"message"`
}

type HomeAssistantWebhookPayload struct {
	
}