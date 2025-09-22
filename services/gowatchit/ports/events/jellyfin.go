package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/iloveicedgreentea/gowatchit/pkg/jellyfin"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/events"
	"go.uber.org/zap"
)

// TODO: output a JF struct to be converted later
func ParseJellyfinWebhook(ctx context.Context, req *http.Request) (events.Event, error) {
	var webhook jellyfin.JellyfinWebhook
	log := logger.GetLoggerFromContext(ctx)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return events.Event{}, err
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Error("Failed to close request body",
				zap.Error(err),
			)
		}
	}()

	// Unmarshal the JSON into the webhook struct
	err = json.Unmarshal(body, &webhook)
	if err != nil {
		return events.Event{}, err
	}

	// Check if the request is a Jellyfin webhook
	if !jellyfin.IsValidWebhook(&webhook) {
		return events.Event{}, fmt.Errorf("failed to parse Jellyfin webhook due to missing fields: %#v", webhook)
	}

	var action events.Action
	switch webhook.NotificationType {
	case string(jellyfin.ActionStart):
		action = events.ActionPlay
	case string(jellyfin.ActionStop):
		action = events.ActionStop
		// JF does not send pause/resume events only "progress" events which include pause/resume status
	}
	paused, err := strconv.ParseBool(webhook.IsPaused)
	if err != nil {
		log.Error("Failed to parse isPaused value",
			zap.Error(err),
		)
	}
	year, err := strconv.Atoi(webhook.Year)
	if err != nil {
		log.Error("Failed to parse year value",
			zap.Error(err),
		)
	}
	// TODO: call client for JellyfinMetadata
	// urls := webhook.ExternalUrls
	// log.Debugf("External urls: %#v", urls)
	// for _, u := range urls {
	// 	if u.Name == "TheMovieDb" {
	// 		s := strings.Replace(u.URL, "https://www.themoviedb.org/", "", -1)
	// 		// extract the numbers
	// 		re, err := regexp.Compile(`\d+$`)
	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		return re.FindString(s), nil
	// 	}
	// }

	return events.Event{
		Action:      action,
		AccountID:   webhook.UserID,
		PlayerUUID:  webhook.DeviceID,
		PlayerTitle: webhook.DeviceName,
		EventType:   events.EventTypeJellyfin,
		Metadata: events.Metadata{
			Key:      webhook.ItemID,
			Type:     events.MediaType(webhook.ItemType),
			IsPaused: paused,
			Year:     year,
			// TODO: tmdb for jellyfin
		},
	}, nil
}
