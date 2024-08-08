package plex

type plexItemTitle string

const (
	showItemTitle  plexItemTitle = "episode"
	movieItemTitle plexItemTitle = "movie"
)

type APIPath string

const (
	APIStatusSession APIPath = "/status/sessions"
)

