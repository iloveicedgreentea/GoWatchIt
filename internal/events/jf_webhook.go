package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

func parseJellyfinWebhook(ctx context.Context, req *http.Request) (models.Event, error) {
	var webhook models.JellyfinWebhook
	log := logger.GetLoggerFromContext(ctx)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return models.Event{}, err
	}
	defer req.Body.Close()

	// Unmarshal the JSON into the webhook struct
	err = json.Unmarshal(body, &webhook)
	if err != nil {
		return models.Event{}, err
	}

	// Check if the request is a Jellyfin webhook
	if !isValidWebhook(webhook) {
		return models.Event{}, fmt.Errorf("failed to parse Jellyfin webhook due to missing fields: %#v", webhook)
	}
	var action models.Action
	switch webhook.NotificationType {
	case "PlaybackStart":
		action = models.ActionPlay
	case "PlaybackStop":
		action = models.ActionStop
		// TODO: is paused and stuff
	}
	paused, err := strconv.ParseBool(webhook.IsPaused)
	if err != nil {
		log.Error("Failed to parse isPaused value",
			slog.Any("error", err),
		)
	}
	year, err := strconv.Atoi(webhook.Year)
	if err != nil {
		log.Error("Failed to parse year value",
			slog.Any("error", err),
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

	return models.Event{
		Action:      action,
		AccountID:   models.IntOrString{StringValue: webhook.UserID},
		PlayerUUID:  webhook.DeviceID,
		PlayerTitle: webhook.DeviceName,
		Metadata: models.Metadata{
			Key:      webhook.ItemID,
			Type:     models.MediaType(webhook.ItemType),
			IsPaused: paused,
			Year:     year,
			// TODO: tmdb for jellyfin
		},
	}, nil
}

func isValidWebhook(s models.JellyfinWebhook) bool {
	return s.ItemID != "" && s.DeviceID != "" && s.DeviceName != "" && s.ItemType != "" && s.NotificationType != ""
}

func isJellyfinWebhook(req *http.Request) bool {
	// TODO: look at a webhook request and dump headers to validate
	return true
}
