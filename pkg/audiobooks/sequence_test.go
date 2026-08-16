package audiobooks

import (
	"encoding/json"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		isRange  bool
		noun     string
	}{
		{"Whole number", "1", "1", false, "Book"},
		{"Decimal", "1.5", "1.5", false, "Book"},
		{"Negative", "-1", "-1", false, "Book"},
		{"Surrounded by spaces", "  2  ", "2", false, "Book"},
		{"Omnibus", "1-3", "1-3", true, "Books"},
		{"Omnibus of decimals", "1.5-2.5", "1.5-2.5", true, "Books"},
		{"Omnibus spaced around the hyphen", "1 - 3", "1-3", true, "Books"},
		{"Omnibus starting below zero", "-1-3", "-1-3", true, "Books"},
		// Neither end is checked against the other, so an omnibus of one book and a
		// pair written backwards are both taken at face value.
		{"Omnibus of a single book", "1-1", "1-1", true, "Books"},
		{"Omnibus in descending order", "3-1", "3-1", true, "Books"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sequence, err := ParseSequence(testCase.input)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, sequence.String())
			assert.Equal(t, testCase.isRange, sequence.IsRange())
			assert.Equal(t, testCase.noun, sequence.Noun())
		})
	}
}

func TestParseSequence_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"Only spaces", "   "},
		{"Words", "book one"},
		{"Missing the end of the range", "1-"},
		{"Missing the start of the range", "- 3"},
		{"Three ends", "1-3-5"},
		{"Trailing text", "1 and 2"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := ParseSequence(testCase.input)
			require.Error(t, err)
			assert.ErrorIs(t, err, errParsingSequence)
		})
	}
}

func TestSequence_TextRoundTrip(t *testing.T) {
	for _, input := range []string{"1", "1.5", "-1", "1-3", "1.5-2.5"} {
		t.Run(input, func(t *testing.T) {
			var sequence Sequence
			require.NoError(t, sequence.UnmarshalText([]byte(input)))

			text, err := sequence.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, input, string(text))
		})
	}
}

func TestSequence_UnmarshalTextInvalid(t *testing.T) {
	var sequence Sequence
	require.Error(t, sequence.UnmarshalText([]byte("omnibus")))
}

// The written form is what the metadata files already contain, so a sequence has
// to survive TOML as the same string a hand-written file would hold.
func TestSequence_TOMLRoundTrip(t *testing.T) {
	for _, sequence := range []string{"1", "1.5", "1-3"} {
		t.Run(sequence, func(t *testing.T) {
			document := "Sequence = \"" + sequence + "\"\nTitle = \"Earthsea\"\n"

			var series Series
			require.NoError(t, toml.Unmarshal([]byte(document), &series))
			assert.Equal(t, sequence, series.Sequence.String())
			assert.Equal(t, "Earthsea", series.Title)

			// Marshalling picks its own quote style, so the value is what is checked.
			data, err := toml.Marshal(series)
			require.NoError(t, err)
			assert.Contains(t, string(data), sequence)

			var remarshaled Series
			require.NoError(t, toml.Unmarshal(data, &remarshaled))
			assert.Equal(t, series, remarshaled)
		})
	}
}

// A sequence that will not parse fails the whole book, which is what makes scan
// skip it rather than write a book with a missing place in its series.
func TestSequence_TOMLInvalid(t *testing.T) {
	var series Series
	require.Error(t, toml.Unmarshal([]byte("Sequence = \"one\"\nTitle = \"Earthsea\"\n"), &series))
}

func TestSequence_JSONOmnibus(t *testing.T) {
	sequence, err := ParseSequence("1-3")
	require.NoError(t, err)

	data, err := json.Marshal(sequence)
	require.NoError(t, err)
	assert.JSONEq(t, `"1-3"`, string(data))

	var unmarshaled Sequence
	require.NoError(t, json.Unmarshal(data, &unmarshaled))
	assert.Equal(t, sequence, unmarshaled)
}
