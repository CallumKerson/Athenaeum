package feed

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
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

// Filenames are data, not URL syntax. Each of these characters is meaningful in
// a URL and would otherwise truncate or corrupt the path — and the enclosure URL
// is also the GUID, so getting one wrong makes a book permanently unreachable.
func TestMediaURLEscapesCharactersWithURLMeaning(t *testing.T) {
	renderer := &Renderer{Host: "https://example.com", MediaPath: "/media"}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "question mark would start the query string",
			path:     "/Randall Munroe/What If?2/What If?2.m4b",
			expected: "https://example.com/media/Randall%20Munroe/What%20If%3F2/What%20If%3F2.m4b",
		},
		{
			name:     "hash would start the fragment and drop the rest",
			path:     "/Artist/Album #2/Track #2.m4b",
			expected: "https://example.com/media/Artist/Album%20%232/Track%20%232.m4b",
		},
		{
			name:     "percent would be read as an escape sequence",
			path:     "/Author/100% Cotton/100% Cotton.m4b",
			expected: "https://example.com/media/Author/100%25%20Cotton/100%25%20Cotton.m4b",
		},
		{
			name:     "a literal %20 must not decode to a space",
			path:     "/Author/Odd%20Name/Odd%20Name.m4b",
			expected: "https://example.com/media/Author/Odd%2520Name/Odd%2520Name.m4b",
		},
		{
			name: "ampersands and apostrophes keep their existing encoding",
			path: "/M. A. Carrick/Rook & Rose/2 The Liar's Knot/The Liar's Knot.m4b",
			expected: "https://example.com/media/M.%20A.%20Carrick/Rook%20&%20Rose/" +
				"2%20The%20Liar%27s%20Knot/The%20Liar%27s%20Knot.m4b",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := renderer.mediaURL(testCase.path)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)

			// The escaped URL must name the file we started from.
			parsed, err := url.Parse(actual)
			require.NoError(t, err)
			assert.Equal(t, "/media"+testCase.path, parsed.Path)
			assert.Empty(t, parsed.RawQuery)
			assert.Empty(t, parsed.Fragment)
		})
	}
}

// A host carrying a path prefix must keep it.
func TestMediaURLKeepsHostPathPrefix(t *testing.T) {
	renderer := &Renderer{Host: "http://www.example-podcast.com/audiobooks/", MediaPath: "/media/"}

	actual, err := renderer.mediaURL("/Author/Book/Book.m4b")

	require.NoError(t, err)
	assert.Equal(t, "http://www.example-podcast.com/audiobooks/media/Author/Book/Book.m4b", actual)
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

// The channel-level iTunes metadata is optional, so the golden feeds above are
// rendered without any of it. These cover the branches that set it.
func TestRenderSetsChannelMetadata(t *testing.T) {
	renderer := &Renderer{
		Host:      "http://www.example-podcast.com",
		MediaPath: "/media/",
		Author:    "Athenaeum",
		Email:     "books@example.com",
		Explicit:  true,
		ImageLink: "http://www.example-podcast.com/static/itunes_image.jpg",
		Language:  "EN",
	}

	rendered, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")

	require.NoError(t, err)
	feed := string(rendered)
	assert.Contains(t, feed, "<itunes:author>Athenaeum</itunes:author>")
	assert.Contains(t, feed, "<itunes:name>Athenaeum</itunes:name>")
	assert.Contains(t, feed, "<itunes:email>books@example.com</itunes:email>")
	assert.Contains(t, feed, "<itunes:explicit>yes</itunes:explicit>")
	assert.Contains(t, feed, `<itunes:image href="http://www.example-podcast.com/static/itunes_image.jpg">`)
	assert.Contains(t, feed, "<language>EN</language>")
}

// A Host that url.Parse rejects has to fail the build rather than produce a feed
// full of broken enclosure URLs, which are also the GUIDs.
func TestRenderRejectsUnparseableHost(t *testing.T) {
	renderer := &Renderer{Host: "http://[::1", MediaPath: "/media/"}

	_, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")

	require.Error(t, err)
}

// Podcast clients disagree about negative timestamps, so a release date before
// 1970 is clamped up to the epoch — plus the eight hour offset every pubDate gets.
func TestPubDateClampsPreUnixEpochDates(t *testing.T) {
	preEpoch := toml.LocalDate{Year: 1968, Month: 11, Day: 1}
	book := audiobooks.Audiobook{Title: "A Wizard of Earthsea", ReleaseDate: &preEpoch}

	clamping := &Renderer{HandlePreUnixEpoch: true}
	assert.Equal(t, unixEpoch.Add(8*time.Hour), clamping.pubDate(&book))

	passingThrough := &Renderer{}
	assert.Equal(t, preEpoch.AsTime(time.UTC).Add(8*time.Hour), passingThrough.pubDate(&book))
}

func TestSummaryHTMLRepeatsTitleForSubtitledBook(t *testing.T) {
	summary := summaryHTML(&audiobooks.Audiobook{
		Title:    "What If?",
		Subtitle: "Serious Scientific Answers to Absurd Hypothetical Questions",
		Authors:  []string{"Randall Munroe"},
	})

	assert.Contains(t, summary, "<h1>What If?</h1>")
	// Reproduces the original: the title, not the subtitle, is repeated.
	assert.Contains(t, summary, "<h4>What If?</h4>")
	assert.NotContains(t, summary, "Absurd Hypothetical")
}

func TestSummaryHTMLDescriptionFormats(t *testing.T) {
	tests := []struct {
		name   string
		format description.Format
		text   string
		want   string
	}{
		{
			name:   "markdown becomes HTML",
			format: description.Markdown,
			text:   "A *terrible* shadow.",
			want:   "<em>terrible</em>",
		},
		{
			name:   "HTML passes through",
			format: description.HTML,
			text:   "<p>A <em>terrible</em> shadow.</p>",
			want:   "<p>A <em>terrible</em> shadow.</p>",
		},
		{
			name:   "plain text becomes paragraphs",
			format: description.Plain,
			text:   "First line.\nSecond line.",
			want:   "<p>First line.</p><p>Second line.</p>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			summary := summaryHTML(&audiobooks.Audiobook{
				Title:       "A Wizard of Earthsea",
				Authors:     []string{"Ursula K. Le Guin"},
				Description: &description.Description{Text: test.text, Format: test.format},
			})

			assert.Contains(t, summary, test.want)
		})
	}
}

func TestSummaryHTMLIncludesSeriesAndNarrator(t *testing.T) {
	summary := summaryHTML(&audiobooks.Audiobook{
		Title:     "A Wizard of Earthsea",
		Authors:   []string{"Ursula K. Le Guin"},
		Narrators: []string{"Kobna Holdbrook-Smith"},
		Series:    &audiobooks.Series{Title: "Earthsea", Sequence: decimal.NewFromInt(1)},
	})

	assert.Contains(t, summary, "<h2>By Ursula K. Le Guin</h2>")
	assert.Contains(t, summary, "<h4>Earthsea Book 1</h4>")
	assert.Contains(t, summary, "<h4>Narrated by Kobna Holdbrook-Smith</h4>")
}
