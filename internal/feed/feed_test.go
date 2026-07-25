package feed

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// The golden files here are copies of the ones the server's podcast service was
// tested against. Matching them byte for byte is the parity guarantee: an
// existing subscriber must not be able to tell the generator from the server.
func TestRenderMatchesServerOutput(t *testing.T) {
	renderer := &Renderer{
		Host:      "http://www.example-podcast.com/audiobooks/",
		MediaPath: "/media/",
	}

	tests := []struct {
		name        string
		books       []audiobooks.Audiobook
		title       string
		description string
		golden      string
	}{
		{
			name:        "all audiobooks",
			books:       testbooks.Audiobooks,
			title:       "Audiobooks",
			description: "Like movies in your mind!",
			golden:      "feed.rss",
		},
		{
			name:        "genre",
			books:       testbooks.AudiobooksFilteredBy(testbooks.GenreFilter(audiobooks.SciFi)),
			title:       audiobooks.SciFi.String(),
			description: "Science Fiction Audiobooks",
			golden:      "scifi_feed.rss",
		},
		{
			name:        "author",
			books:       testbooks.AudiobooksFilteredBy(testbooks.AuthorFilter("Amal El-Mohtar")),
			title:       "Amal El-Mohtar",
			description: "Audiobooks by Amal El-Mohtar",
			golden:      "el_mohtar_feed.rss",
		},
		{
			name:        "narrator",
			books:       testbooks.AudiobooksFilteredBy(testbooks.NarratorFilter("Kobna Holdbrook-Smith")),
			title:       "Kobna Holdbrook-Smith",
			description: "Audiobooks Narrated by Kobna Holdbrook-Smith",
			golden:      "holdbrook_smith_feed.rss",
		},
		{
			name:        "tag",
			books:       testbooks.AudiobooksFilteredBy(testbooks.TagFilter("Hugo Awards")),
			title:       "Hugo Awards",
			description: "Hugo Awards Audiobooks",
			golden:      "hugo_awards_feed.rss",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			expected, err := os.ReadFile(filepath.Join("testdata", testCase.golden))
			require.NoError(t, err)

			rendered, err := renderer.Render(testCase.books, testCase.title, testCase.description)
			require.NoError(t, err)

			assert.Equal(t, strings.TrimSpace(string(expected)), string(rendered))
		})
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	renderer := &Renderer{Host: "https://example.com", MediaPath: "/media"}

	first, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")
	require.NoError(t, err)
	second, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")
	require.NoError(t, err)

	assert.Equal(t, string(first), string(second))
}

// A book with no release date must not panic the renderer, even though every
// book in the real library has one.
func TestRenderWithoutReleaseDate(t *testing.T) {
	renderer := &Renderer{Host: "https://example.com", MediaPath: "/media"}

	rendered, err := renderer.Render([]audiobooks.Audiobook{{
		Title:    "Undated",
		Authors:  []string{"Nobody"},
		Path:     "/Nobody/Undated/Undated.m4b",
		MIMEType: "audio/mp4a-latm",
	}}, "Audiobooks", "Like movies in your mind!")

	require.NoError(t, err)
	assert.Contains(t, string(rendered), "<title>Undated</title>")
}
