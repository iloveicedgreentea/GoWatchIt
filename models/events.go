package models

import "encoding/json"

type Action string

const (
	ActionPlay  Action = "play"
	ActionPause Action = "pause"
	ActionStop  Action = "stop"
)

type IntOrString struct {
	IntValue    int
	StringValue string
}

func (ios *IntOrString) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &ios.IntValue); err == nil {
		return nil
	}

	if err := json.Unmarshal(data, &ios.StringValue); err != nil {
		return err
	}
	return nil
}

// Event is a generic container for events
type Event struct {
	Action Action `json:"action"`
	User   bool   `json:"user"`
	Owner  bool   `json:"owner"`
	// JF and plex have different account structures
	AccountID   IntOrString `json:"id"`
	ServerUUID  string      `json:"server_uuid"`
	PlayerUUID  string      `json:"player_uuid"`
	PlayerTitle string      `json:"player_title"`
	ServerTitle string      `json:"server_title"`
	PlayerIP    string      `json:"player_ip"`
	ServerIP    string      `json:"server_ip"`
	Metadata Metadata `json:"metadata"`
}

// Metadata is a generic container for media metadata
type Metadata struct {
	LibrarySectionType  string    `json:"librarySectionType"`
	Key                 string    `json:"key"`
	GUID                string    `json:"guid"`
	Type                MediaType `json:"type"`
	Title               string    `json:"title"`
	LibrarySectionTitle string    `json:"librarySectionTitle"`
	LibrarySectionID    int       `json:"librarySectionID"`
	LibrarySectionKey   string    `json:"librarySectionKey"`
	Year                int       `json:"year"`
	ItemID              string    `json:"ItemId"`
	IsPaused            bool      `json:"isPaused"`
	Codec               string    `json:"codec"`
	FileName            string    `json:"filename"`
}

type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeShow  MediaType = "episode"
)
