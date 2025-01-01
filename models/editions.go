package models

// Edition is an enum for different editions of a movie
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
