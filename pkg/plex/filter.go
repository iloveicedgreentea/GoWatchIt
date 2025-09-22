package plex

import (
	"fmt"
	"strings"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
)

// check each one and return errors as needed
func CheckAllFilters(payload PlexWebhookPayload) (bool, error) {
	mediaType := payload.Metadata.Type
	if !checkItemTypeMatch(mediaType) {
		return false, fmt.Errorf("item type %s is not supported", mediaType)
	}
	uuidFilter := config.GetPlexDeviceUUIDFilter()
	uuid := payload.Player.UUID

	if !checkUUIDFilterMatch(uuid, uuidFilter) {
		return false, fmt.Errorf("uuid did not match expected %s but got %s", uuidFilter, uuid)
	}

	userIDFilter := config.GetPlexOwnerNameFilter()
	user := payload.Account.Title

	if !checkUserFilterMatch(user, userIDFilter) {
		return false, fmt.Errorf("user did not match expected %s but got %s", user, userIDFilter)
	}

	return true, nil
}

// check item type filter
func checkItemTypeMatch(mediaType string) bool {
	// check media type one of supported
	return strings.EqualFold(mediaType, string(MediaTypeMovie)) || strings.EqualFold(mediaType, string(MediaTypeShow))
}

// check if device uuid matches
func checkUUIDFilterMatch(uuid, uuidFilter string) bool {
	// check if device filter is set and if it matches the device UUID
	if uuidFilter == "" {
		return true
	}
	return strings.EqualFold(uuidFilter, uuid)
}

// check if user filter matches
func checkUserFilterMatch(username, userFilter string) bool {
	if userFilter == "" {
		return true
	}
	// only respond to events on a particular account
	// TODO: decodedPayload.Account.Title seems to always map to server owner not player account?
	return strings.EqualFold(username, userFilter)
}
