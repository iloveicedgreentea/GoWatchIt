package models

// MediaPayload is a union type that can hold payloads from different media players
type MediaPayload struct {
    PlexPayload    *PlexWebhookPayload
    JellyfinPayload *JellyfinWebhook
}

type DataMediaContainer struct {
    PlexPayload *MediaContainer
    JellyfinPayload *JellyfinMetadata
}