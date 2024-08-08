package common

// PlaybackInterface is a common interface to control Client
func PlaybackInterface(action string, c Client) error {
	return c.DoPlaybackAction(action)
}

