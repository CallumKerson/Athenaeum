package site

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func seriesTitles(series *Series) []string {
	names := make([]string, 0, len(series.Books))
	for index := range series.Books {
		names = append(names, series.Books[index].Title)
	}
	return names
}

func TestPlanGroupsSeriesInReadingOrder(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		inSeries("The Farthest Shore", "Earthsea", 3),
		inSeries("Tales from Earthsea", "Earthsea", 0),
		inSeries("A Wizard of Earthsea", "Earthsea", 1),
		inSeries("The Tombs of Atuan", "Earthsea", 2),
		book("The Left Hand of Darkness", nil, []string{"Ursula K. Le Guin"}, nil, nil),
	}, Exclusions{})

	require.Len(t, content.Series, 1)
	assert.Equal(t, "earthsea", content.Series[0].Slug)
	assert.Equal(
		t,
		[]string{"A Wizard of Earthsea", "The Tombs of Atuan", "The Farthest Shore", "Tales from Earthsea"},
		seriesTitles(&content.Series[0]),
	)
}

// Two authors can use the same series title, and each gets their own series
// and page, while one author's books join up even through a co-written book.
func TestPlanSplitsASeriesTitleByAuthor(t *testing.T) {
	mine := inSeries("Mine One", "Chronicles", 1)
	cowritten := inSeries("Ours Two", "Chronicles", 2)
	cowritten.Authors = []string{"Ursula K. Le Guin", "A Friend"}
	friends := inSeries("Friend's Three", "Chronicles", 3)
	friends.Authors = []string{"A Friend"}
	theirs := inSeries("Theirs One", "Chronicles", 1)
	theirs.Authors = []string{"Someone Else"}

	content := Plan([]audiobooks.Audiobook{theirs, friends, mine, cowritten}, Exclusions{})

	require.Len(t, content.Series, 2)
	assert.Equal(t, []string{"Mine One", "Ours Two", "Friend's Three"}, seriesTitles(&content.Series[0]))
	assert.Equal(t, "chronicles", content.Series[0].Slug, "the group with the first path keeps the plain slug")
	assert.Equal(t, []string{"Theirs One"}, seriesTitles(&content.Series[1]))
	assert.Equal(t, "chronicles-2", content.Series[1].Slug)
	assert.Equal(t, "Chronicles (Ursula K. Le Guin, A Friend)", content.Series[0].Name,
		"lists name the authors when two series share a title")
	assert.Equal(t, "Chronicles (Someone Else)", content.Series[1].Name)

	root := t.TempDir()
	_, err := Build(root, content, testRenderer(t), Options{}, testLogger())
	require.NoError(t, err)
	page := readPage(t, root, filepath.Join(booksDir, content.BookSlugs[theirs.Path]))
	assert.Contains(t, page, `<a href="../../series/chronicles-2/">Chronicles book 1</a>`)
	page = readPage(t, root, filepath.Join(booksDir, content.BookSlugs[mine.Path]))
	assert.Contains(t, page, `<a href="../../series/chronicles/">The whole series</a>`)
}

func TestPlanMatchesSeriesTitlesLikeOtherNames(t *testing.T) {
	first := inSeries("Leviathan Wakes", "The Expanse", 1)
	second := inSeries("Caliban's War", "the expanse", 2)
	content := Plan([]audiobooks.Audiobook{second, first}, Exclusions{})

	require.Len(t, content.Series, 1)
	assert.Equal(t, "The Expanse", content.Series[0].Title, "the first book's spelling")
	assert.Equal(t, "The Expanse", content.Series[0].Name)
}

func TestPlanGroupsSeriesBooksWithoutAuthors(t *testing.T) {
	first := inSeries("One", "Anonymous Tales", 1)
	first.Authors = nil
	second := inSeries("Two", "Anonymous Tales", 2)
	second.Authors = nil
	content := Plan([]audiobooks.Audiobook{first, second}, Exclusions{})

	require.Len(t, content.Series, 1)
	assert.Equal(t, []string{"One", "Two"}, seriesTitles(&content.Series[0]))
}

func TestBuildWritesSeriesPages(t *testing.T) {
	root := t.TempDir()
	content := testPlan(t)
	_, err := Build(root, content, testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	index := readPage(t, root, seriesDir)
	assert.Contains(t, index, `href="../series/earthsea/"`)
	assert.Contains(t, index, `aria-current="page">Series</a>`)

	page := readPage(t, root, filepath.Join(seriesDir, "earthsea"))
	assert.Contains(t, page, "<h1>Earthsea</h1>")
	assert.Contains(t, page, "By Ursula K. Le Guin. 1 book.")
	assert.Contains(t, page, `<span class="seq">Book 1</span>`)
	assert.Contains(t, page, `href="../../books/a-wizard-of-earthsea/"`)
}

func TestBuildWithoutSeriesLeavesOutTheSection(t *testing.T) {
	root := t.TempDir()
	content := Plan([]audiobooks.Audiobook{book("Alone", nil, []string{"A"}, nil, nil)}, Exclusions{})
	_, err := Build(root, content, testRenderer(t), Options{}, testLogger())
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, seriesDir))
	assert.NotContains(t, readPage(t, root, "."), ">Series</a>")
}
