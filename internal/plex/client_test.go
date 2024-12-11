package plex

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	l "log"

	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iloveicedgreentea/go-plex/internal/database"
)

var (
	db     *sql.DB
	dbOnce sync.Once
	c      *PlexClient
)

// TestMain sets up the database and runs the tests
func TestMain(m *testing.M) {
	var code int
	dbOnce.Do(func() {
		// Setup code before tests
		var err error

		// Open SQLite database connection
		db, err = database.GetDB(":memory:")
		if err != nil {
			l.Fatalf("Failed to open database: %v", err)
		}

		// run migrations
		err = database.RunMigrations(db)
		if err != nil {
			l.Fatalf("Failed to run migrations: %v", err)
		}

		// Initialize the config with the database
		err = config.InitConfig(db)
		if err != nil {
			l.Fatalf("Failed to initialize config: %v", err)
		}

		cf := config.GetConfig()

		// populate test data
		plexCfg := models.PlexConfig{
			Enabled:              true,
			URL:                  os.Getenv("PLEX_URL"),
			Port:                 "443",
			Scheme:               "https",
			DeviceUUIDFilter:     "device_uuid",
			EnableTrailerSupport: false,
			OwnerNameFilter:      "owner_name",
		}
		err = cf.SaveConfig(&plexCfg)
		if err != nil {
			l.Fatalf("Failed to save ezbeq config: %v", err)
		}

		c, err = NewClient(config.GetPlexScheme(), config.GetPlexUrl(), config.GetPlexPort())
		if err != nil {
			l.Fatalf("Failed to create plex client: %v for values - %s, %s, %s", err, config.GetPlexScheme(), config.GetPlexUrl(), config.GetPlexPort())
		}

		// Run the tests
		code = m.Run()

		// Cleanup code after tests
		err = db.Close()
		if err != nil {
			l.Printf("Error closing database: %v", err)
		}
	})
	// Exit with the test result code
	os.Exit(code)
}

// test to ensure server is white listed
func TestGetPlexReq(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, err := c.makePlexReq(ctx, "/library/metadata/70390")
	require.NoError(t, err)
	res := string(d)

	if !assert.NotContains(t, res, "Unauthorized", "Client is not authorized in plex server") {
		t.Fatal(err)
	}
}

func TestGetEdition(t *testing.T) {
	t.Parallel()
	type test struct {
		dataObj  models.MediaContainer
		expected models.Edition
	}
	tests := []test{
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					EditionTitle: "Unrated",
				},
			},
			expected: models.EditionUnrated,
		},
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					EditionTitle: "unrated",
				},
			},
			expected: models.EditionUnrated,
		},
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					Media: struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					}{
						Part: struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						}{
							File: "Black Hawk Down (2001) {edition-Extended Edition} Remux 2160p HDR TrueHD Atmos",
						},
					},
				},
			},
			expected: models.EditionExtended,
		},
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					Media: struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					}{
						Part: struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						}{
							File: "Close Encounters of the Third Kind (1977) {edition-Director's Cut} Remux 2160p HDR DTS-HD MA",
						},
					},
				},
			},
			expected: models.EditionDirectorsCut,
		},
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					Media: struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					}{
						Part: struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						}{
							File: "Close Encounters of the Third Kind (1977) {edition-Special Edition} Remux 2160p HDR DTS-HD MA",
						},
					},
				},
			},
			expected: models.EditionSpecialEdition,
		},
		{
			dataObj: models.MediaContainer{
				Video: struct {
					RatingKey             string "xml:\"ratingKey,attr\""
					Key                   string "xml:\"key,attr\""
					ParentRatingKey       string "xml:\"parentRatingKey,attr\""
					GrandparentRatingKey  string "xml:\"grandparentRatingKey,attr\""
					AttrGuid              string "xml:\"guid,attr\""
					ParentGuid            string "xml:\"parentGuid,attr\""
					GrandparentGuid       string "xml:\"grandparentGuid,attr\""
					EditionTitle          string "xml:\"editionTitle,attr\""
					Type                  string "xml:\"type,attr\""
					Title                 string "xml:\"title,attr\""
					GrandparentKey        string "xml:\"grandparentKey,attr\""
					ParentKey             string "xml:\"parentKey,attr\""
					LibrarySectionTitle   string "xml:\"librarySectionTitle,attr\""
					LibrarySectionID      string "xml:\"librarySectionID,attr\""
					LibrarySectionKey     string "xml:\"librarySectionKey,attr\""
					GrandparentTitle      string "xml:\"grandparentTitle,attr\""
					ParentTitle           string "xml:\"parentTitle,attr\""
					ContentRating         string "xml:\"contentRating,attr\""
					Summary               string "xml:\"summary,attr\""
					Index                 string "xml:\"index,attr\""
					ParentIndex           string "xml:\"parentIndex,attr\""
					AudienceRating        string "xml:\"audienceRating,attr\""
					Thumb                 string "xml:\"thumb,attr\""
					Art                   string "xml:\"art,attr\""
					ParentThumb           string "xml:\"parentThumb,attr\""
					GrandparentThumb      string "xml:\"grandparentThumb,attr\""
					GrandparentArt        string "xml:\"grandparentArt,attr\""
					GrandparentTheme      string "xml:\"grandparentTheme,attr\""
					Duration              string "xml:\"duration,attr\""
					OriginallyAvailableAt string "xml:\"originallyAvailableAt,attr\""
					AddedAt               string "xml:\"addedAt,attr\""
					UpdatedAt             string "xml:\"updatedAt,attr\""
					AudienceRatingImage   string "xml:\"audienceRatingImage,attr\""
					Media                 struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					} "xml:\"Media\""
					Director struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Director\""
					Writer struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
					} "xml:\"Writer\""
					Guid []struct {
						ID string "xml:\"id,attr\""
					} "xml:\"Guid\""
					Rating struct {
						Image string "xml:\"image,attr\""
						Value string "xml:\"value,attr\""
						Type  string "xml:\"type,attr\""
					} "xml:\"Rating\""
					Role []struct {
						ID     string "xml:\"id,attr\""
						Filter string "xml:\"filter,attr\""
						Tag    string "xml:\"tag,attr\""
						TagKey string "xml:\"tagKey,attr\""
						Role   string "xml:\"role,attr\""
						Thumb  string "xml:\"thumb,attr\""
					} "xml:\"Role\""
				}{
					Media: struct {
						ID              string "xml:\"id,attr\""
						Duration        string "xml:\"duration,attr\""
						Bitrate         string "xml:\"bitrate,attr\""
						Width           string "xml:\"width,attr\""
						Height          string "xml:\"height,attr\""
						AspectRatio     string "xml:\"aspectRatio,attr\""
						AudioChannels   string "xml:\"audioChannels,attr\""
						AudioCodec      string "xml:\"audioCodec,attr\""
						VideoCodec      string "xml:\"videoCodec,attr\""
						VideoResolution string "xml:\"videoResolution,attr\""
						Container       string "xml:\"container,attr\""
						VideoFrameRate  string "xml:\"videoFrameRate,attr\""
						AudioProfile    string "xml:\"audioProfile,attr\""
						VideoProfile    string "xml:\"videoProfile,attr\""
						Part            struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						} "xml:\"Part\""
					}{
						Part: struct {
							ID           string "xml:\"id,attr\""
							Key          string "xml:\"key,attr\""
							Duration     string "xml:\"duration,attr\""
							File         string "xml:\"file,attr\""
							Size         string "xml:\"size,attr\""
							AudioProfile string "xml:\"audioProfile,attr\""
							Container    string "xml:\"container,attr\""
							VideoProfile string "xml:\"videoProfile,attr\""
							Stream       []struct {
								ID                   string "xml:\"id,attr\""
								StreamType           string "xml:\"streamType,attr\""
								Default              string "xml:\"default,attr\""
								Codec                string "xml:\"codec,attr\""
								Index                string "xml:\"index,attr\""
								Bitrate              string "xml:\"bitrate,attr\""
								BitDepth             string "xml:\"bitDepth,attr\""
								ChromaLocation       string "xml:\"chromaLocation,attr\""
								ChromaSubsampling    string "xml:\"chromaSubsampling,attr\""
								CodedHeight          string "xml:\"codedHeight,attr\""
								CodedWidth           string "xml:\"codedWidth,attr\""
								ColorRange           string "xml:\"colorRange,attr\""
								FrameRate            string "xml:\"frameRate,attr\""
								Height               string "xml:\"height,attr\""
								Level                string "xml:\"level,attr\""
								Profile              string "xml:\"profile,attr\""
								RefFrames            string "xml:\"refFrames,attr\""
								Width                string "xml:\"width,attr\""
								DisplayTitle         string "xml:\"displayTitle,attr\""
								ExtendedDisplayTitle string "xml:\"extendedDisplayTitle,attr\""
								Selected             string "xml:\"selected,attr\""
								Channels             string "xml:\"channels,attr\""
								Language             string "xml:\"language,attr\""
								LanguageTag          string "xml:\"languageTag,attr\""
								LanguageCode         string "xml:\"languageCode,attr\""
								AudioChannelLayout   string "xml:\"audioChannelLayout,attr\""
								SamplingRate         string "xml:\"samplingRate,attr\""
								Title                string "xml:\"title,attr\""
							} "xml:\"Stream\""
						}{
							File: "Close Encounters of the Third Kind (1977) Remux 2160p HDR DTS-HD MA",
						},
					},
				},
			},
			expected: models.EditionNone,
		},
	}

	for _, test := range tests {
		edition, err := getEdition(&test.dataObj)
		require.NoError(t, err)
		assert.Equal(t, test.expected, edition, "Expected: ", string(test.expected), "Got: ", edition, "for ", test.dataObj.Video.Media.Part.File, test.dataObj.Video.EditionTitle)
	}
}

func TestGetMediaData(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// no time to die
	event := &models.Event{
		Metadata: models.Metadata{
			Key: "/library/metadata/70390",
		},
	}
	med, err := c.getMediaData(ctx, event)
	require.NoError(t, err)

	assert.Equal(t, "Atmos", med.Video.Media.AudioCodec)
}

func TestGetCodecFromSession(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()

	codec, err := c.GetCodecFromSession(ctx, config.GetPlexDeviceUUIDFilter())
	require.NoError(t, err)

	t.Log(codec)
}

type codecTest struct {
	codec     string
	fullcodec string
	expected  string
}

func TestMapCodecs(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	ctx := context.Background()
	tests := []codecTest{
		{
			codec:     "EAC3",
			fullcodec: "EAC3",
			expected:  "DD+",
		},
		{
			codec:     "EAC3 5.1",
			fullcodec: "German (German EAC3 5.1)",
			expected:  "DD+Atmos5.1Maybe",
		},
		{
			codec:     "DDP 5.1 Atmos",
			fullcodec: "DDP 5.1 Atmos (Engelsk EAC3)",
			expected:  "DD+ Atmos",
		},
		{
			codec:     "English (TRUEHD 7.1)",
			fullcodec: "Surround 7.1 (English TRUEHD)",
			expected:  "AtmosMaybe",
		},
		{
			codec:     "English (TRUEHD 5.1)",
			fullcodec: "Dolby TrueHD Audio / 5.1 / 48 kHz / 1541 kbps / 16-bit (English)",
			expected:  "TrueHD 5.1",
		},
		{
			codec:     "English (DTS-HD MA 5.1)",
			fullcodec: "DTS-HD Master Audio / 5.1 / 48 kHz / 3887 kbps / 24-bit (English)",
			expected:  "DTS-HD MA 5.1",
		},
		{
			codec:     "English (TRUEHD 7.1)",
			fullcodec: "TrueHD Atmos 7.1 (English)",
			expected:  "Atmos",
		},
		{
			codec:     "English (DTS-HD MA 7.1)",
			fullcodec: "DTS:X / 7.1 / 48 kHz / 4213 kbps / 24-bit (English DTS-HD MA)",
			expected:  "DTS-X",
		},
		// TODO: verify other codecs without using extended display title
	}
	// execute each test
	for _, test := range tests {
		s := MapPlexToBeqAudioCodec(ctx, test.codec, test.fullcodec)
		a.Equal(test.expected, s)
	}
}

// lets me print out every codec I have in a given library
// func getCodecTemp(c *PlexClient, libraryKey string) string {
// 	data, err := c.GetMediaData(ctx, libraryKey)
// 	if err != nil {
// 		return "fail"
// 	}
// 	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
// 	// loop instead of index because of edge case with two or more video streams
// 	for _, val := range data.Video.Media.Part.Stream {
// 		if val.StreamType == "2" {
// 			return fmt.Sprintf("%s --- %s \n", val.DisplayTitle, val.ExtendedDisplayTitle)
// 		}
// 	}

// 	return "fail"
// }
