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
	// The server was tested without channel artwork, so Artwork is left off.
	renderer := newTestRenderer(t, "http://www.example-podcast.com/audiobooks/", Channel{})

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
	renderer := newTestRenderer(t, "https://example.com", Channel{})

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
			actual := renderer.mediaURL(testCase.path)
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

// The GUID is the enclosure URL, so the published URLs must not depend on how
// the host was written: a trailing slash added to the config must not make
// every subscriber re-download the library. A path prefix must be kept, and
// the artwork, channel link and enclosures must all share the same prefix.
func TestPublishedURLsAreStableAcrossHostSpellings(t *testing.T) {
	tests := []struct {
		name      string
		hosts     []string
		published string
	}{
		{
			name: "bare host",
			hosts: []string{
				"https://books.example.com",
				"https://books.example.com/",
				"https://books.example.com//",
				// The slash is only trimmed if it is trimmed after parsing: in
				// the raw string it is not the last character.
				"https://books.example.com/#",
				"HTTPS://books.example.com/",
			},
			published: "https://books.example.com",
		},
		{
			name: "host with a path prefix",
			hosts: []string{
				"http://books.example.com:8080/audiobooks",
				"http://books.example.com:8080/audiobooks/",
				"http://books.example.com:8080/audiobooks/#",
				"HTTP://books.example.com:8080/audiobooks/",
			},
			published: "http://books.example.com:8080/audiobooks",
		},
	}

	book := audiobooks.Audiobook{
		Title:     "Earthsea",
		Authors:   []string{"Ursula K. Le Guin"},
		Path:      "/Ursula K Le Guin/Earthsea/Earthsea.m4b",
		ImagePath: "/Ursula K Le Guin/Earthsea/Earthsea.jpg",
		MIMEType:  "audio/mp4a-latm",
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			for _, host := range testCase.hosts {
				renderer := newTestRenderer(t, host, Channel{Artwork: true})
				rendered, err := renderer.Render([]audiobooks.Audiobook{book}, "Audiobooks", "")
				require.NoError(t, err, host)

				feed := string(rendered)
				media := testCase.published + "/media/Ursula%20K%20Le%20Guin/Earthsea/"
				enclosure := media + "Earthsea.m4b"
				assert.Contains(t, feed, "<guid>"+enclosure+"</guid>", host)
				assert.Contains(t, feed, `<enclosure url="`+enclosure+`"`, host)
				assert.Contains(t, feed, `<itunes:image href="`+media+`Earthsea.jpg">`, host)
				assert.Contains(t, feed, "<link>"+testCase.published+"</link>", host)
				assert.Contains(t, feed, `<itunes:image href="`+testCase.published+`/static/itunes_image.jpg">`, host)
				assert.Equal(t, testCase.published, renderer.Host(), host)
				assert.NotContains(t, feed, "//media", host)
				assert.NotContains(t, feed, "//static", host)
			}
		})
	}
}

// Each of these would publish enclosure URLs, and so GUIDs, that no client can
// fetch, or that carry a query, fragment or credentials into every GUID, so New
// refuses them rather than letting the build write the feeds.
func TestNewRejectsInvalidHosts(t *testing.T) {
	for _, host := range []string{
		"", "books.example.com", "/books", "ftp://books.example.com", "https://", "://x", "http://[::1",
		"https://books.example.com?q=1", "https://books.example.com/?q=1", "https://books.example.com?",
		"https://books.example.com#f", "https://books.example.com/#f",
		"https://u:p@books.example.com", "https://u@books.example.com",
	} {
		renderer, err := New(host, Channel{})

		require.ErrorIs(t, err, ErrInvalidHost, host)
		assert.Contains(t, err.Error(), host)
		assert.Nil(t, renderer)
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	renderer := newTestRenderer(t, "https://example.com", Channel{})

	first, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")
	require.NoError(t, err)
	second, err := renderer.Render(testbooks.Audiobooks, "Audiobooks", "Like movies in your mind!")
	require.NoError(t, err)

	assert.Equal(t, string(first), string(second))
}

// An omnibus stands in for several books in its series, and its heading says so.
// A series with no sequence is named without a place in it.
func TestRenderSeriesHeading(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		expected string
	}{
		{"Single book", "1", "<h4>Earthsea Book 1</h4>"},
		{"Prequel", "0", "<h4>Earthsea Book 0</h4>"},
		{"Omnibus", "1-3", "<h4>Earthsea Books 1-3</h4>"},
		{"No sequence", "", "<h4>Earthsea</h4>"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var sequence *audiobooks.Sequence
			if testCase.sequence != "" {
				parsed, err := audiobooks.ParseSequence(testCase.sequence)
				require.NoError(t, err)
				sequence = &parsed
			}

			renderer := newTestRenderer(t, "https://example.com", Channel{})
			rendered, err := renderer.Render([]audiobooks.Audiobook{{
				Title:    "Earthsea",
				Authors:  []string{"Ursula K. Le Guin"},
				Path:     "/Ursula K Le Guin/Earthsea/Earthsea.m4b",
				MIMEType: "audio/mp4a-latm",
				Series:   &audiobooks.Series{Title: "Earthsea", Sequence: sequence},
			}}, "Audiobooks", "Like movies in your mind!")

			require.NoError(t, err)
			assert.Contains(t, string(rendered), testCase.expected)
		})
	}
}

// A book with no release date must not panic the renderer, even though every
// book in the real library has one.
func TestRenderWithoutReleaseDate(t *testing.T) {
	renderer := newTestRenderer(t, "https://example.com", Channel{})

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
	renderer := newTestRenderer(t, "http://www.example-podcast.com", Channel{
		Author:   "Athenaeum",
		Email:    "books@example.com",
		Explicit: true,
		Language: "EN",
		Artwork:  true,
	})

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

// Podcast clients disagree about negative timestamps, so a release date before
// 1970 is clamped up to the epoch — plus the eight hour offset every pubDate gets.
func TestPubDateClampsPreUnixEpochDates(t *testing.T) {
	preEpoch := toml.LocalDate{Year: 1968, Month: 11, Day: 1}
	book := audiobooks.Audiobook{Title: "A Wizard of Earthsea", ReleaseDate: &preEpoch}

	clamping := newTestRenderer(t, "https://example.com", Channel{HandlePreUnixEpoch: true})
	assert.Equal(t, unixEpoch.Add(8*time.Hour), clamping.pubDate(&book))

	passingThrough := newTestRenderer(t, "https://example.com", Channel{})
	assert.Equal(t, preEpoch.AsTime(time.UTC).Add(8*time.Hour), passingThrough.pubDate(&book))
}

func newTestRenderer(t *testing.T, host string, channel Channel) *Renderer {
	t.Helper()
	renderer, err := New(host, channel)
	require.NoError(t, err)
	return renderer
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
		Series:    &audiobooks.Series{Title: "Earthsea", Sequence: &audiobooks.Sequence{First: decimal.NewFromInt(1)}},
	})

	assert.Contains(t, summary, "<h2>By Ursula K. Le Guin</h2>")
	assert.Contains(t, summary, "<h4>Earthsea Book 1</h4>")
	assert.Contains(t, summary, "<h4>Narrated by Kobna Holdbrook-Smith</h4>")
}
