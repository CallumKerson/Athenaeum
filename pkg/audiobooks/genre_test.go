package audiobooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenre_String(t *testing.T) {
	tests := []struct {
		genre    Genre
		expected string
	}{
		{Literary, "Literary"},
		{Mystery, "Mystery"},
		{Romance, "Romance"},
		{Comedy, "Comedy"},
		{Childrens, "Children's"},
		{YoungAdult, "Young Adult"},
		{SciFi, "Science Fiction"},
		{Fantasy, "Fantasy"},
		{NonFiction, "Non-fiction"},
		{Biography, "Biography"},
		{Historical, "Historical Fiction"},
		{Thriller, "Thriller"},
		{Horror, "Horror"},
		{LGBT, "LGBT+"},
		{Erotica, "Erotica"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.genre.String())
		})
	}
}

func TestGenre_String_Undefined(t *testing.T) {
	var genre Genre
	assert.Equal(t, "", genre.String())
}

func TestParseGenre_ValidInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected Genre
	}{
		// Exact matches
		{"literary", Literary},
		{"mystery", Mystery},
		{"romance", Romance},
		{"comedy", Comedy},
		{"children's", Childrens},
		{"children", Childrens},
		{"childrens", Childrens},
		{"young adult", YoungAdult},
		{"youngadult", YoungAdult},
		{"ya", YoungAdult},
		{"science fiction", SciFi},
		{"sciencefiction", SciFi},
		{"sci-fi", SciFi},
		{"scifi", SciFi},
		{"fantasy", Fantasy},
		{"non-fiction", NonFiction},
		{"nonfiction", NonFiction},
		{"biography", Biography},
		{"historical", Historical},
		{"historical fiction", Historical},
		{"historicalfiction", Historical},
		{"thriller", Thriller},
		{"horror", Horror},
		{"lgbt", LGBT},
		{"lgbt+", LGBT},
		{"erotica", Erotica},

		// Case variations
		{"LITERARY", Literary},
		{"Mystery", Mystery},
		{"ROMANCE", Romance},
		{"SciFi", SciFi},
		{"FANTASY", Fantasy},
		{"LGBT+", LGBT},

		// With whitespace
		{" literary ", Literary},
		{"\tmystery\n", Mystery},
		{" romance ", Romance},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseGenre(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseGenre_InvalidInputs(t *testing.T) {
	tests := []string{
		"",
		"unknown",
		"invalid",
		"not-a-genre",
		"123",
		"fiction-mystery", // hyphenated combination
		"scifi-fantasy",   // hyphenated combination
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseGenre(input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not a valid genre")
		})
	}
}

func TestGenre_MarshalText(t *testing.T) {
	tests := []struct {
		genre    Genre
		expected string
	}{
		{Literary, "Literary"},
		{SciFi, "Science Fiction"},
		{LGBT, "LGBT+"},
		{YoungAdult, "Young Adult"},
		{NonFiction, "Non-fiction"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			data, err := tt.genre.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestGenre_UnmarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected Genre
	}{
		{"literary", Literary},
		{"scifi", SciFi},
		{"lgbt+", LGBT},
		{"youngadult", YoungAdult},
		{"nonfiction", NonFiction},
		{"sci-fi", SciFi},           // alias
		{"young adult", YoungAdult}, // with space
		{"MYSTERY", Mystery},        // uppercase
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var genre Genre
			err := genre.UnmarshalText([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, genre)
		})
	}
}

func TestGenre_UnmarshalText_Error(t *testing.T) {
	var genre Genre
	err := genre.UnmarshalText([]byte("invalid-genre"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid genre")
}

func TestGenre_RoundTrip(t *testing.T) {
	// Test that marshalling and unmarshalling preserves the genre
	allGenres := []Genre{
		Literary, Mystery, Romance, Comedy, Childrens, YoungAdult,
		SciFi, Fantasy, NonFiction, Biography, Historical,
		Thriller, Horror, LGBT, Erotica,
	}

	for _, original := range allGenres {
		t.Run(original.String(), func(t *testing.T) {
			// Marshal to text
			data, err := original.MarshalText()
			require.NoError(t, err)

			// Unmarshal back
			var unmarshaled Genre
			err = unmarshaled.UnmarshalText(data)
			require.NoError(t, err)

			// Should be equal
			assert.Equal(t, original, unmarshaled)
		})
	}
}
