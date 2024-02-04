package common

import (
	"sync"
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/avr"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/ezbeq"
	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
	"github.com/iloveicedgreentea/go-plex/internal/jellyfin"
	"github.com/iloveicedgreentea/go-plex/internal/plex"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

type TestDeps struct {
	BeqClient *ezbeq.BeqClient
	HaClient  *homeassistant.HomeAssistantClient
	AvrClient avr.AVRClient
}

func InitializeTestDependencies() TestDeps {
	beq, err := ezbeq.NewClient(config.GetString("ezbeq.url"), config.GetString("ezbeq.port"))
	if err != nil {
		log.Fatalf("Error initializing BeqClient: %s", err)
	}
	ha := homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), "")
	avr := avr.GetAVRClient(config.GetString("ezbeq.avrUrl"))
	return TestDeps{
		BeqClient: beq,
		HaClient:  ha,
		AvrClient: avr,
	}
}

func InitializeMediaClient(payload interface{}) Client {
	switch p := payload.(type) {
	case models.JellyfinWebhook:
		// Initialize and return a Jellyfin client
		return jellyfin.NewClient(config.GetString(""), config.GetString(""), config.GetString(""), config.GetString("")) // Replace with the actual constructor function
	case models.PlexWebhookPayload:
		// Initialize and return a Plex client
		return plex.NewClient(config.GetString(""), config.GetString(""), config.GetString(""), config.GetString("")) // Replace with the actual constructor function
	default:
		// Handle unsupported payload types
		log.Fatalf("Unsupported payload type: %T", p)
		return nil
	}
}

func TestCommonPlay(t *testing.T) {
	deps := InitializeTestDependencies()

	// Example payload data - ensure these are suitable for real clients
	// TODO: add test data
	jellyfinPayload := models.JellyfinWebhook{
		Year:     "2021",
		ItemType: "movie",
		/* Initialize with test data */
	}
	plexPayload := models.PlexWebhookPayload{
		Metadata: models.Metadata{
			Year:  2021,
			Type:  "movie",
			Title: "Test Movie",
		},
		/* Initialize with test data */
	}

	testCases := []struct {
		description string
		payload     interface{}
	}{
		{
			description: "Handle Jellyfin Payload",
			payload:     jellyfinPayload,
		},
		{
			description: "Handle Plex Payload",
			payload:     plexPayload,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var wg sync.WaitGroup
			skipActions := false
			searchRequest := new(models.SearchRequest)
			mediaClient := InitializeMediaClient(tc.payload)
            if mediaClient == nil {
                t.Fatalf("Failed to initialize media client for payload: %v", tc.payload)
            }

			// Call the function under test with real clients
			commonPlay(deps.BeqClient, deps.HaClient, mediaClient, deps.AvrClient, tc.payload, searchRequest, &skipActions, &wg)

			// Assertions to verify the behavior of commonPlay with real clients
			assert.Equal(t, true, skipActions, "skipActions should be set to true after commonPlay")
			// Add more assertions as needed based on the expected behavior of commonPlay and real clients
		})
	}
}
