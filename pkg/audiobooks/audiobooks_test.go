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
		{name: "Single author", expectsNarrators: false, audiobook: Audiobook{Authors: []string{"Ursula K LeGuin"}}, expectedPersons: "Ursula K LeGuin"},
		{name: "Single narrator", expectsNarrators: true, audiobook: Audiobook{Narrators: []string{"Kobna Holdbrook-Smith"}}, expectedPersons: "Kobna Holdbrook-Smith"},
		{name: "Two authors", expectsNarrators: false, audiobook: Audiobook{Authors: []string{"Amal El-Mohtar", "Max Gladstone"}}, expectedPersons: "Amal El-Mohtar & Max Gladstone"},
		{name: "Two narrators", expectsNarrators: true, audiobook: Audiobook{Narrators: []string{"Cynthia Farrell", "Emily Woo Zeller"}}, expectedPersons: "Cynthia Farrell & Emily Woo Zeller"},
		{name: "Multiple narrators", expectsNarrators: true,
			audiobook: Audiobook{Narrators: []string{
				"Jay Snyder", "Brandon Rubin", "Fred Berman", "Lauren Fortgang", "Roger Clark", "Elizabeth Evans", "Tristan Morris"},
			},
			expectedPersons: "Jay Snyder, Brandon Rubin, Fred Berman, Lauren Fortgang, Roger Clark, Elizabeth Evans & Tristan Morris"},
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

func TestNewBook(t *testing.T) {
	title := "The Left Hand of Darkness"
	authors := []string{"Ursula K. Le Guin"}

	book := NewBook(title, nil, authors, nil, nil, nil)

	assert.Equal(t, title, book.Title)
	assert.Equal(t, authors, book.Authors)
	assert.Empty(t, book.Narrators)

	// Check that other fields are initialised to zero values
	assert.Empty(t, book.Path)
	assert.Empty(t, book.Genres)
	assert.Empty(t, book.Tags)
	assert.Zero(t, book.Duration)
	assert.Zero(t, book.FileSize)
	assert.Empty(t, book.MIMEType)
}

func TestNewBook_EmptyInputs(t *testing.T) {
	book := NewBook("", nil, nil, nil, nil, nil)

	assert.Empty(t, book.Title)
	assert.Nil(t, book.Authors)
	assert.Empty(t, book.Narrators)
}

func TestAudiobook_Equal(t *testing.T) {
	book1 := Audiobook{
		Title:     "Test Book",
		Authors:   []string{"Author One", "Author Two"},
		Narrators: []string{"Narrator One"},
		Genres:    []Genre{SciFi, Fantasy},
		Tags:      []string{"award-winner", "classic"},
		Path:      "/path/to/book.m4b",
		Duration:  10 * time.Hour,
		FileSize:  1024 * 1024 * 500, // 500MB
		MIMEType:  "audio/mp4a-latm",
		Series: &Series{
			Title:    "Test Series",
			Sequence: decimal.NewFromFloat(1.5),
		},
		Description: &description.Description{
			Text:   "A great book",
			Format: description.Markdown,
		},
	}

	book2 := book1 // Copy

	// Should be equal
	assert.True(t, Equal(&book1, &book2))
	assert.True(t, Equal(&book2, &book1))
}

func TestAudiobook_Equal_Different(t *testing.T) {
	book1 := Audiobook{
		Title:   "Book One",
		Authors: []string{"Author One"},
	}

	book2 := Audiobook{
		Title:   "Book Two",
		Authors: []string{"Author One"},
	}

	// Should not be equal due to different titles
	assert.False(t, Equal(&book1, &book2))
}

func TestAudiobook_Equal_NilPointers(t *testing.T) {
	book1 := Audiobook{Title: "Test"}

	// Both nil
	assert.True(t, Equal(nil, nil))

	// One nil, one not
	assert.False(t, Equal(&book1, nil))
	assert.False(t, Equal(nil, &book1))
}

func TestAudiobook_Equal_NilSeries(t *testing.T) {
	book1 := Audiobook{
		Title:  "Test Book",
		Series: nil,
	}

	book2 := Audiobook{
		Title:  "Test Book",
		Series: &Series{Title: "Test Series", Sequence: decimal.NewFromInt(1)},
	}

	book3 := Audiobook{
		Title:  "Test Book",
		Series: nil,
	}

	// book1 and book3 should be equal (both nil series)
	assert.True(t, Equal(&book1, &book3))

	// book1 and book2 should not be equal (nil vs non-nil series)
	assert.False(t, Equal(&book1, &book2))
	assert.False(t, Equal(&book2, &book1))
}

func TestAudiobook_Equal_DifferentSeries(t *testing.T) {
	book1 := Audiobook{
		Title:  "Test Book",
		Series: &Series{Title: "Series One", Sequence: decimal.NewFromInt(1)},
	}

	book2 := Audiobook{
		Title:  "Test Book",
		Series: &Series{Title: "Series Two", Sequence: decimal.NewFromInt(1)},
	}

	book3 := Audiobook{
		Title:  "Test Book",
		Series: &Series{Title: "Series One", Sequence: decimal.NewFromInt(2)},
	}

	// Different series name
	assert.False(t, Equal(&book1, &book2))

	// Different sequence
	assert.False(t, Equal(&book1, &book3))
}

func TestAudiobook_Equal_NilDescription(t *testing.T) {
	book1 := Audiobook{
		Title:       "Test Book",
		Description: nil,
	}

	book2 := Audiobook{
		Title:       "Test Book",
		Description: &description.Description{Text: "Test", Format: description.Plain},
	}

	book3 := Audiobook{
		Title:       "Test Book",
		Description: nil,
	}

	// book1 and book3 should be equal (both nil description)
	assert.True(t, Equal(&book1, &book3))

	// book1 and book2 should not be equal (nil vs non-nil description)
	assert.False(t, Equal(&book1, &book2))
	assert.False(t, Equal(&book2, &book1))
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
			Sequence: decimal.NewFromFloat(1.5),
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

	// Should be equal
	assert.True(t, Equal(&book, &unmarshaled))
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
			Sequence: decimal.NewFromInt(2),
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

	// Should be equal
	assert.True(t, Equal(&book, &unmarshaled))
}

func TestAudiobook_EmptySlicesMarshaling(t *testing.T) {
	book := Audiobook{
		Title:     "Minimal Book",
		Authors:   []string{},
		Narrators: []string{},
		Genres:    []Genre{},
		Tags:      []string{},
	}

	// JSON round trip
	jsonData, err := json.Marshal(book)
	require.NoError(t, err)

	var jsonUnmarshaled Audiobook
	err = json.Unmarshal(jsonData, &jsonUnmarshaled)
	require.NoError(t, err)
	assert.True(t, Equal(&book, &jsonUnmarshaled))

	// TOML round trip
	tomlData, err := toml.Marshal(book)
	require.NoError(t, err)

	var tomlUnmarshaled Audiobook
	err = toml.Unmarshal(tomlData, &tomlUnmarshaled)
	require.NoError(t, err)

	// Handle nil vs empty slice differences in TOML marshalling
	if tomlUnmarshaled.Authors == nil {
		tomlUnmarshaled.Authors = []string{}
	}
	if tomlUnmarshaled.Narrators == nil {
		tomlUnmarshaled.Narrators = []string{}
	}
	if tomlUnmarshaled.Genres == nil {
		tomlUnmarshaled.Genres = []Genre{}
	}
	if tomlUnmarshaled.Tags == nil {
		tomlUnmarshaled.Tags = []string{}
	}

	assert.True(t, Equal(&book, &tomlUnmarshaled))
}

func TestSeries_DecimalSequence(t *testing.T) {
	series := Series{
		Title:    "Test Series",
		Sequence: decimal.NewFromFloat(1.5),
	}

	// Test that decimal sequence is preserved in JSON
	data, err := json.Marshal(series)
	require.NoError(t, err)

	var unmarshaled Series
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, series.Title, unmarshaled.Title)
	assert.True(t, series.Sequence.Equal(unmarshaled.Sequence))
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
