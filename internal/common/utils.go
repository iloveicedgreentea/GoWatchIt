package common

import (
	"strings"
)

func InsensitiveContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}
