package plex

import (
	"fmt"
	"strings"
)

// parseXMLError extracts the error message from an XML response
func parseXMLError(err error, payload string) error {
	if err == nil {
		return nil
	}

	switch {
	case IsHTTPSToHTTPError(payload):
		return NewErrHTTPSToHTTP()
	case IsPlexUnauthorizedError(payload):
		return NewPlexUnauthorized()
	default:
		return fmt.Errorf("failed to parse XML response: %w - %s", err, payload)
	}
}

type ErrHTTPSToHTTP struct{}

// Error implements the error interface
func (e *ErrHTTPSToHTTP) Error() string {
	return "The plain HTTP request was sent to HTTPS port"
}

// NewErrHTTPSToHTTP creates a new ErrHTTPSToHTTP error
func NewErrHTTPSToHTTP() *ErrHTTPSToHTTP {
	return &ErrHTTPSToHTTP{}
}

// IsHTTPSToHTTPError checks if the error is ErrHTTPSToHTTP
func IsHTTPSToHTTPError(payload string) bool {
	return strings.Contains(payload, "The plain HTTP request was sent to HTTPS port")
}

type PlexUnauthorized struct{}

// Error implements the error interface
func (e *PlexUnauthorized) Error() string {
	return "Unauthorized"
}

// NewPlexUnauthorized creates a new PlexUnauthorized error
func NewPlexUnauthorized() *PlexUnauthorized {
	return &PlexUnauthorized{}
}

// IsHTTPSToHTTPError checks if the error is PlexUnauthorized
func IsPlexUnauthorizedError(payload string) bool {
	return strings.Contains(payload, "Unauthorized")
}
