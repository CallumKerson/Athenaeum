package audiobooks

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

func TestAudiobooks_Persons(t *testing.T) {
	tests := []struct {
		name             string
		expectsNarrators bool
		audiobook        Audiobook
		expectedPersons  string
	}{
		{
			name:             "Single author",
			expectsNarrators: false,
			audiobook:        Audiobook{Authors: []string{"Ursula K LeGuin"}},
			expectedPersons:  "Ursula K LeGuin",
		},
		{
			name:             "Single narrator",
			expectsNarrators: true,
			audiobook:        Audiobook{Narrators: []string{"Kobna Holdbrook-Smith"}},
			expectedPersons:  "Kobna Holdbrook-Smith",
		},
		{
			name:             "Two authors",
			expectsNarrators: false,
			audiobook:        Audiobook{Authors: []string{"Amal El-Mohtar", "Max Gladstone"}},
			expectedPersons:  "Amal El-Mohtar & Max Gladstone",
		},
		{
			name:             "Two narrators",
			expectsNarrators: true,
			audiobook:        Audiobook{Narrators: []string{"Cynthia Farrell", "Emily Woo Zeller"}},
			expectedPersons:  "Cynthia Farrell & Emily Woo Zeller",
		},
		{
			name:             "Multiple narrators",
			expectsNarrators: true,
			audiobook: Audiobook{
				Narrators: []string{
					"Jay Snyder",
					"Brandon Rubin",
					"Fred Berman",
					"Lauren Fortgang",
					"Roger Clark",
					"Elizabeth Evans",
					"Tristan Morris",
				},
			},
			expectedPersons: "Jay Snyder, Brandon Rubin, Fred Berman, Lauren Fortgang, Roger Clark, Elizabeth Evans & Tristan Morris",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var personString string
			if testCase.expectsNarrators {
				personString = testCase.audiobook.GetNarrator()
			} else {
				personString = testCase.audiobook.GetAuthor()
			}
			assert.Equal(t, testCase.expectedPersons, personString)
		})
	}
}

func TestAudiobook_JSONMarshaling(t *testing.T) {
	book := Audiobook{
		Title:     "Test Audiobook",
		Authors:   []string{"Test Author"},
		Narrators: []string{"Test Narrator"},
		Genres:    []Genre{SciFi, Fantasy},
		Tags:      []string{"test", "audiobook"},
		Path:      "/test/path.m4b",
		Duration:  2 * time.Hour,
		FileSize:  1024 * 1024 * 100, // 100MB
		MIMEType:  "audio/mp4a-latm",
		Series: &Series{
			Title:    "Test Series",
			Sequence: Sequence{First: decimal.NewFromFloat(1.5)},
		},
		Description: &description.Description{
			Text:   "A test audiobook for unit testing",
			Format: description.Markdown,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(book)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Audiobook
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, book, unmarshaled)
}

func TestAudiobook_TOMLMarshaling(t *testing.T) {
	book := Audiobook{
		Title:     "Test Audiobook TOML",
		Authors:   []string{"TOML Author"},
		Narrators: []string{"TOML Narrator"},
		Genres:    []Genre{Biography, Historical},
		Tags:      []string{"toml", "test"},
		Path:      "/toml/test/path.m4b",
		Duration:  3*time.Hour + 30*time.Minute,
		FileSize:  1024 * 1024 * 200, // 200MB
		MIMEType:  "audio/mp4a-latm",
		Series: &Series{
			Title:    "TOML Test Series",
			Sequence: Sequence{First: decimal.NewFromInt(2)},
		},
		Description: &description.Description{
			Text:   "A test audiobook for TOML serialisation testing",
			Format: description.HTML,
		},
	}

	// Marshal to TOML
	data, err := toml.Marshal(book)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Audiobook
	err = toml.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, book, unmarshaled)
}

func TestAudiobook_EmptySlicesMarshaling(t *testing.T) {
	book := Audiobook{
		Title:     "Minimal Book",
		Authors:   []string{},
		Narrators: []string{},
		Genres:    []Genre{},
		Tags:      []string{},
	}

	// An empty slice does not survive a round trip through either format when the
	// field is omitempty — it is dropped on the way out and comes back nil. Nothing
	// distinguishes "no genres" from "an empty list of genres", so this is only
	// worth pinning down, not fixing.
	jsonData, err := json.Marshal(book)
	require.NoError(t, err)

	var jsonUnmarshaled Audiobook
	err = json.Unmarshal(jsonData, &jsonUnmarshaled)
	require.NoError(t, err)
	assert.Equal(t, []string{}, jsonUnmarshaled.Authors, "Authors is not omitempty in JSON")
	assert.Equal(t, []string{}, jsonUnmarshaled.Narrators, "Narrators is not omitempty in JSON")
	assert.Nil(t, jsonUnmarshaled.Genres)
	assert.Nil(t, jsonUnmarshaled.Tags)

	tomlData, err := toml.Marshal(book)
	require.NoError(t, err)

	var tomlUnmarshaled Audiobook
	err = toml.Unmarshal(tomlData, &tomlUnmarshaled)
	require.NoError(t, err)
	assert.Equal(t, []string{}, tomlUnmarshaled.Authors, "Authors is not omitempty in TOML")
	assert.Nil(t, tomlUnmarshaled.Narrators)
	assert.Nil(t, tomlUnmarshaled.Genres)
	assert.Nil(t, tomlUnmarshaled.Tags)

	// Everything that is not a slice still survives both round trips.
	assert.Equal(t, book.Title, jsonUnmarshaled.Title)
	assert.Equal(t, book.Title, tomlUnmarshaled.Title)
}

func TestSeries_SequenceJSON(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		expected string
	}{
		{"Whole number", "2", `{"sequence":"2","title":"Test Series"}`},
		{"Decimal", "1.5", `{"sequence":"1.5","title":"Test Series"}`},
		{"Omnibus", "1-3", `{"sequence":"1-3","title":"Test Series"}`},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sequence, err := ParseSequence(testCase.sequence)
			require.NoError(t, err)
			series := Series{Title: "Test Series", Sequence: sequence}

			data, err := json.Marshal(series)
			require.NoError(t, err)
			assert.JSONEq(t, testCase.expected, string(data))

			var unmarshaled Series
			require.NoError(t, json.Unmarshal(data, &unmarshaled))
			assert.Equal(t, series, unmarshaled)
		})
	}
}

func TestGetPersonsString_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		persons  []string
		expected string
	}{
		{"Empty slice", []string{}, ""},
		{"Nil slice", nil, ""},
		{"Single empty string", []string{""}, ""},
		{"Multiple empty strings", []string{"", "", ""}, ",  & "},
		{"Mix of empty and valid", []string{"", "Valid Author", ""}, ", Valid Author & "},
		{"Single person with spaces", []string{" Author Name "}, " Author Name "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPersonsString(tt.persons)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAuthor_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		book     Audiobook
		expected string
	}{
		{"Nil authors", Audiobook{Authors: nil}, ""},
		{"Empty authors", Audiobook{Authors: []string{}}, ""},
		{"Authors with empty strings", Audiobook{Authors: []string{"", "Valid Author"}}, " & Valid Author"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.book.GetAuthor()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetNarrator_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		book     Audiobook
		expected string
	}{
		{"Nil narrators", Audiobook{Narrators: nil}, ""},
		{"Empty narrators", Audiobook{Narrators: []string{}}, ""},
		{"Narrators with empty strings", Audiobook{Narrators: []string{"", "Valid Narrator"}}, " & Valid Narrator"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.book.GetNarrator()
			assert.Equal(t, tt.expected, result)
		})
	}
}
