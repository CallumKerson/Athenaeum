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

func feedAt(t *testing.T, content *Content, path string) Feed {
	t.Helper()
	planned := feeds(content)
	index := slices.IndexFunc(planned, func(feed *Feed) bool { return slices.Contains(feed.Paths, path) })
	require.GreaterOrEqual(t, index, 0, "no feed published at %s", path)
	return *planned[index]
}

func TestPlanExcludesGenresFromMainFeedOnly(t *testing.T) {
	books := []audiobooks.Audiobook{
		book("Clean", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
		book("Spicy", []audiobooks.Genre{audiobooks.Erotica}, []string{"B"}, nil, nil),
	}

	pages := Plan(books, []audiobooks.Genre{audiobooks.Erotica})

	main := feedAt(t, pages, "podcast/feed.rss")
	require.Len(t, main.Books, 1)
	assert.Equal(t, "Clean", main.Books[0].Title)

	// The excluded genre still gets its own feed, and the author of an excluded
	// book still gets theirs — only the main feed filters.
	assert.Len(t, feedAt(t, pages, "podcast/genre/Erotica/feed.rss").Books, 1)
	assert.Len(t, feedAt(t, pages, "podcast/authors/B/feed.rss").Books, 1)
}

// A valid genre with no books returned an empty feed rather than a 404, so an
// existing subscription to one must not start failing.
func TestPlanPublishesEveryGenreEvenWhenEmpty(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
	}, nil)

	for _, genre := range audiobooks.AllGenres() {
		page := feedAt(t, pages, "podcast/genre/"+genre.String()+"/feed.rss")
		assert.Equal(t, genre.String(), page.Title)
	}
	assert.Empty(t, feedAt(t, pages, "podcast/genre/Horror/feed.rss").Books)
}

// Authors, narrators and tags 404 when they have no books, which a static tree
// expresses by not writing the file at all.
func TestPlanOmitsNamesWithNoBooks(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", []audiobooks.Genre{audiobooks.Fantasy}, []string{"A"}, nil, nil),
	}, nil)

	published := map[string]bool{}
	for _, planned := range feeds(pages) {
		for _, path := range planned.Paths {
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

	page := feedAt(t, pages, "podcast/narrators/Full Cast/feed.rss")
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

	page := feedAt(t, pages, "podcast/authors/V.E. Schwab/feed.rss")
	assert.Len(t, page.Books, 2)
	assert.Equal(t, "V.E. Schwab", page.Title)
	assert.Contains(t, page.Paths, "podcast/authors/V. E. Schwab/feed.rss")
	assert.Contains(t, page.Paths, "podcast/authors/v. e. schwab/feed.rss")
	assert.Contains(t, page.Paths, "podcast/authors/v.e.schwab/feed.rss")
}

// A book listing two spellings of the same name belongs in that feed once. The
// server filtered per book, so an item appearing twice — with the same GUID as
// itself — would be a regression as well as a bug.
func TestPlanCountsABookOncePerFeed(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", nil, []string{"V.E. Schwab", "V. E. Schwab"}, nil, []string{"tag", "tag"}),
	}, nil)

	assert.Len(t, feedAt(t, pages, "podcast/authors/V.E. Schwab/feed.rss").Books, 1)
	assert.Len(t, feedAt(t, pages, "podcast/tags/tag/feed.rss").Books, 1)
}

// Unlike the server, which echoed the requested spelling back as the feed title,
// every alias now carries the canonical name.
func TestPlanAliasesShareCanonicalTitle(t *testing.T) {
	pages := Plan([]audiobooks.Audiobook{
		book("Only", nil, []string{"Adrian Tchaikovsky"}, nil, nil),
	}, nil)

	page := feedAt(t, pages, "podcast/authors/adriantchaikovsky/feed.rss")
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

// slugsOf lists each name in a section with the slug of its page.
func slugsOf(section *Section) map[string]string {
	slugs := map[string]string{}
	for index := range section.Names {
		slugs[section.Names[index].Title] = section.Names[index].Slug
	}
	return slugs
}

func TestPlanAssignsEveryNameAPage(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		book("Only", []audiobooks.Genre{audiobooks.SciFi}, []string{"Ursula K. Le Guin"}, []string{"Full Cast"},
			[]string{"Hugo Awards"}),
	}, nil)

	assert.Equal(t, map[string]string{"Ursula K. Le Guin": "ursula-k-le-guin"}, slugsOf(&content.Authors))
	assert.Equal(t, map[string]string{"Full Cast": "full-cast"}, slugsOf(&content.Narrators))
	assert.Equal(t, map[string]string{"Hugo Awards": "hugo-awards"}, slugsOf(&content.Tags))
	// A genre with nothing filed under it keeps its feed, so it keeps its page.
	assert.Equal(t, "science-fiction", slugsOf(&content.Genres)["Science Fiction"])
	assert.Equal(t, "horror", slugsOf(&content.Genres)["Horror"])
}

// Two names that slugify the same way must not share a page, and which one keeps
// the plain slug must not depend on the order the library lists them in. Names
// are sorted by their normalised form, so "ann-leckie" comes before "annleckie".
func TestPlanGivesCollidingNamesDistinctSlugsInSortedOrder(t *testing.T) {
	collide := func(authors ...string) map[string]string {
		books := make([]audiobooks.Audiobook, 0, len(authors))
		for _, author := range authors {
			books = append(books, book(author, nil, []string{author}, nil, nil))
		}
		content := Plan(books, nil)
		return slugsOf(&content.Authors)
	}

	expected := map[string]string{
		"Ann-Leckie": "ann-leckie",
		"ann.leckie": "ann-leckie-2",
		"Ann Leckie": "ann-leckie-3",
	}
	assert.Equal(t, expected, collide("Ann Leckie", "Ann-Leckie", "ann.leckie"))
	assert.Equal(t, expected, collide("ann.leckie", "Ann Leckie", "Ann-Leckie"))
}

// A book links to a name under whichever spelling it lists, and every spelling
// of one name leads to the same page.
func TestPlanResolvesListedNamesToTheirPages(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		book("First", []audiobooks.Genre{audiobooks.SciFi}, []string{"V.E. Schwab"}, []string{"Full Cast"}, nil),
		book("Second", nil, []string{"V. E. Schwab", "v.e. schwab"}, []string{"full cast"}, []string{"Tag"}),
	}, nil)

	pageOf := func(section *Section, listed string) string {
		t.Helper()
		index, ok := section.Listed[listed]
		require.True(t, ok, "%q has no page in %s", listed, section.Dir)
		return section.Names[index].Slug
	}

	assert.Equal(t, "v-e-schwab", pageOf(&content.Authors, "V.E. Schwab"))
	assert.Equal(t, "v-e-schwab", pageOf(&content.Authors, "V. E. Schwab"))
	assert.Equal(t, "v-e-schwab", pageOf(&content.Authors, "v.e. schwab"))
	assert.Equal(t, "full-cast", pageOf(&content.Narrators, "Full Cast"))
	assert.Equal(t, "full-cast", pageOf(&content.Narrators, "full cast"))
	assert.Equal(t, "tag", pageOf(&content.Tags, "Tag"))
	assert.Equal(t, "science-fiction", pageOf(&content.Genres, audiobooks.SciFi.String()))
}

// A name that is blank once whitespace is removed gets no feed, so it has no
// page to link to either.
func TestPlanGivesBlankNamesNoPage(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		book("Only", nil, []string{"A"}, []string{"  "}, nil),
	}, nil)

	assert.Empty(t, content.Narrators.Names)
	assert.NotContains(t, content.Narrators.Listed, "  ")
}

// A name with no spelling that is safe as a directory still gets a page, which
// just has no feed to link to.
func TestPlanGivesNamesWithNoFeedPathAPage(t *testing.T) {
	content := Plan([]audiobooks.Audiobook{
		book("Only", nil, []string{".."}, nil, nil),
	}, nil)

	require.Len(t, content.Authors.Names, 1)
	assert.Empty(t, content.Authors.Names[0].Paths)
	assert.Empty(t, content.Authors.Names[0].Link)
	assert.Equal(t, "unnamed", content.Authors.Names[0].Slug)
}
