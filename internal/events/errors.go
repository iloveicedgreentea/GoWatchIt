package events

type EventNotSupportedError struct {
	Message string
}

func (e EventNotSupportedError) Error() string {
	return e.Message
}
