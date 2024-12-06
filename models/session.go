package models

// MediaSession is a struct that contains the payload of a media session
type MediaSession struct {
	PlexPayload     *SessionMediaContainer
	JellyfinPayload *JellyfinSession
}
