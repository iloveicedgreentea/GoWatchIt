package jellyfin

import (
	"net/http"
)

func IsValidWebhook(s *JellyfinWebhook) bool {
	return s.ItemID != "" && s.DeviceID != "" && s.DeviceName != "" && s.ItemType != "" && s.NotificationType != ""
}

func IsJellyfinWebhook(req *http.Request) bool {
	// TODO: look at a webhook request and dump headers to validate
	return true
}
