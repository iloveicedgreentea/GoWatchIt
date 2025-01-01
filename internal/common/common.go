package common

// all common actions
import (
	"context"
	"log/slog"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
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
func ReadAttrAndWait(ctx context.Context, waitTime int, entType models.HomeAssistantEntity, entName string, attrResp homeassistant.HAAttributeResponse, haClient *homeassistant.HomeAssistantClient) (bool, error) {
	var err error
	var isSignal bool
	var attributes models.Attributes
	log := logger.GetLoggerFromContext(ctx)

	// read attributes until its not nosignal
	for i := 0; i < waitTime; i++ {
		attributes, err = haClient.ReadAttributes(entName, attrResp, entType)
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
