package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/database"
	"github.com/iloveicedgreentea/go-plex/internal/mediaplayer"
	"github.com/iloveicedgreentea/go-plex/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// MockMediaClient implements MediaAPIClient for testing
type MockMediaClient struct {
	GetEditionCalled    bool
	GetAudioCodecCalled bool
	EditionToReturn     models.Edition
	CodecToReturn       models.CodecName
	ErrorToReturn       error
}

func (m *MockMediaClient) GetEdition(ctx context.Context, payload *models.Event) (models.Edition, error) {
	m.GetEditionCalled = true
	if m.ErrorToReturn != nil {
		return models.EditionNone, m.ErrorToReturn
	}
	return m.EditionToReturn, nil
}

func (m *MockMediaClient) GetAudioCodec(ctx context.Context, data *models.Event) (models.CodecName, error) {
	m.GetAudioCodecCalled = true
	if m.ErrorToReturn != nil {
		return "", m.ErrorToReturn
	}
	return m.CodecToReturn, nil
}

func (m *MockMediaClient) DoPlaybackAction(ctx context.Context, action models.Action) error {
	return m.ErrorToReturn
}

func setupTestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = database.RunMigrations(db)
	assert.NoError(t, err)

	err = config.InitConfig(db)
	assert.NoError(t, err)

	// Insert test EZBEQ config with all fields from the model
	_, err = db.Exec(`INSERT INTO EZBEQConfig (
        enabled, url, port, scheme, adjust_master_volume_with_profile,
        denon_ip, denon_port, dry_run, enable_tv_beq, notify_on_load,
        preferred_author, slots, stop_plex_if_mismatch,
        use_avr_codec_search, avr_brand, avr_url,
        loose_edition_matching, skip_edition_matching
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		true, "localhost", "8080", "http", false,
		"", "", false, false, false,
		"", "[]", false,
		false, "", "",
		false, false)
	assert.NoError(t, err)

	// Insert test Plex config with all fields from the model
	_, err = db.Exec(`INSERT INTO PlexConfig (
        url, port, scheme, enabled, device_uuid_filter,
        owner_name_filter, token, enable_trailer_support
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"localhost", "32400", "http", true, "",
		"", "test-token", false)
	assert.NoError(t, err)
	assert.True(t, config.IsPlexEnabled())
}

func TestGetClient(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name       string
		payload    models.Event
		wantErr    bool
		errMessage string
	}{
		{
			name: "valid_plex_event",
			payload: models.Event{
				ServerUUID: "fakeuuidtesting",
				EventType:  models.EventTypePlex,
				Metadata: models.Metadata{
					Type:  models.MediaTypeShow,
					Title: "Test Show",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_client_type",
			payload: models.Event{
				ServerUUID: "fakeuuidtesting",
				Client:     "not a media client", // This will fail type assertion
			},
			wantErr:    true,
			errMessage: "error checking client is MediaAPIClient",
		},
		{
			name: "unknown_server_type",
			payload: models.Event{
				ServerUUID: "unknown-server",
			},
			wantErr:    true,
			errMessage: "unsupported server type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := getClient(&tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Implements(t, (*mediaplayer.MediaAPIClient)(nil), client)
			}
		})
	}
}

func TestInitializeMediaClient(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name      string
		event     *models.Event
		expectErr bool
	}{
		{
			name: "valid_plex_server",
			event: &models.Event{
				ServerUUID: "plex-123",
				EventType:  models.EventTypePlex,
			},
			expectErr: false,
		},
		{
			name: "unknown_server",
			event: &models.Event{
				ServerUUID: "unknown-123",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := initializeMediaClient(tt.event)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Implements(t, (*mediaplayer.MediaAPIClient)(nil), client)
			}
		})
	}
}
