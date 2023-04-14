package audiobooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

// Audiobook - representation of a book.
type Audiobook struct {
	Path        string                   `json:"path" toml:",omitempty"`
	Title       string                   `json:"title"`
	Subtitle    string                   `json:"subtitle" toml:",omitempty"`
	Authors     []string                 `json:"authors"`
	Description *description.Description `json:"description,omitempty" toml:",omitempty"`
	ReleaseDate *toml.LocalDate          `json:"releaseDate,omitempty" toml:",omitempty"`
	Genres      []Genre                  `json:"genres,omitempty" toml:",omitempty"`
	Series      *Series                  `json:"series,omitempty" toml:",omitempty"`
	Narrators   []string                 `json:"narrators" toml:",omitempty"`
	Tags        []string                 `json:"tags,omitempty" toml:",omitempty"`
	Duration    time.Duration            `json:"duration" toml:",omitempty"`
	FileSize    uint64                   `json:"fileSize" toml:",omitempty"`
	MIMEType    string                   `json:"mimeType" toml:",omitempty"`
}

// Person - representation of a person, for example an author or audiobook narrator.
type Person string

// Series - representation of a series of books.
type Series struct {
	Sequence decimal.Decimal `json:"sequence"`
	Title    string          `json:"title"`
}

func NewBook(title string, desc *description.Description, authors []string, releaseDate *toml.LocalDate,
	genreList []Genre, series *Series) Audiobook {
	return Audiobook{
		Title:       title,
		Authors:     authors,
		Description: desc,
		ReleaseDate: releaseDate,
		Genres:      genreList,
		Series:      series,
	}
}

func (b *Audiobook) GetAuthor() string {
	return GetPersonsString(b.Authors)
}

func (b *Audiobook) GetNarrator() string {
	return GetPersonsString(b.Narrators)
}

func GetPersonsString(persons []string) string {
	switch len(persons) {
	case 0:
		return ""
	case 1:
		return persons[0]
	default:
		return fmt.Sprintf("%s & %s", strings.Join(persons[:len(persons)-1], ", "), persons[len(persons)-1])
	}
}

func Equal(a, b *Audiobook) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return bytes.Equal(aBytes, bBytes)
}
