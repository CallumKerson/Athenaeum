package audiobooks

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var errParsingSequence = errors.New("cannot parse sequence")

// rangeSeparator joins the ends of an omnibus, as in "1-3".
const rangeSeparator = "-"

// Sequence is a book's place in its series: a single entry — "1", or "1.5" for a
// novella that sits between two books — or, when Last is set, the span an
// omnibus covers, written "1-3".
//
// A range is whatever the metadata says it is. Nothing here insists the ends are
// in order or that they differ, because the only cost of an odd pair is an odd
// heading on one book's page, and rejecting it would drop that book from the
// library entirely — scan skips a book whose metadata will not parse.
type Sequence struct {
	First decimal.Decimal
	Last  *decimal.Decimal
}

// ParseSequence reads the written form: a decimal, or two joined by a hyphen.
// Whitespace around either end is ignored, so "1 - 3" is the range 1 to 3.
func ParseSequence(text string) (Sequence, error) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return Sequence{}, fmt.Errorf("%w: sequence is empty", errParsingSequence)
	}

	// A leading hyphen is the sign of the first entry rather than the separator, so
	// the search for the separator starts past it. Slicing by byte is safe: a
	// hyphen cannot appear inside a multi-byte rune.
	if index := strings.Index(trimmed[1:], rangeSeparator); index >= 0 {
		index++
		first, err := parseEnd(trimmed[:index], text)
		if err != nil {
			return Sequence{}, err
		}
		last, err := parseEnd(trimmed[index+len(rangeSeparator):], text)
		if err != nil {
			return Sequence{}, err
		}
		return Sequence{First: first, Last: &last}, nil
	}

	first, err := parseEnd(trimmed, text)
	if err != nil {
		return Sequence{}, err
	}
	return Sequence{First: first}, nil
}

// parseEnd reports the whole sequence in the error, not just the end that failed,
// because that is what the metadata file actually contains.
func parseEnd(end, sequence string) (decimal.Decimal, error) {
	value, err := decimal.NewFromString(strings.TrimSpace(end))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf(
			"%w: %q is not a number or a range of numbers",
			errParsingSequence,
			sequence,
		)
	}
	return value, nil
}

// IsRange reports whether the sequence covers more than one entry.
func (s Sequence) IsRange() bool {
	return s.Last != nil
}

// Noun is the word that introduces the sequence in a heading, so that an omnibus
// reads "Books 1-3" where a single book reads "Book 1".
func (s Sequence) Noun() string {
	if s.IsRange() {
		return "Books"
	}
	return "Book"
}

func (s Sequence) String() string {
	if !s.IsRange() {
		return s.First.String()
	}
	return s.First.String() + rangeSeparator + s.Last.String()
}

func (s Sequence) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Sequence) UnmarshalText(data []byte) (err error) {
	if *s, err = ParseSequence(string(data)); err != nil {
		return err
	}
	return nil
}
