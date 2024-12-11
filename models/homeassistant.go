package models

type HomeAssistantEntity string

const (
	HomeAssistantEntityMediaPlayer  HomeAssistantEntity = "media_player"
	HomeAssistantEntityBinarySensor HomeAssistantEntity = "binary_sensor"
	HomeAssistantEntityRemote       HomeAssistantEntity = "remote"
)

type HomeAssistantMediaPlayerState string

const (
	HomeAssistantMediaPlayerStatePlaying HomeAssistantMediaPlayerState = "playing"
	HomeAssistantMediaPlayerStateStandby HomeAssistantMediaPlayerState = "standby"
	HomeAssistantMediaPlayerStatePaused  HomeAssistantMediaPlayerState = "paused"
	HomeAssistantMediaPlayerStateIdle    HomeAssistantMediaPlayerState = "idle"
	HomeAssistantMediaPlayerStateOff     HomeAssistantMediaPlayerState = "off"
	HomeAssistantMediaPlayerStateOn      HomeAssistantMediaPlayerState = "on"
)

type HomeAssistantScriptReq struct {
	EntityID string `json:"entity_id"`
}

type HomeAssistantNotificationReq struct {
	Message string `json:"message"`
}

type HomeAssistantWebhookPayload struct{}

// HAMediaPlayerResponse is a struct that contains the payload of a media player
type HAMediaPlayerResponse struct {
	EntityID   string                        `json:"entity_id"`
	State      HomeAssistantMediaPlayerState `json:"state"`
	Attributes Attributes                    `json:"attributes"`
}

// Attributes are media player attributes
type Attributes struct {
	SignalStatus     bool             `json:"is_signal"`
	MediaContentID   MediaContentID   `json:"media_content_id"`
	MediaContentType MediaContentType `json:"media_content_type"`
	MediaTitle       string           `json:"media_title"`
}

type MediaContentType string

const (
	MediaContentTypeMovie MediaContentType = "movie"
	MediaContentTypeShow  MediaContentType = "show"
)

type MediaContentID struct {
	IMDB string `json:"imdb"`
	TMDB string `json:"tmdb"`
	TVDB string `json:"tvdb"`
}

type HABinaryResponse struct {
	State string `json:"state"`
}

func (r *HABinaryResponse) GetState() string {
	return r.State
}

func (r *HABinaryResponse) GetSignalStatus() bool {
	return false // TODO: implement
}

func (r *HAMediaPlayerResponse) GetState() HomeAssistantMediaPlayerState {
	return r.State
}

func (r *HAMediaPlayerResponse) GetSignalStatus() bool {
	return r.Attributes.SignalStatus
}

func (r *HAMediaPlayerResponse) GetAttributes() Attributes {
	return r.Attributes
}
