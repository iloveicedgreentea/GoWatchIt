package beq

// BEQAuthor represents a BEQ author from the catalogue
type BEQAuthor struct {
	// Id Database ID (internal use only)
	Id *int64 `db:"id" json:"-"`

	// Name Author name from BEQ catalogue
	Name string `json:"name" db:"name"`

	// UpdatedAt Last updated timestamp
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}
