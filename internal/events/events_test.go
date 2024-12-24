package events

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockMultipartRequest(rawBody string) *http.Request {
	req, _ := http.NewRequest("POST", "http://localhost:9999/plexwebhook", bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------9a74d1aa5a7ac807")
	return req
}

func createJellyfinWebhookTestRequest() *http.Request {
	// Create a sample JSON payload
	jsonPayload := []byte(`{
		"DeviceId": "123456789",
		"DeviceName": "Living Room TV",
		"ClientName": "Jellyfin Web",
		"UserId": "user123",
		"ItemId": "item789",
		"ItemType": "Movie",
		"NotificationType": "PlaybackStart",
		"Year": "2023",
		"PlayedToCompletion": "false",
		"IsPaused": "false"
	}`)

	// Create a new request with the JSON payload
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(jsonPayload))

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestRequestToEvent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Plex event", func(t *testing.T) {
		req := createMockMultipartRequest(plexPayload)
		event, err := RequestToEvent(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, models.ActionPlay, event.Action)
		assert.Equal(t, "/library/metadata/3019", event.Metadata.Key)
		assert.Equal(t, "Player", event.PlayerTitle)
		assert.Equal(t, "player-id", event.PlayerUUID)
	})

	t.Run("Jellyfin event", func(t *testing.T) {
		req := createJellyfinWebhookTestRequest()
		event, err := RequestToEvent(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, models.ActionPlay, event.Action)
		assert.Equal(t, "Living Room TV", event.PlayerTitle)
		assert.Equal(t, "123456789", event.PlayerUUID)
	})

	t.Run("Unsupported event", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "text/plain")
		_, err := RequestToEvent(ctx, req)
		assert.Error(t, err)
	})
}

func TestIsPlexType(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Valid Plex request", func(t *testing.T) {
		req := createMockMultipartRequest(plexPayload)
		assert.True(t, IsPlexType(ctx, req))
	})

	t.Run("Invalid Plex request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", http.NoBody)
		assert.False(t, IsPlexType(ctx, req))
	})

	t.Run("Non-multipart request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "application/json")
		assert.False(t, IsPlexType(ctx, req))
	})
}

func TestIsJellyfinType(t *testing.T) {
	t.Parallel()
	t.SkipNow() // TODO: skipped
	t.Run("Valid Jellyfin request", func(t *testing.T) {
		req := createJellyfinWebhookTestRequest()
		assert.True(t, IsJellyfinType(req))
	})

	t.Run("Invalid Jellyfin request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "text/plain")
		assert.False(t, IsJellyfinType(req))
	})

	t.Run("Non-JSON request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		assert.False(t, IsJellyfinType(req))
	})
}

func TestParseJellyfinWebhook(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Valid Jellyfin webhook", func(t *testing.T) {
		req := createJellyfinWebhookTestRequest()
		event, err := parseJellyfinWebhook(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, models.ActionPlay, event.Action)
		assert.Equal(t, "Living Room TV", event.PlayerTitle)
		assert.Equal(t, "123456789", event.PlayerUUID)
		assert.Equal(t, models.MediaType("Movie"), event.Metadata.Type)
		assert.Equal(t, 2023, event.Metadata.Year)
		assert.False(t, event.Metadata.IsPaused)
	})

	t.Run("Invalid Jellyfin webhook", func(t *testing.T) {
		invalidPayload := []byte(`{"InvalidField": "InvalidValue"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		_, err := parseJellyfinWebhook(ctx, req)
		assert.Error(t, err)
	})

	t.Run("Non-JSON payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("Not JSON"))
		req.Header.Set("Content-Type", "text/plain")
		_, err := parseJellyfinWebhook(ctx, req)
		assert.Error(t, err)
	})
}

func TestIsValidWebhook(t *testing.T) {
	t.Parallel()
	t.SkipNow() // TODO: skipped
	t.Run("Valid webhook", func(t *testing.T) {
		webhook := models.JellyfinWebhook{
			DeviceID:   "123",
			DeviceName: "TestDevice",
			ItemType:   "Movie",
		}
		assert.True(t, isValidWebhook(&webhook))
	})

	t.Run("Invalid webhook - empty struct", func(t *testing.T) {
		webhook := models.JellyfinWebhook{}
		assert.False(t, isValidWebhook(&webhook))
	})

	t.Run("Invalid webhook - missing required fields", func(t *testing.T) {
		webhook := models.JellyfinWebhook{
			DeviceID: "123",
		}
		assert.False(t, isValidWebhook(&webhook))
	})
}

func TestIsHomeassistantType(t *testing.T) {
	t.Parallel()
	t.Run("Not implemented", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webhook", http.NoBody)
		assert.False(t, IsHomeassistantType(req))
	})
}

// Mock data
const plexPayload = "--------------------------9a74d1aa5a7ac807\r\nContent-Disposition: form-data; name=\"payload\"\r\nContent-Type: application/json\r\n\r\n{\"event\":\"media.play\",\"user\":false,\"owner\":true,\"Account\":{\"id\":1234,\"thumb\":\"https://plex.tv/users/1234/avatar?c=1234\",\"title\":\"o\"},\"Server\":{\"title\":\"123\",\"uuid\":\"fakeuuidtesting\"},\"Player\":{\"local\":false,\"publicAddress\":\"192.168.1.1\",\"title\":\"Player\",\"uuid\":\"player-id\"},\"Metadata\":{\"librarySectionType\":\"show\",\"ratingKey\":\"3019\",\"key\":\"/library/metadata/3019\",\"parentRatingKey\":\"3009\",\"grandparentRatingKey\":\"2958\",\"guid\":\"plex://episode/5d9c12a208fddd001f318b56\",\"parentGuid\":\"plex://season/602e680b9b7e9c002d71a5e0\",\"grandparentGuid\":\"plex://show/5d9c086d2192ba001f3101c6\",\"type\":\"episode\",\"title\":\"The One Where Rachel Quits\",\"titleSort\":\"One Where Rachel Quits\",\"grandparentKey\":\"/library/metadata/2958\",\"parentKey\":\"/library/metadata/3009\",\"librarySectionTitle\":\"TV Shows\",\"librarySectionID\":2,\"librarySectionKey\":\"/library/sections/2\",\"grandparentTitle\":\"Friends\",\"parentTitle\":\"Season 3\",\"contentRating\":\"TV-14\",\"summary\":\"Rachel makes a rash decision after Gunther tells her she needs to be retrained.  Phoebe tries to help Joey when he gets a job selling Christmas trees.  And Ross accidentally breaks a girl's leg and tries to make it up to her. [Christmas Episode]\",\"index\":10,\"parentIndex\":3,\"audienceRating\":8.1,\"viewCount\":1,\"lastViewedAt\":1696800036,\"year\":1996,\"thumb\":\"/library/metadata/3019/thumb/1687815756\",\"art\":\"/library/metadata/2958/art/1695823887\",\"parentThumb\":\"/library/metadata/3009/thumb/1687815755\",\"grandparentThumb\":\"/library/metadata/2958/thumb/1695823887\",\"grandparentArt\":\"/library/metadata/2958/art/1695823887\",\"grandparentTheme\":\"/library/metadata/2958/theme/1695823887\",\"duration\":1320000,\"originallyAvailableAt\":\"1996-12-12\",\"addedAt\":1669266057,\"updatedAt\":1687815756,\"audienceRatingImage\":\"themoviedb://image.rating\",\"Guid\":[{\"id\":\"imdb://tt0583474\"},{\"id\":\"tmdb://86334\"},{\"id\":\"tvdb://303878\"}],\"Rating\":[{\"image\":\"themoviedb://image.rating\",\"value\":8.1,\"type\":\"audience\"}],\"Director\":[{\"id\":17953,\"filter\":\"director=17953\",\"tag\":\"Terry Hughes\",\"tagKey\":\"5d7768384de0ee001fccc190\",\"thumb\":\"https://image.tmdb.org/t/p/original/ffU0D0Yn6RIjdufcviD3e5tn7Hu.jpg\"}],\"Writer\":[{\"id\":17812,\"filter\":\"writer=17812\",\"tag\":\"Michael Curtis\",\"tagKey\":\"5e1635494c78f7003e7f44ba\"},{\"id\":17813,\"filter\":\"writer=17813\",\"tag\":\"Greg Malins\",\"tagKey\":\"5d7768760ea56a001e2a5a4c\",\"thumb\":\"https://metadata-static.plex.tv/b/people/b9a7830f2754cca651abbefe7d64fdd1.jpg\"}],\"Role\":[{\"id\":15772,\"filter\":\"actor=15772\",\"tag\":\"Mae Whitman\",\"tagKey\":\"5d776831103a2d001f566b27\",\"role\":\"Sarah Tuttle\",\"thumb\":\"https://metadata-static.plex.tv/8/people/848114147b5a88bf0a6fab205d9524dc.jpg\"},{\"id\":17648,\"filter\":\"actor=17648\",\"tag\":\"Shelley Berman\",\"tagKey\":\"5d776827103a2d001f564674\",\"role\":\"Mr. Kaplan Jr.\",\"thumb\":\"https://metadata-static.plex.tv/f/people/fa5ceaa3e423b6ec48b116f19cd2a625.jpg\"},{\"id\":17695,\"filter\":\"actor=17695\",\"tag\":\"Kyla Pratt\",\"tagKey\":\"5d77682d8718ba001e31307a\",\"role\":\"Charla Nichols\",\"thumb\":\"https://metadata-static.plex.tv/8/people/801a2079ce5ddc1000a0373f6d353f2c.jpg\"},{\"id\":17698,\"filter\":\"actor=17698\",\"tag\":\"Romy Rosemont\",\"tagKey\":\"5d77682b999c64001ec2d66b\",\"role\":\"Troop Leader\",\"thumb\":\"https://metadata-static.plex.tv/people/5d77682b999c64001ec2d66b.jpg\"},{\"id\":17955,\"filter\":\"actor=17955\",\"tag\":\"Sandra Gould\",\"tagKey\":\"5d77683aeb5d26001f1e1db5\",\"role\":\"Old Woman (voice)\",\"thumb\":\"https://metadata-static.plex.tv/c/people/c361791218a21938dfa1bfa7e379afd5.jpg\"},{\"id\":17612,\"filter\":\"actor=17612\",\"tag\":\"James Michael Tyler\",\"tagKey\":\"5d776b0ffb0d55001f55a7fb\",\"role\":\"Gunther\",\"thumb\":\"https://metadata-static.plex.tv/3/people/3570d61e44686f5d15724609d9e5d059.jpg\"}]}}\r\n--------------------------9a74d1aa5a7ac807--\r\n"
