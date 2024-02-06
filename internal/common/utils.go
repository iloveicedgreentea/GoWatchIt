package common

import (
	"strings"
)

func InsensitiveContains(s string, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}