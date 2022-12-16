package models

// PlexPayload is a typical payload for webhooks
type PlexWebhookPayload struct {
	Event    string   `json:"event"`
	User     bool     `json:"user"`
	Owner    bool     `json:"owner"`
	Account  Account  `json:"Account"`
	Server   Server   `json:"Server"`
	Player   Player   `json:"Player"`
	Metadata Metadata `json:"Metadata"`
}
type Account struct {
	ID    int    `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
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
type Genre struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	Count  int    `json:"count"`
}
type Director struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
}
type Writer struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	Count  int    `json:"count,omitempty"`
}
type Producer struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	Count  int    `json:"count"`
}
type Country struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	Count  int    `json:"count"`
}
type GUID0 struct {
	ID string `json:"id"`
}
type Rating0 struct {
	Image string  `json:"image"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
	Count int     `json:"count"`
}
type Collection struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	Count  int    `json:"count"`
	GUID   string `json:"guid"`
}
type Role struct {
	ID     int    `json:"id"`
	Filter string `json:"filter"`
	Tag    string `json:"tag"`
	TagKey string `json:"tagKey"`
	Count  int    `json:"count,omitempty"`
	Role   string `json:"role"`
	Thumb  string `json:"thumb,omitempty"`
}
type Metadata struct {
	LibrarySectionType    string       `json:"librarySectionType"`
	RatingKey             string       `json:"ratingKey"`
	Key                   string       `json:"key"`
	GUID                  string       `json:"guid"`
	Studio                string       `json:"studio"`
	Type                  string       `json:"type"`
	Title                 string       `json:"title"`
	LibrarySectionTitle   string       `json:"librarySectionTitle"`
	LibrarySectionID      int          `json:"librarySectionID"`
	LibrarySectionKey     string       `json:"librarySectionKey"`
	ContentRating         string       `json:"contentRating"`
	Summary               string       `json:"summary"`
	Rating                float64      `json:"rating"`
	AudienceRating        float64      `json:"audienceRating"`
	ViewCount             int          `json:"viewCount"`
	SkipCount             int          `json:"skipCount"`
	LastViewedAt          int          `json:"lastViewedAt"`
	Year                  int          `json:"year"`
	Tagline               string       `json:"tagline"`
	Thumb                 string       `json:"thumb"`
	Art                   string       `json:"art"`
	Duration              int          `json:"duration"`
	OriginallyAvailableAt string       `json:"originallyAvailableAt"`
	AddedAt               int          `json:"addedAt"`
	UpdatedAt             int          `json:"updatedAt"`
	AudienceRatingImage   string       `json:"audienceRatingImage"`
	ChapterSource         string       `json:"chapterSource"`
	PrimaryExtraKey       string       `json:"primaryExtraKey"`
	RatingImage           string       `json:"ratingImage"`
	Genre                 []Genre      `json:"Genre"`
	Director              []Director   `json:"Director"`
	Writer                []Writer     `json:"Writer"`
	Producer              []Producer   `json:"Producer"`
	Country               []Country    `json:"Country"`
	GUID0                 []GUID0      `json:"Guid"`
	Rating0               []Rating0    `json:"Rating"`
	Collection            []Collection `json:"Collection"`
	Role                  []Role       `json:"Role"`
}
