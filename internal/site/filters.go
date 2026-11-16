package site

import (
	"math"
	"slices"
	"strings"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// sortOption is a way the filter script can reorder a list. The page's own
// order has an empty Key, so it is what the list shows without script.
type sortOption struct {
	Key   string
	Label string
}

// bookFilters are the controls on the All books page beyond the search box.
type bookFilters struct {
	Genres []linkView
	Sorts  []sortOption
}

var (
	bookSorts = []sortOption{
		{"", "Newest"},
		{"title", "Title"},
		{"longest", "Longest"},
		{"shortest", "Shortest"},
	}
	entrySorts = []sortOption{
		{"", "Name"},
		{"books", "Most books"},
		{"hours", "Most hours"},
	}
	// seriesSorts adds the series with the latest release, to find what has
	// moved on since you last listened.
	seriesSorts = append(slices.Clip(entrySorts), sortOption{"newest", "Newest release"})
)

// allBooksFilters offers a chip for each genre that has books, in the genres'
// own order, and each chip filters by the slug the books carry.
func allBooksFilters(lib *library) *bookFilters {
	filters := &bookFilters{Sorts: bookSorts}
	for index := range lib.Genres.Names {
		genre := &lib.Genres.Names[index]
		if len(genre.Books) > 0 {
			filters.Genres = append(filters.Genres, linkView{Name: genre.Title, Href: genre.Slug})
		}
	}
	return filters
}

// describeEntry describes a name or series with books, including what the
// filter script sorts it by.
func describeEntry(name, href, feed string, books []audiobooks.Audiobook) entryView {
	entry := entryView{
		Name:     name,
		Href:     href,
		FeedHref: feed,
		Search:   strings.ToLower(name),
		Count:    len(books),
		Hours:    totalHours(books),
		Minutes:  int(math.Round(totalDuration(books).Minutes())),
	}
	for index := range books {
		if released := releaseTime(&books[index]); !released.IsZero() {
			entry.Newest = max(entry.Newest, released.Format("2006-01-02"))
		}
	}
	return entry
}

// genreSlugs lists the slugs of a book's genre pages, for the genre chips.
func genreSlugs(lib *library, book *audiobooks.Audiobook) string {
	slugs := make([]string, 0, len(book.Genres))
	for _, name := range genreNames(book.Genres) {
		if slug, ok := lib.Genres.Slugs[name]; ok {
			slugs = append(slugs, slug)
		}
	}
	return strings.Join(slugs, " ")
}
