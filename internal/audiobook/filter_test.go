package audiobook

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestFilters(t *testing.T) {
	tests := []struct {
		name          string
		filter        Filter
		book          audiobooks.Audiobook
		expectedMatch bool
	}{
		{name: "match author", filter: AuthorFilter("Max Gladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match author ignoring case", filter: AuthorFilter("Max gladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match author ignoring whitespace", filter: AuthorFilter("MaxGladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match author ignoring case and whitespace", filter: AuthorFilter("maxgladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match author", filter: AuthorFilter("Max Gladstone"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring case", filter: NarratorFilter("Emily woo zeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring whitespace", filter: NarratorFilter("EmilyWooZeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring case and whitespace", filter: NarratorFilter("emilywoozeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match tag", filter: TagFilter("Hugo Awards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring case", filter: TagFilter("Hugo awards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring whitespace", filter: TagFilter("HugoAwards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring case and whitespace", filter: TagFilter("hugoawards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match tag", filter: TagFilter("Nebula Awards"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "use not filter with author", filter: NotFilter(AuthorFilter("Max Gladstone")), book: testbooks.Audiobooks[1], expectedMatch: true},
		{name: "use not filter with author fails", filter: NotFilter(AuthorFilter("Max Gladstone")), book: testbooks.Audiobooks[0], expectedMatch: false},
		{name: "combine multiple filters", filter: AndFilter(AuthorFilter("Max Gladstone"), GenreFilter(audiobooks.SciFi)), book: testbooks.Audiobooks[0], expectedMatch: true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedMatch, testCase.filter(&testCase.book))
		})
	}
}

func TestAuthorFilter_NilAudiobook(t *testing.T) {
	filter := AuthorFilter("Test Author")
	result := filter(nil)
	assert.False(t, result)
}

func TestGenreFilter_NilAudiobook(t *testing.T) {
	filter := GenreFilter(audiobooks.SciFi)
	result := filter(nil)
	assert.False(t, result)
}

func TestNarratorFilter_NilAudiobook(t *testing.T) {
	filter := NarratorFilter("Test Narrator")
	result := filter(nil)
	assert.False(t, result)
}

func TestTagFilter_NilAudiobook(t *testing.T) {
	filter := TagFilter("Test Tag")
	result := filter(nil)
	assert.False(t, result)
}

func TestAndFilter_EmptyFilters(t *testing.T) {
	filter := AndFilter()
	result := filter(&testbooks.Audiobooks[0])
	assert.True(t, result) // Empty filters should match everything
}

func TestAndFilter_SingleFilter(t *testing.T) {
	filter := AndFilter(AuthorFilter("Max Gladstone"))
	result := filter(&testbooks.Audiobooks[0])
	assert.True(t, result)
}

func TestAndFilter_MultipleFiltersAllMatch(t *testing.T) {
	filter := AndFilter(
		AuthorFilter("Max Gladstone"),
		GenreFilter(audiobooks.SciFi),
		TagFilter("Hugo Awards"),
	)
	result := filter(&testbooks.Audiobooks[0])
	assert.True(t, result)
}

func TestAndFilter_MultipleFiltersOneDoesNotMatch(t *testing.T) {
	filter := AndFilter(
		AuthorFilter("Max Gladstone"),
		GenreFilter(audiobooks.Romance), // This won't match
		TagFilter("Hugo Awards"),
	)
	result := filter(&testbooks.Audiobooks[0])
	assert.False(t, result)
}

func TestContains_Generic(t *testing.T) {
	tests := []struct {
		name     string
		slice    []audiobooks.Genre
		item     audiobooks.Genre
		expected bool
	}{
		{"item found", []audiobooks.Genre{audiobooks.SciFi, audiobooks.Fantasy}, audiobooks.SciFi, true},
		{"item not found", []audiobooks.Genre{audiobooks.SciFi, audiobooks.Fantasy}, audiobooks.Romance, false},
		{"empty slice", []audiobooks.Genre{}, audiobooks.SciFi, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsIgnoringCaseAndWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"exact match", []string{"Max Gladstone", "Amal El-Mohtar"}, "Max Gladstone", true},
		{"case insensitive match", []string{"Max Gladstone", "Amal El-Mohtar"}, "max gladstone", true},
		{"whitespace insensitive match", []string{"Max Gladstone", "Amal El-Mohtar"}, "MaxGladstone", true},
		{"case and whitespace insensitive match", []string{"Max Gladstone", "Amal El-Mohtar"}, "maxgladstone", true},
		{"no match", []string{"Max Gladstone", "Amal El-Mohtar"}, "Unknown Author", false},
		{"empty slice", []string{}, "Max Gladstone", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsIgnoringCaseAndWhitespace(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormaliseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "hello", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"mixed case", "HeLLo", "hello"},
		{"with spaces", "Hello World", "helloworld"},
		{"with tabs", "Hello\tWorld", "helloworld"},
		{"with newlines", "Hello\nWorld", "helloworld"},
		{"multiple whitespace", "Hello   World", "helloworld"},
		{"special characters", "Hello-World_123!", "hello-world_123!"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normaliseString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
