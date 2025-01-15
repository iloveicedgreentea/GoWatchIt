package events

type EventNotSupportedError struct {
	Message string
}

func (e EventNotSupportedError) Error() string {
	return e.Message
}

type FilterDoesNotMatchError struct {
	Message string
}

func (e FilterDoesNotMatchError) Error() string {
	return e.Message
}
