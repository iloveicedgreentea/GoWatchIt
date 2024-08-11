package actions

// all common actions
import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/internal/mqtt"
	"github.com/iloveicedgreentea/go-plex/models"
)

// trigger HA for MV change per type
func ChangeMasterVolume(ctx context.Context, mediaType models.MediaType) {
	if config.IsHomeAssistantTriggerAVRMasterVolumeChangeOnEvent() {
		log := logger.GetLoggerFromContext(ctx)
		log.Debug("Changing volume")
		err := mqtt.Publish([]byte(fmt.Sprintf("{\"type\":\"%s\"}", mediaType)), config.GetString("mqtt.topicvolume"))
		if err != nil {
			log.Error("Error changing volume", slog.String("error", err.Error()))
		}
	}
}

// trigger HA for light change given entity and desired state
func ChangeLight(ctx context.Context, state string) {
	if config.IsHomeAssistantTriggerLightsOnEvent() {
		log := logger.GetLoggerFromContext(ctx)
		log.Debug("Changing light")
		err := mqtt.Publish([]byte(fmt.Sprintf("{\"state\":\"%s\"}", state)), config.GetString("mqtt.topiclights"))
		if err != nil {
			log.Error("Error changing light", slog.String("error", err.Error()))
		}
	}
}
