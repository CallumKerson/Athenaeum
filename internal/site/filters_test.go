package site

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestAllBooksPageCarriesWhatTheFilterScriptNeeds(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	page := readPage(t, root, booksDir)
	assert.Contains(t, page, `<script src="../static/filter.js"></script>`)
	assert.FileExists(t, filepath.Join(root, "static", "filter.js"))
	assert.Contains(t, page, `data-title="A Wizard of Earthsea"`)
	assert.Contains(t, page, `data-genres="children-s fantasy"`)
	assert.Contains(
		t,
		page,
		`<button type="button" role="radio" data-genre="fantasy" aria-checked="false" tabindex="-1">Fantasy</button>`,
	)
	assert.Contains(t, page, `<option value="longest">Longest</option>`)
	assert.Contains(t, page, `<select id="length">`)
	assert.Contains(t, page, `data-view="list"`)
}

func TestIndexPagesOfferTheirSorts(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	authors := readPage(t, root, authorsDir)
	assert.Contains(t, authors, `<option value="hours">Most hours</option>`)
	assert.NotContains(t, authors, `value="newest"`, "only series sort by their latest release")
	assert.NotContains(t, authors, `<select id="length">`, "length and genre filters are for books")

	assert.Contains(t, readPage(t, root, seriesDir), `<option value="newest">Newest release</option>`)

	genres := readPage(t, root, "genre")
	assert.Contains(
		t,
		genres,
		`<option value="">Usual order</option>`,
		"genres are listed in their own order, not by name",
	)
	assert.Contains(t, genres, `<option value="name">Name</option>`)
}

func TestDescribeEntryKeepsItsSortKeys(t *testing.T) {
	older := dated("Older", &toml.LocalDate{Year: 2001, Month: 1, Day: 1})
	older.Duration = 90 * time.Minute
	newer := dated("Newer", &toml.LocalDate{Year: 2019, Month: 1, Day: 1})
	newer.Duration = 30 * time.Minute
	undated := book("Undated", nil, []string{"A"}, nil, nil)

	entry := describeEntry("Someone", "someone/", "", []audiobooks.Audiobook{older, undated, newer})

	assert.Equal(t, 3, entry.Count)
	assert.Equal(t, 120, entry.Minutes)
	assert.Equal(t, "2019-01-01", entry.Newest)
	assert.Equal(t, "someone", entry.Search)
}
