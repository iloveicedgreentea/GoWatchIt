package models

// PayloadType is an enum for the type of payload
type PayloadType int

const (
    JellyfinPayloadType PayloadType = iota  // JellyfinPayloadType is 0
    PlexPayloadType                         // PlexPayloadType is 1
)
// Allow generic payload for shared functions
type PayloadTypeUnion interface {
    JellyfinWebhook | PlexWebhookPayload
}

// WebhookPayload holds a payload that is either JellyfinWebhook or PlexWebhookPayload.
type WebhookPayload[T PayloadTypeUnion] struct {
    Type    PayloadType
    Payload T
}