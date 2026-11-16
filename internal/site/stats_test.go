package site

import (
	"strconv"
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
	assert.Contains(t, page, "<h2>The longest listen</h2>")
	assert.Contains(t, page, `<span class="t">This Is How You Lose the Time War</span>`)
}

func TestStatsRankOnlyTheTopNarratorsWithBooks(t *testing.T) {
	books := make([]audiobooks.Audiobook, 0, topNarrators+2)
	for index := range topNarrators + 2 {
		books = append(books, released(
			"Book "+strconv.Itoa(index), 2000, time.Duration(index+1), "Narrator "+strconv.Itoa(index),
		))
	}
	content := Plan(books, nil)
	content.Narrators.Names = append(content.Narrators.Names, Name{Feed: Feed{Title: "Nobody"}, Slug: "nobody"})

	narrators := statsView(newLibrary(content, nil, nil)).Stats.Narrators

	require.Len(t, narrators, topNarrators)
	assert.Equal(t, "Narrator 11", narrators[0].Label)
	assert.Equal(t, "12 hours", narrators[0].Value)
	assert.Equal(t, "Narrator 2", narrators[topNarrators-1].Label)
}

func TestHoursLabel(t *testing.T) {
	assert.Equal(t, "1 hour", hoursLabel(time.Hour))
	assert.Equal(t, "1,234 hours", hoursLabel(1234*time.Hour))
	assert.Equal(t, "20 minutes", hoursLabel(20*time.Minute))
	assert.Equal(t, "1 minute", hoursLabel(time.Minute))
}
