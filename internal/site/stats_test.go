package site

import (
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func released(title string, year int, hours time.Duration, narrators ...string) audiobooks.Audiobook {
	audiobook := dated(title, &toml.LocalDate{Year: year, Month: 1, Day: 1})
	audiobook.Duration = hours * time.Hour
	audiobook.Narrators = narrators
	return audiobook
}

func TestStatsCountTheWholeLibrary(t *testing.T) {
	books := []audiobooks.Audiobook{
		released("Old", 1968, 10, "Kobna"),
		released("Older", 1961, 2, "Kobna"),
		released("New", 2019, 30, "Rosamund"),
		book("Undated", nil, []string{"A"}, nil, nil),
	}
	content := Plan(books, nil)
	stats := statsView(newLibrary(content, nil, nil)).Stats

	assert.Equal(t, statTotal{"4", "books"}, stats.Totals[0])
	assert.Equal(t, statTotal{"42", "hours of listening"}, stats.Totals[1])

	require.Len(t, stats.Narrators, 2)
	assert.Equal(t, "Rosamund", stats.Narrators[0].Label, "most hours first")
	assert.Equal(t, 100, stats.Narrators[0].Share)
	assert.Equal(t, "Kobna", stats.Narrators[1].Label)
	assert.Equal(t, 40, stats.Narrators[1].Share)
	assert.Equal(t, "2 books", stats.Narrators[1].Note)

	require.Len(t, stats.Decades, 2, "undated books are left out")
	assert.Equal(t, statBar{Label: "1960s", Value: "2 books", Share: 100}, stats.Decades[0])
	assert.Equal(t, statBar{Label: "2010s", Value: "1 book", Share: 50}, stats.Decades[1])

	require.NotNil(t, stats.Longest)
	assert.Equal(t, "New", stats.Longest.Title)
}

func TestShareShowsEverythingWithAnyValue(t *testing.T) {
	assert.Equal(t, 0, share(0, 10))
	assert.Equal(t, 1, share(0.01, 10))
	assert.Equal(t, 100, share(10, 10))
	assert.Equal(t, 0, share(5, 0))
}

func TestBuildWritesTheStatsPage(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	page := readPage(t, root, statsDir)
	assert.Contains(t, page, `aria-current="page">Stats</a>`)
	assert.Contains(t, page, "<h2>Hours by genre</h2>")
	assert.Contains(t, page, `style="--share: 100%"`)
	assert.Contains(t, page, `href="../genre/fantasy/"`)
}
