package audiobooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Audiobook - representation of a book.
type Audiobook struct {
	Path             string        `json:"path"`
	Title            string        `json:"title"`
	Subtitle         string        `json:"subtitle"`
	Authors          []string      `json:"authors"`
	Description      *Description  `json:"description,omitempty"`
	ReleaseDate      *ReleaseDate  `json:"releaseDate,omitempty"`
	Genres           []Genre       `json:"genres,omitempty"`
	Series           *Series       `json:"series,omitempty"`
	AudiobookMediaID string        `json:"audiobookMediaId"`
	Narrators        []string      `json:"narrators"`
	Duration         time.Duration `json:"duration"`
	FileSize         uint64        `json:"fileSize"`
	MIMEType         string        `json:"mimeType"`
}

// Person - representation of a person, for example an author or audiobook narrator.
type Person string

// Series - representation of a series of books.
type Series struct {
	Sequence decimal.Decimal `json:"sequence"`
	Title    string          `json:"title"`
}

func NewBook(title string, description *Description, authors []string, releaseDate *ReleaseDate,
	genreList []Genre, series *Series) Audiobook {
	return Audiobook{
		Title:       title,
		Authors:     authors,
		Description: description,
		ReleaseDate: releaseDate,
		Genres:      genreList,
		Series:      series,
	}
}

// Description - representation of a blurb of a book.
type Description struct {
	Text   string `json:"text"`
	Format Format `json:"format,omitempty"`
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

type ReleaseDate struct {
	time.Time
}

const releaseDateLayout = "2006-01-02"

func NewReleaseDate(date string) (*ReleaseDate, error) {
	layout := "2006-01-02"

	dateTime, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}

	return &ReleaseDate{Time: dateTime}, nil
}

func (d *ReleaseDate) UnmarshalJSON(b []byte) (err error) {
	str := strings.Trim(string(b), "\"")
	if str == "null" {
		d.Time = time.Time{}
		return
	}

	d.Time, err = time.Parse(releaseDateLayout, str)
	return
}

func (d *ReleaseDate) MarshalJSON() ([]byte, error) {
	if d.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("%q", d.Time.Format(releaseDateLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (d *ReleaseDate) IsSet() bool {
	return d.UnixNano() != nilTime
}
