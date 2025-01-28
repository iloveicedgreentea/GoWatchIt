package common

// import (
// 	// "os"
// 	"sync"
// 	"testing"

// 	// "fmt"

// 	// "github.com/iloveicedgreentea/go-plex/internal/avr"
// 	"github.com/iloveicedgreentea/go-plex/internal/config"
// 	// "github.com/iloveicedgreentea/go-plex/internal/ezbeq"
// 	"github.com/iloveicedgreentea/go-plex/internal/homeassistant"
// 	// "github.com/iloveicedgreentea/go-plex/internal/jellyfin"
// 	"github.com/iloveicedgreentea/go-plex/internal/plex"
// 	"github.com/iloveicedgreentea/go-plex/models"
// 	"github.com/stretchr/testify/assert"
// )

// func TestReadAndWait(t *testing.T) {
// 	waitTime := 30
// 	haClient := homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), config.GetString("homeAssistant.remoteentityname"))
// 	sync, err := readAttrAndWait(waitTime, "remote", "envy", &models.HAEnvyResponse{}, haClient)
// 	assert.NoError(t, err)
// 	t.Log(sync)
// }

// func TestWaitForHDMI(t *testing.T) {
// 	haClient := homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), config.GetString("homeAssistant.remoteentityname"))
// 	wg := new(sync.WaitGroup)
// 	wg.Add(1)
// 	skipActions := false
// 	plexClient := plex.NewClient(config.GetString("plex.url"), config.GetString("plex.port"))
// 	t.Log(plexClient.MachineID)
// 	WaitForHDMISync(wg, &skipActions, haClient, plexClient)
// }

// type TestDeps struct {
// 	BeqClient *ezbeq.BeqClient
// 	HaClient  *homeassistant.HomeAssistantClient
// 	AvrClient avr.AVRClient
// }

// func InitializeTestDependencies() TestDeps {
// 	beq, err := ezbeq.NewClient(config.GetString("ezbeq.url"), config.GetString("ezbeq.port"))
// 	if err != nil {
// 		log.Fatalf("Error initializing BeqClient: %s", err)
// 	}
// 	ha := homeassistant.NewClient(config.GetString("homeAssistant.url"), config.GetString("homeAssistant.port"), config.GetString("homeAssistant.token"), "")
// 	avr := avr.GetAVRClient(config.GetString("ezbeq.avrUrl"))
// 	return TestDeps{
// 		BeqClient: beq,
// 		HaClient:  ha,
// 		AvrClient: avr,
// 	}
// }

// func InitializeMediaClient(payload interface{}) Client {
// 	switch p := payload.(type) {
// 	case models.JellyfinWebhook:
// 		return nil
// 		// Initialize and return a Jellyfin client
// 		return jellyfin.NewClient(config.GetString("jellyfin.url"), config.GetString("jellyfin.port"), config.GetString(""), config.GetString("")) // Replace with the actual constructor function
// 	case models.PlexWebhookPayload:
// 		// Initialize and return a Plex client
// 		return plex.NewClient(config.GetString("plex.url"), config.GetString("plex.port"), config.GetString("plex.playermachineidentifier"), config.GetString("plex.playerip")) // Replace with the actual constructor function
// 	default:
// 		// Handle unsupported payload types
// 		log.Fatalf("Unsupported payload type: %T", p)
// 		return nil
// 	}
// }

// func getJFPayload() (out models.JellyfinWebhook, err error) {
// 	// open testing file
// 	jsonFile, err := os.ReadFile("testdata/jf_pause.json")
// 	if err != nil {
// 		return out, err
// 	}

// 	// mock request
// 	payload, err := decodeJfWebhook(jsonFile)
// 	if err != nil {
// 		return out, err
// 	}
// 	log.Debugf("JF Test Payload: %#v", payload)
// 	return payload, nil

// }

// func TestJFPayload(t *testing.T) {
// 	payload, err := getJFPayload()
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, payload.Year)
// 	assert.NotEmpty(t, payload.ItemID)
// 	assert.NotEmpty(t, payload.ItemType)
// }

// func TestCommonPlay(t *testing.T) {
// 	deps := InitializeTestDependencies()

// 	// Example payload data - ensure these are suitable for real clients
// 	jellyfinPayload, err := getJFPayload()
// 	if err != nil {
// 		t.Fatalf("Error getting Jellyfin payload: %v", err)
// 	}
// 	plexPayload, err := getPlexWebhook()
// 	if err != nil {
// 		t.Fatalf("Error getting Plex payload: %v", err)
// 	}

// 	testCases := []struct {
// 		description string
// 		payload     interface{}
// 	}{
// 		{
// 			description: "Handle Jellyfin Payload",
// 			payload:     jellyfinPayload,
// 		},
// 		{
// 			description: "Handle Plex Payload",
// 			payload:     plexPayload,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.description, func(t *testing.T) {
// 			var wg sync.WaitGroup
// 			skipActions := false
// 			searchRequest := new(models.BeqSearchRequest)
// 			mediaClient := InitializeMediaClient(tc.payload)
// 			if mediaClient == nil {
// 				t.Fatalf("Failed to initialize media client for payload: %v", tc.payload)
// 			}

// 			// Call the function under test with real clients
// 			commonPlay(deps.BeqClient, deps.HaClient, mediaClient, deps.AvrClient, tc.payload, searchRequest, &skipActions, &wg)

// 			// Assertions to verify the behavior of commonPlay with real clients
// 			assert.Equal(t, true, skipActions, "skipActions should be set to true after commonPlay")
// 			// Add more assertions as needed based on the expected behavior of commonPlay and real clients
// 		})
// 	}
// }
