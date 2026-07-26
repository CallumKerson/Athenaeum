package site

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestNameAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multi-word name gets all three forms",
			input:    "Ursula K. Le Guin",
			expected: []string{"Ursula K. Le Guin", "ursula k. le guin", "ursulak.leguin"},
		},
		{
			name:     "single lowercase word needs only one form",
			input:    "sleepy",
			expected: []string{"sleepy"},
		},
		{
			name:     "single capitalised word has no whitespace to strip",
			input:    "Fantasy",
			expected: []string{"Fantasy", "fantasy"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.ElementsMatch(t, testCase.expected, nameAliases(testCase.input))
		})
	}
}

// Every genre URL the old server accepted must still resolve, including the
// synonyms ParseGenre understands.
func TestGenreAliasesCoverEveryAcceptedSpelling(t *testing.T) {
	assert.ElementsMatch(t,
		[]string{"Science Fiction", "science fiction", "sciencefiction", "sci-fi", "scifi"},
		genreAliases(audiobooks.SciFi))

	for _, genre := range audiobooks.AllGenres() {
		for _, alias := range genreAliases(genre) {
			parsed, err := audiobooks.ParseGenre(alias)
			assert.NoError(t, err, alias)
			assert.Equal(t, genre, parsed, alias)
		}
	}
}

// Names come from user-editable metadata, so they must never be able to direct
// a write outside the output directory.
func TestSafeSegment(t *testing.T) {
	for _, name := range []string{"Ursula K. Le Guin", "sci-fi", "Nettle & Bone"} {
		assert.True(t, safeSegment(name), name)
	}
	for _, name := range []string{"", ".", "..", "../../etc", "a/b"} {
		assert.False(t, safeSegment(name), name)
	}
}

func TestFeedPathsDropUnsafeNames(t *testing.T) {
	canonical, paths := feedPaths("authors", "ok", []string{"..", "ok", "a/b"})

	assert.Equal(t, []string{"podcast/authors/ok/feed.rss"}, paths)
	assert.Equal(t, "podcast/authors/ok/feed.rss", canonical)
}

// A name that cannot be a directory at all has no feed, so nothing to link to.
func TestFeedPathsWithNoUsableSpelling(t *testing.T) {
	canonical, paths := feedPaths("authors", "a/b", []string{"a/b"})

	assert.Empty(t, paths)
	assert.Empty(t, canonical)
}
