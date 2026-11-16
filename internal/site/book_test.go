package site

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

func inSeries(title, series string, position int64) audiobooks.Audiobook {
	audiobook := book(title, nil, []string{"Ursula K. Le Guin"}, nil, nil)
	audiobook.Series = &audiobooks.Series{Title: series}
	if position > 0 {
		audiobook.Series.Sequence = &audiobooks.Sequence{First: decimal.NewFromInt(position)}
	}
	return audiobook
}

func TestPlanGivesEveryBookItsOwnSlug(t *testing.T) {
	first := book("Dune", nil, []string{"Frank Herbert"}, nil, nil)
	first.Path = "/Frank Herbert/1965/Dune.m4b"
	second := book("Dune", nil, []string{"Frank Herbert"}, nil, nil)
	second.Path = "/Frank Herbert/2006/Dune.m4b"
	other := book("What If? 2", nil, []string{"Randall Munroe"}, nil, nil)

	slugs := Plan([]audiobooks.Audiobook{second, other, first}, nil).BookSlugs

	assert.Equal(t, "dune", slugs[first.Path], "the first path in order keeps the plain slug")
	assert.Equal(t, "dune-2", slugs[second.Path])
	assert.Equal(t, "what-if-2", slugs[other.Path])
}

func TestBuildWritesAPageForEveryBook(t *testing.T) {
	root := t.TempDir()
	content := testPlan(t)
	_, err := Build(root, content, testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	for _, slug := range content.BookSlugs {
		assert.FileExists(t, filepath.Join(root, booksDir, slug, indexName))
	}
	assert.Contains(t, readPage(t, root, booksDir), `href="../books/a-wizard-of-earthsea/"`)
}

func TestBookPageShowsTheBook(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	page := readPage(t, root, filepath.Join(booksDir, "a-wizard-of-earthsea"))
	assert.Contains(t, page, "<h1>A Wizard of Earthsea</h1>")
	assert.Contains(
		t,
		page,
		`<source src="http://www.example-podcast.com/audiobooks/media/Ursula%20K%20Le%20Guin/Earthsea/1%20A%20Wizard%20of%20Earthsea/A%20Wizard%20of%20Earthsea.m4b" type="audio/mp4a-latm" />`,
	)
	assert.Contains(
		t,
		page,
		`<p class="description">Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.</p>`,
	)
	assert.NotContains(t, page, "&lt;p&gt;", "the description's own markup is dropped, not escaped onto the page")
	assert.Contains(t, page, `Subscribe to Ursula K. Le Guin</a>`)
	assert.Contains(t, page, `<img class="cover" src="../../covers/`)
	assert.Contains(t, page, "<dd>Earthsea book 1</dd>")
}

func TestBookPageEscapesTheDescription(t *testing.T) {
	audiobook := book("Risky", nil, []string{"A"}, nil, nil)
	audiobook.Description = &description.Description{
		Text:   "<p>Fine &lt;script&gt;alert(1)&lt;/script&gt;</p>",
		Format: description.HTML,
	}
	content := Plan([]audiobooks.Audiobook{audiobook}, nil)

	view := bookPageView(newLibrary(content, nil, nil), seriesInOrder(content.Books), &content.Books[0])
	require.Equal(t, []string{"Fine <script>alert(1)</script>"}, view.Book.Paragraphs)

	root := t.TempDir()
	_, err := Build(root, content, testRenderer(t), Options{}, testLogger())
	require.NoError(t, err)
	page := readPage(t, root, filepath.Join(booksDir, "risky"))
	assert.Contains(t, page, "Fine &lt;script&gt;alert(1)&lt;/script&gt;")
	assert.NotContains(t, page, "<script>alert")
}

func TestBookPageListsTheSeriesInReadingOrder(t *testing.T) {
	books := []audiobooks.Audiobook{
		inSeries("The Farthest Shore", "Earthsea", 3),
		inSeries("Tales from Earthsea", "Earthsea", 0),
		inSeries("A Wizard of Earthsea", "Earthsea", 1),
		inSeries("The Tombs of Atuan", "Earthsea", 2),
		book("The Left Hand of Darkness", nil, []string{"Ursula K. Le Guin"}, nil, nil),
		book("Someone Else's Book", nil, []string{"Someone Else"}, nil, nil),
	}
	content := Plan(books, nil)
	lib := newLibrary(content, nil, nil)
	current := &content.Books[3]

	view := bookPageView(lib, seriesInOrder(content.Books), current)

	order := make([]string, 0, len(view.Book.SeriesBooks))
	for index := range view.Book.SeriesBooks {
		entry := &view.Book.SeriesBooks[index]
		order = append(order, entry.Title)
		assert.Equal(t, entry.Title == current.Title, entry.Current, entry.Title)
	}
	assert.Equal(
		t,
		[]string{"A Wizard of Earthsea", "The Tombs of Atuan", "The Farthest Shore", "Tales from Earthsea"},
		order,
		"numbered books in order, then the unnumbered one",
	)
	assert.Equal(t, []string{"The Left Hand of Darkness"}, titles(view.Books),
		"more by the author, leaving out the series and other authors")
}

func TestBookPageWithoutASeriesListsNone(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{inSeries("Alone", "Solo", 1)}, nil)
	view := bookPageView(newLibrary(content, nil, nil), seriesInOrder(content.Books), &content.Books[0])
	assert.Empty(t, view.Book.SeriesBooks, "a series of one has nothing else to show")
	assert.Empty(t, view.Book.AudioURL, "no renderer, no player")
	assert.False(t, strings.Contains(view.Heading, "Solo"))
}

func titles(views []bookView) []string {
	names := make([]string, 0, len(views))
	for index := range views {
		names = append(names, views[index].Title)
	}
	return names
}

// Two authors can name a series the same thing, and neither's books belong on
// the other's page.
func TestBookPageSeriesIsTheAuthorsOwn(t *testing.T) {
	mine := inSeries("Mine One", "Chronicles", 1)
	mineToo := inSeries("Mine Two", "Chronicles", 2)
	theirs := inSeries("Theirs One", "Chronicles", 1)
	theirs.Authors = []string{"Someone Else"}
	content := Plan([]audiobooks.Audiobook{mine, mineToo, theirs}, nil)

	view := bookPageView(newLibrary(content, nil, nil), seriesInOrder(content.Books), &content.Books[0])

	order := make([]string, 0, len(view.Book.SeriesBooks))
	for index := range view.Book.SeriesBooks {
		order = append(order, view.Book.SeriesBooks[index].Title)
	}
	assert.Equal(t, []string{"Mine One", "Mine Two"}, order)
}
