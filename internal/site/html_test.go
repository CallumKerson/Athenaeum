package site

import (
	"html"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var linkPattern = regexp.MustCompile(`(?:href|src)="([^"]+)"`)

func dated(title string, release *toml.LocalDate) audiobooks.Audiobook {
	audiobook := book(title, nil, []string{"A"}, nil, nil)
	audiobook.ReleaseDate = release
	return audiobook
}

// Nothing resolves relative links for a static site, so a page that links to a
// path the build did not write is a dead end with no server to notice.
func TestPagesLinkOnlyToFilesThatExist(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	var checked int
	err = filepath.WalkDir(root, func(fullPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(fullPath) != ".html" {
			return walkErr
		}
		contents, readErr := os.ReadFile(fullPath)
		require.NoError(t, readErr)

		for _, match := range linkPattern.FindAllStringSubmatch(string(contents), -1) {
			link, parseErr := url.Parse(html.UnescapeString(match[1]))
			require.NoError(t, parseErr, match[1])
			if link.IsAbs() {
				continue
			}
			target := filepath.Join(filepath.Dir(fullPath), filepath.FromSlash(link.Path))
			if strings.HasSuffix(link.Path, "/") {
				target = filepath.Join(target, indexName)
			}
			assert.FileExists(t, target, "%s links to %s", fullPath, link.Path)
			checked++
		}
		return nil
	})
	require.NoError(t, err)
	assert.Positive(t, checked)
}

// Every feed has a page, so a name can be browsed to as well as subscribed to.
func TestEveryFeedHasAPage(t *testing.T) {
	root := t.TempDir()
	content := testPlan(t)
	_, err := Build(root, content, testRenderer(t), Options{Sweep: true}, testLogger())
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(root, indexName))
	assert.FileExists(t, filepath.Join(root, booksDir, indexName))
	for _, section := range sections(content) {
		assert.FileExists(t, filepath.Join(root, section.Dir, indexName))
		for index := range section.Names {
			name := &section.Names[index]
			require.NotEmpty(t, name.Slug, name.Title)
			assert.FileExists(t, filepath.Join(root, section.Dir, name.Slug, indexName))
		}
	}
	// A genre with nothing filed under it keeps its feed, so it keeps its page.
	assert.FileExists(t, filepath.Join(root, genresDir, "horror", indexName))
}

// Every section the plan publishes feeds for must also be shown, in the same
// order, with a heading and blurb to introduce it.
func TestEverySectionIsShown(t *testing.T) {
	content := Plan(nil, Exclusions{})
	shown := shownSections(content)
	planned := sections(content)

	require.Len(t, shown, len(planned))
	for index, section := range shown {
		assert.Same(t, planned[index], section.Section)
		assert.NotEmpty(t, section.Title, section.Dir)
		assert.NotEmpty(t, section.Blurb, section.Dir)
	}
}

// The genre exclusion is a rule about the main feed, not about the library, so
// the browsable list still holds everything.
func TestAllBooksPageListsBooksExcludedFromTheMainFeed(t *testing.T) {
	root := t.TempDir()
	content := Plan([]audiobooks.Audiobook{
		book("Clean", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
		book("Spicy", []audiobooks.Genre{audiobooks.Erotica}, []string{"B"}, nil, nil),
	}, Exclusions{MainFeed: []audiobooks.Genre{audiobooks.Erotica}})

	_, err := Build(root, content, testRenderer(t), Options{Sweep: true}, testLogger())
	require.NoError(t, err)

	page, err := os.ReadFile(filepath.Join(root, booksDir, indexName))
	require.NoError(t, err)
	assert.Contains(t, string(page), "Clean")
	assert.Contains(t, string(page), "Spicy")
}

// A genre excluded from the site leaves no trace in the pages: not its books,
// not its own page, and not the pages of names that have no other books.
func TestPagesLeaveOutGenresExcludedFromTheSite(t *testing.T) {
	root := t.TempDir()
	spicy := book("Spicy", []audiobooks.Genre{audiobooks.Erotica}, []string{"B"}, []string{"N"}, []string{"T"})
	spicy.Series = &audiobooks.Series{Title: "Steamy"}
	content := Plan([]audiobooks.Audiobook{
		book("Clean", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, []string{"N"}, nil),
		spicy,
	}, Exclusions{Site: []audiobooks.Genre{audiobooks.Erotica}})

	_, err := Build(root, content, testRenderer(t), Options{Sweep: true}, testLogger())
	require.NoError(t, err)

	err = filepath.WalkDir(root, func(fullPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(fullPath) != ".html" {
			return walkErr
		}
		page, readErr := os.ReadFile(fullPath)
		require.NoError(t, readErr)
		assert.NotContains(t, string(page), "Spicy", fullPath)
		assert.NotContains(t, string(page), "Erotica", fullPath)
		assert.NotContains(t, string(page), "Steamy", fullPath)
		return nil
	})
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, booksDir, "spicy"))
	assert.NoDirExists(t, filepath.Join(root, seriesDir, "steamy"))
	assert.NoDirExists(t, filepath.Join(root, genresDir, "erotica"))
	assert.NoDirExists(t, filepath.Join(root, authorsDir, "b"))
	assert.NoDirExists(t, filepath.Join(root, tagsDir, "t"))
	assert.FileExists(t, filepath.Join(root, narratorsDir, "n", indexName))
	assert.FileExists(t, filepath.Join(root, podcastDir, genresDir, "Erotica", feedName))
	assert.FileExists(t, filepath.Join(root, podcastDir, authorsDir, "B", feedName))
}

// Pages lead with the newest book; a book with no release date has nothing to
// sort on and goes last.
func TestBookListsAreNewestFirst(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		dated("Undated", nil),
		dated("Oldest", &toml.LocalDate{Year: 1968, Month: 11, Day: 1}),
		dated("Newest", &toml.LocalDate{Year: 2023, Month: 4, Day: 6}),
		dated("Middle", &toml.LocalDate{Year: 2019, Month: 7, Day: 16}),
	}, Exclusions{})

	titles := func(views []bookView) []string {
		found := make([]string, 0, len(views))
		for _, view := range views {
			found = append(found, view.Title)
		}
		return found
	}

	assert.Equal(
		t,
		[]string{"Newest", "Middle", "Oldest", "Undated"},
		titles(allBooksView(newLibrary(content, nil, nil)).Books),
	)
	// Every page that lists books, not just the whole-library one.
	require.Len(t, content.Authors.Names, 1)
	assert.Equal(
		t,
		[]string{"Newest", "Middle", "Oldest", "Undated"},
		titles(
			nameView(
				newLibrary(content, nil, nil),
				shownSection{Section: &content.Authors},
				&content.Authors.Names[0],
			).Books,
		),
	)
}

// The feeds are rendered from the same slices the pages list, and their item
// order is the one subscribers already have.
func TestBuildingPagesDoesNotReorderFeeds(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		dated("Oldest", &toml.LocalDate{Year: 1968, Month: 11, Day: 1}),
		dated("Newest", &toml.LocalDate{Year: 2023, Month: 4, Day: 6}),
	}, Exclusions{})

	before, err := testRenderer(t).Render(content.Main.Books, "Audiobooks", "")
	require.NoError(t, err)

	_, err = buildHTML(newLibrary(content, nil, nil))
	require.NoError(t, err)

	after, err := testRenderer(t).Render(content.Main.Books, "Audiobooks", "")
	require.NoError(t, err)
	assert.Equal(t, string(before), string(after))
}

// The date shown is the day the feed gives the book as its pubDate.
func TestBookDateMatchesTheFeed(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		dated("Dated", &toml.LocalDate{Year: 1968, Month: 11, Day: 1}),
		dated("Undated", nil),
	}, Exclusions{})

	books := allBooksView(newLibrary(content, nil, nil)).Books
	require.Len(t, books, 2)

	assert.Equal(t, "1 Nov 1968", books[0].Date)
	assert.Contains(t, books[0].Search, "1968", "the year should be searchable")
	assert.Empty(t, books[1].Date)
}

// An omnibus stands in for several books in its series, and the page says so.
func TestBookSeriesLine(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		expected string
	}{
		{"Single book", "1", "Earthsea book 1"},
		{"Omnibus", "1-3", "Earthsea books 1-3"},
		{"No sequence", "", "Earthsea"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var sequence *audiobooks.Sequence
			if testCase.sequence != "" {
				parsed, err := audiobooks.ParseSequence(testCase.sequence)
				require.NoError(t, err)
				sequence = &parsed
			}

			audiobook := book("Earthsea", nil, []string{"Ursula K. Le Guin"}, nil, nil)
			audiobook.Series = &audiobooks.Series{Title: "Earthsea", Sequence: sequence}

			content := Plan([]audiobooks.Audiobook{audiobook}, Exclusions{})
			books := allBooksView(newLibrary(content, nil, nil)).Books
			require.Len(t, books, 1)

			assert.Equal(t, testCase.expected, books[0].Series)
			assert.Contains(t, books[0].Search, strings.ToLower(testCase.expected))
		})
	}
}

// A book is found by anything a reader might remember about it.
func TestBookSearchTextCoversEveryField(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{{
		Title:     "A Wizard of Earthsea",
		Authors:   []string{"Ursula K. Le Guin"},
		Narrators: []string{"Kobna Holdbrook-Smith"},
		Genres:    []audiobooks.Genre{audiobooks.Fantasy},
		Tags:      []string{"Hugo Awards"},
	}}, Exclusions{})

	books := allBooksView(newLibrary(content, nil, nil)).Books
	require.Len(t, books, 1)

	for _, term := range []string{"earthsea", "le guin", "holdbrook", "fantasy", "hugo"} {
		assert.Contains(t, books[0].Search, term)
	}
	// Absent fields must not leave gaps that break substring matching.
	assert.NotContains(t, books[0].Search, "  ")
}

// A book links each name it lists to that name's page, whichever spelling it
// uses, and leaves a name with no page as plain text.
func TestBookLinksToTheNamesItLists(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		book("First", nil, []string{"Ann Leckie"}, nil, nil),
		book("Second", []audiobooks.Genre{audiobooks.SciFi}, []string{"Ann-Leckie", "ann leckie"}, []string{"  "},
			[]string{"Hugo Awards"}),
	}, Exclusions{})

	books := allBooksView(newLibrary(content, nil, nil)).Books
	require.Len(t, books, 2)
	second := books[1]
	require.Equal(t, "Second", second.Title)

	assert.Equal(t, []linkView{
		{Name: "Ann-Leckie", Href: "../authors/ann-leckie/"},
		{Name: "ann leckie", Href: "../authors/ann-leckie-2/"},
	}, second.Authors)
	assert.Equal(t, []linkView{{Name: "  "}}, second.Narrators)
	assert.Equal(t, []linkView{{Name: "Science Fiction", Href: "../genre/science-fiction/"}}, second.Genres)
	assert.Equal(t, []linkView{{Name: "Hugo Awards", Href: "../tags/hugo-awards/"}}, second.Tags)
}

func TestSlugify(t *testing.T) {
	for name, expected := range map[string]string{
		"Ursula K. Le Guin": "ursula-k-le-guin",
		"V.E. Schwab":       "v-e-schwab",
		"Children's":        "children-s",
		"LGBT+":             "lgbt",
		"  Nettle & Bone  ": "nettle-bone",
		"???":               "unnamed",
	} {
		assert.Equal(t, expected, slugify(name), name)
	}
}

func TestFormatDuration(t *testing.T) {
	for duration, expected := range map[time.Duration]string{
		0:                                  "",
		30 * time.Second:                   "",
		9 * time.Minute:                    "9m",
		time.Hour + 5*time.Minute:          "1h 05m",
		16*time.Hour + 8*time.Minute:       "16h 08m",
		16*time.Hour + 8*time.Minute + 999: "16h 08m",
	} {
		assert.Equal(t, expected, formatDuration(duration), duration.String())
	}
}
