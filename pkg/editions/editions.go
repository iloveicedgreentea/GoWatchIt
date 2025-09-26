package editions

import "strings"

// Edition represents the edition of the media being played
// Edition is based on EzBEQ/BEQ strings
// All users need to normalize to these in their own GetEdition mappers
type Edition string

const (
	EditionExtended       Edition = "Extended"
	EditionUnrated        Edition = "Unrated"
	EditionTheatrical     Edition = "Theatrical"
	EditionSpecialEdition Edition = "Special"
	EditionUltimate       Edition = "Ultimate"
	EditionDirectorsCut   Edition = "Director"
	EditionCriterion      Edition = "Criterion"
	EditionUnknown        Edition = "Unknown"
	EditionNone           Edition = "None"
)

// MapToEdition maps a string to a Edition
func MapSToEdition(edition string) Edition {
	s := strings.ToLower(edition)

	switch {
	case strings.Contains(s, "extended"):
		return EditionExtended
	case strings.Contains(s, "unrated"):
		return EditionUnrated
	// case strings.Contains(s, "theatrical"):
	// 	return EditionTheatrical
	case strings.Contains(s, "ultimate"):
		return EditionUltimate
	case strings.Contains(s, "director"):
		return EditionDirectorsCut
	case strings.Contains(s, "criterion"):
		return EditionCriterion
	case strings.Contains(s, "special"):
		return EditionSpecialEdition
	default:
		return EditionNone // IMO should include theatrical too
	}
}
