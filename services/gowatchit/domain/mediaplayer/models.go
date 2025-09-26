package mediaplayer

type Status struct{}

type MetaData struct{}

type Player string

const (
	PlayerPlex          Player = "Plex"
	PlayerJellyfin      Player = "Jellyfin"
	PlayerHomeAssistant Player = "HomeAssistant"
)
