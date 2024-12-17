package models

// MediaPayload is a union type that can hold payloads from different media players for webhooks
type MediaPayload struct {
	PlexPayload     *PlexWebhookPayload
	JellyfinPayload *JellyfinWebhook
}

// DataMediaContainer is a union type that can hold payloads from different media servers
type DataMediaContainer struct {
	PlexPayload     *MediaContainer
	JellyfinPayload *JellyfinMetadata
}
