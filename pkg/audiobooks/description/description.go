package description

// Description - representation of a blurb of a book.
type Description struct {
	Text   string `json:"text"`
	Format Format `json:"format,omitempty" toml:",omitempty"`
}
