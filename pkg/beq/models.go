package beq

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/iloveicedgreentea/gowatchit/pkg/codecs"
	"github.com/iloveicedgreentea/gowatchit/pkg/editions"
	"github.com/iloveicedgreentea/gowatchit/pkg/events"
)

// BEQPayload is used for searching aka loading
type BEQPayload struct {
	// normalized codec
	Codec      codecs.Codec
	DryrunMode bool
	Edition    editions.Edition
	EntryID    string
	// show, movie, etc
	MediaType       events.MediaType
	MVAdjust        float64
	PreferredAuthor string
	SkipSearch      bool
	Slots           []int32
	Title           string
	TMDB            string
	Year            int
}

type BeqCatalog struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	SortTitle  string   `json:"sortTitle"`
	Year       int      `json:"year"`
	AudioTypes []string `json:"audioTypes"`
	Digest     string   `json:"digest"`
	MvAdjust   float64  `json:"mvAdjust"`
	Edition    string   `json:"edition"`
	MovieDbID  string   `json:"theMovieDB"`
	Author     string   `json:"author"`
}

type BeqDevices struct {
	Name         string     `json:"name"`
	MasterVolume float64    `json:"mastervolume"`
	Mute         bool       `json:"mute"`
	Slots        []BeqSlots `json:"slots"`
}

type BeqSlots struct {
	ID     string  `json:"id"`
	Last   string  `json:"last"`
	Active bool    `json:"active"`
	Gain1  float64 `json:"gain1"`
	Gain2  float64 `json:"gain2"`
	Mute1  bool    `json:"mute1"`
	Mute2  bool    `json:"mute2"`
}

type BeqPatchV2 struct {
	Mute         bool      `json:"mute"`
	MasterVolume float64   `json:"masterVolume"`
	Slots        []SlotsV2 `json:"slots"`
}
type SlotsV2 struct {
	ID     string    `json:"id"`
	Active bool      `json:"active"`
	Gains  []float64 `json:"gains"`
	Mutes  []bool    `json:"mutes"`
	Entry  string    `json:"entry"`
}

type BeqPatchV1 struct {
	Mute         bool      `json:"mute"`
	MasterVolume float64   `json:"masterVolume"`
	Slots        []SlotsV1 `json:"slots"`
}
type SlotsV1 struct {
	ID     string    `json:"id"`
	Active bool      `json:"active"`
	Gains  []float64 `json:"gains"`
	Mutes  []bool    `json:"mutes"`
	Entry  string    `json:"entry"`
}

type BeqClient struct {
	Scheme              string
	ServerURL           string
	Port                string
	CurrentMasterVolume float64
	CurrentMediaType    string
	MuteStatus          bool
	MasterVolume        float64
	HTTPClient          *retryablehttp.Client
	DeviceInfo          []BeqDevices
}

const (
	// API Endpoints (prefix) - Add comments for V1 usage
	apiV1Prefix  = "/api/1"
	apiV2Prefix  = "/api/2"
	profileEmpty = "Empty"
	authorNone   = "none"
)
