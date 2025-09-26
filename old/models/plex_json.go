package models

// PlexPayload is a typical payload for webhooks
type PlexWebhookPayload struct {
	Event    string       `json:"event"`
	User     bool         `json:"user"`
	Owner    bool         `json:"owner"`
	Account  Account      `json:"Account"`
	Server   Server       `json:"Server"`
	Player   Player       `json:"Player"`
	Metadata PlexMetadata `json:"Metadata"`
}

type Account struct {
	ID    IntOrString `json:"id"`
	Title string      `json:"title"`
}
type Server struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}
type Player struct {
	Local         bool   `json:"local"`
	PublicAddress string `json:"publicAddress"`
	Title         string `json:"title"`
	UUID          string `json:"uuid"`
}

type GUID0 struct {
	ID string `json:"id"`
}

type PlexMetadata struct {
	LibrarySectionType string `json:"librarySectionType"`
	RatingKey          string `json:"ratingKey"`
	Key                string `json:"key"`
	Type               string `json:"type"`
	Title              string `json:"title"`
	Year               int    `json:"year"`
	// yes plex returns two fields with the same name
	Guid string `json:"guid"`
	// this is where tmbd is
	GUID0               []GUID0 `json:"Guid"`
	LibrarySectionTitle string  `json:"librarySectionTitle"`
	LibrarySectionID    int     `json:"librarySectionID"`
	LibrarySectionKey   string  `json:"librarySectionKey"`
}
