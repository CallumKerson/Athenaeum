package audiobooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Format - representation of formats of texts.
type Format uint8

const (
	Undefined Format = iota
	Plain
	Markdown
	HTML
)

var (
	errParsingFormat = errors.New("cannot parse Format")
	textFormatName   = map[uint8]string{
		1: "Plain",
		2: "Markdown",
		3: "HTML",
	}
	textFormatValue = map[string]uint8{
		"plain":    1,
		"markdown": 2,
		"md":       2,
		"html":     3,
	}
)

// String allows Format to implement fmt.Stringer.
func (g Format) String() string {
	return textFormatName[uint8(g)]
}

// Convert a string to a Format, returns an error if the string is unknown.
func ParseFormat(s string) (Format, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	value, ok := textFormatValue[s]
	if !ok {
		return Format(0), fmt.Errorf("%w: %q is not a valid textFormat", errParsingFormat, s)
	}

	return Format(value), nil
}

// MarshalJSON allows compatibility with marshalling JSON.
func (g Format) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON allows compatibility with unmarshalling JSON.
func (g *Format) UnmarshalJSON(data []byte) (err error) {
	var textFormat string

	err = json.Unmarshal(data, &textFormat)
	if err != nil {
		return err
	}

	if *g, err = ParseFormat(textFormat); err != nil {
		return err
	}

	return nil
}
