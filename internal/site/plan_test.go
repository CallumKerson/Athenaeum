package site

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func book(title string, genres []audiobooks.Genre, authors, narrators, tags []string) audiobooks.Audiobook {
	return audiobooks.Audiobook{
		Title:     title,
		Authors:   authors,
		Narrators: narrators,
		Tags:      tags,
		Genres:    genres,
		Path:      "/" + title + "/" + title + ".m4b",
		MIMEType:  "audio/mp4a-latm",
	}
}

func pageFor(t *testing.T, pages []Page, path string) Page {
	t.Helper()
	index := slices.IndexFunc(pages, func(page Page) bool { return slices.Contains(page.Paths, path) })
	require.GreaterOrEqual(t, index, 0, "no page published at %s", path)
	return pages[index]
}

func TestPlanExcludesGenresFromMainFeedOnly(t *testing.T) {
	books := []audiobooks.Audiobook{
		book("Clean", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
		book("Spicy", []audiobooks.Genre{audiobooks.Erotica}, []string{"B"}, nil, nil),
	}

	pages := Plan(books, []audiobooks.Genre{audiobooks.Erotica})

	main := pageFor(t, pages, "podcast/feed.rss")
	require.Len(t, main.Books, 1)
	assert.Equal(t, "Clean", main.Books[0].Title)

	// The excluded genre still gets its own feed, and the author of an excluded
	// book still gets theirs — only the main feed filters.
	assert.Len(t, pageFor(t, pages, "podcast/genre/Erotica/feed.rss").Books, 1)
	assert.Len(t, pageFor(t, pages, "podcast/authors/B/feed.rss").Books, 1)
}

// A valid genre with no books returned an empty feed rather than a 404, so an
// existing subscription to one must not start failing.
func TestPlanPublishesEveryGenreEvenWhenEmpty(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
	}, nil)

	for _, genre := range audiobooks.AllGenres() {
		page := pageFor(t, pages, "podcast/genre/"+genre.String()+"/feed.rss")
		assert.Equal(t, genre.String(), page.Title)
	}
	assert.Empty(t, pageFor(t, pages, "podcast/genre/Horror/feed.rss").Books)
}

// Authors, narrators and tags 404 when they have no books, which a static tree
// expresses by not writing the file at all.
func TestPlanOmitsNamesWithNoBooks(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
	}, nil)

	published := map[string]bool{}
	for index := range pages {
		for _, path := range pages[index].Paths {
			published[path] = true
		}
	}

	assert.True(t, published["podcast/authors/A/feed.rss"])
	assert.False(t, published["podcast/narrators/A/feed.rss"])
	assert.False(t, published["podcast/tags/A/feed.rss"])
}

// The server compared names case- and whitespace-insensitively, so metadata
// spelling variants have to end up in one feed, not two.
func TestPlanGroupsNameVariants(t *testing.T) {
	books := []audiobooks.Audiobook{
		book("First", nil, []string{"A"}, []string{"Full Cast"}, nil),
		book("Second", nil, []string{"A"}, []string{"full cast"}, nil),
	}

	pages := Plan(books, nil)

	page := pageFor(t, pages, "podcast/narrators/Full Cast/feed.rss")
	assert.Len(t, page.Books, 2)
	// The first spelling seen becomes the canonical title, and every variant is
	// published as an alias of it.
	assert.Equal(t, "Full Cast", page.Title)
	assert.Contains(t, page.Paths, "podcast/narrators/full cast/feed.rss")
	assert.Contains(t, page.Paths, "podcast/narrators/fullcast/feed.rss")
}

// Two spellings that differ by more than case — the library holds both
// "V.E. Schwab" and "V. E. Schwab" — each need their own directory, because
// neither is derivable from the other.
func TestPlanPublishesEverySpellingSeen(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("First", nil, []string{"V.E. Schwab"}, nil, nil),
		book("Second", nil, []string{"V. E. Schwab"}, nil, nil),
	}, nil)

	page := pageFor(t, pages, "podcast/authors/V.E. Schwab/feed.rss")
	assert.Len(t, page.Books, 2)
	assert.Equal(t, "V.E. Schwab", page.Title)
	assert.Contains(t, page.Paths, "podcast/authors/V. E. Schwab/feed.rss")
	assert.Contains(t, page.Paths, "podcast/authors/v. e. schwab/feed.rss")
	assert.Contains(t, page.Paths, "podcast/authors/v.e.schwab/feed.rss")
}

// Unlike the server, which echoed the requested spelling back as the feed title,
// every alias now carries the canonical name.
func TestPlanAliasesShareCanonicalTitle(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", nil, []string{"Adrian Tchaikovsky"}, nil, nil),
	}, nil)

	page := pageFor(t, pages, "podcast/authors/adriantchaikovsky/feed.rss")
	assert.Equal(t, "Adrian Tchaikovsky", page.Title)
	assert.Equal(t, "Audiobooks by Adrian Tchaikovsky", page.Description)
}

func TestPlanIsDeterministic(t *testing.T) {
	books := []audiobooks.Audiobook{
		book("First", []audiobooks.Genre{audiobooks.Fantasy}, []string{"B", "A"}, []string{"N"}, []string{"T"}),
		book("Second", []audiobooks.Genre{audiobooks.SciFi}, []string{"A"}, []string{"M"}, nil),
	}

	assert.Equal(t, Plan(books, nil), Plan(books, nil))
}
