package plex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemTypeMatchCheck(t *testing.T) {
	type test struct {
		payload  string
		expected bool
	}

	tt := []test{
		{
			payload:  string(MediaTypeMovie),
			expected: true,
		},
		{
			payload:  string(MediaTypeShow),
			expected: true,
		},
		{
			payload:  "other",
			expected: false,
		},
	}

	for _, p := range tt {
		t.Run(p.payload, func(t *testing.T) {
			if p.expected {
				assert.True(t, checkItemTypeMatch(p.payload))
			} else {
				assert.False(t, checkItemTypeMatch(p.payload))
			}
		})
	}
}

func TestUUIDCheck(t *testing.T) {
	type test struct {
		uuid     string
		filter   string
		expected bool
	}

	tt := []test{
		{
			uuid:     "12345",
			filter:   "12345",
			expected: true,
		},
		{
			uuid:     "12345",
			filter:   "1234",
			expected: false,
		},
		{
			uuid:     "12345",
			filter:   "",
			expected: true,
		},
	}

	for _, p := range tt {
		t.Run(p.uuid, func(t *testing.T) {
			if p.expected {
				assert.True(t, checkUUIDFilterMatch(p.uuid, p.filter))
			} else {
				assert.False(t, checkUUIDFilterMatch(p.uuid, p.filter))
			}
		})
	}
}

func TestUserCheck(t *testing.T) {
	type test struct {
		user     string
		filter   string
		expected bool
	}

	tt := []test{
		{
			user:     "abc",
			filter:   "abc",
			expected: true,
		},
		{
			user:     "abc",
			filter:   "1234",
			expected: false,
		},
		{
			user:     "abc",
			filter:   "",
			expected: true,
		},
	}

	for _, p := range tt {
		t.Run(p.user, func(t *testing.T) {
			if p.expected {
				assert.True(t, checkUserFilterMatch(p.user, p.filter))
			} else {
				assert.False(t, checkUserFilterMatch(p.user, p.filter))
			}
		})
	}
}
