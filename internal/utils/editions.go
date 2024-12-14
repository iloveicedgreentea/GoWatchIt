package utils

import (
	"strings"

	"github.com/iloveicedgreentea/go-plex/models"
)

// MapToEdition maps a string to a models.Edition
func MapSToEdition(edition string) models.Edition {
	s := strings.ToLower(edition)

	switch {
	case strings.Contains(s, "extended"):
		return models.EditionExtended
	case strings.Contains(s, "unrated"):
		return models.EditionUnrated
	// case strings.Contains(s, "theatrical"):
	// 	return models.EditionTheatrical
	case strings.Contains(s, "ultimate"):
		return models.EditionUltimate
	case strings.Contains(s, "director"):
		return models.EditionDirectorsCut
	case strings.Contains(s, "criterion"):
		return models.EditionCriterion
	case strings.Contains(s, "special"):
		return models.EditionSpecialEdition
	default:
		return models.EditionNone // IMO should include theatrical too
	}
}
