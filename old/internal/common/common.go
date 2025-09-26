package common

// all common actions
import (
	"context"

	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/models"
)

// IsAtmosodecPlaying checks if Atmos (mapped and normalized from the player -> eg plex codec name into BEQ name) is being decoded instead of multi ch in (plex bug I believe)
func IsAtmosCodecPlaying(codec, expectedCodec string) (bool, error) {
	if codec == expectedCodec {
		return true, nil
	}

	return false, nil
}

// readAttrAndWait is a generic func to read attr from HA
func ReadAttrAndWait(ctx context.Context, waitTime int, entType models.HomeAssistantEntity, entName string, haClient *homeassistant.HomeAssistantClient) (bool, error) {
	return haClient.ReadAttrAndWait(ctx, waitTime, entType, entName)
}
