package site

import (
	"fmt"
	"path"
	"slices"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	podcastDir  = "podcast"
	feedName    = "feed.rss"
	allTitle    = "Audiobooks"
	allSummary  = "Like movies in your mind!"
	byFormat    = "Audiobooks by %s"
	readFormat  = "Audiobooks Narrated by %s"
	plainFormat = "%s Audiobooks"

	// The directory each kind of feed lives in, under podcastDir. The HTML tree
	// reuses these names one level up, so /authors/ and /podcast/authors/ are
	// obviously the same set of feeds.
	authorsDir   = "authors"
	narratorsDir = "narrators"
	genresDir    = "genre"
	tagsDir      = "tags"
)

// Page is one rendered feed together with every path it is published at. The
// aliases of a name all serve identical bytes, so they share a Page and the
// feed is rendered once.
type Page struct {
	Paths       []string
	Title       string
	Description string
	Books       []audiobooks.Audiobook
	// Section is the feed directory this page belongs to, or "" for the main feed.
	Section string
	// FeedPath is the one of Paths that the HTML pages link to.
	FeedPath string
}

// Content is everything a build publishes: the library the HTML pages list, and
// the feeds planned from it.
type Content struct {
	Books []audiobooks.Audiobook
	Feeds []Page
}

// Plan works out the complete set of feeds to publish for a library.
//
// Books are assumed to be sorted already; every feed preserves that order.
func Plan(books []audiobooks.Audiobook, excludedFromMainFeed []audiobooks.Genre) Content {
	pages := []Page{mainFeed(books, excludedFromMainFeed)}
	pages = append(pages, genreFeeds(books)...)
	pages = append(pages, personFeeds(books, authorsDir, byFormat, func(b *audiobooks.Audiobook) []string {
		return b.Authors
	})...)
	pages = append(pages, personFeeds(books, narratorsDir, readFormat, func(b *audiobooks.Audiobook) []string {
		return b.Narrators
	})...)
	pages = append(pages, personFeeds(books, tagsDir, plainFormat, func(b *audiobooks.Audiobook) []string {
		return b.Tags
	})...)
	return Content{Books: books, Feeds: pages}
}

func mainFeed(books []audiobooks.Audiobook, excluded []audiobooks.Genre) Page {
	included := make([]audiobooks.Audiobook, 0, len(books))
	for index := range books {
		if !hasAnyGenre(&books[index], excluded) {
			included = append(included, books[index])
		}
	}
	mainPath := path.Join(podcastDir, feedName)
	return Page{
		Paths:       []string{mainPath},
		Title:       allTitle,
		Description: allSummary,
		Books:       included,
		FeedPath:    mainPath,
	}
}

// genreFeeds emits a feed for every defined genre, including those with no
// books. The server returned an empty feed rather than a 404 for a valid but
// unused genre, and an existing subscription to one must not start 404ing.
func genreFeeds(books []audiobooks.Audiobook) []Page {
	pages := make([]Page, 0, len(audiobooks.AllGenres()))
	for _, genre := range audiobooks.AllGenres() {
		matching := []audiobooks.Audiobook{}
		for index := range books {
			if slices.Contains(books[index].Genres, genre) {
				matching = append(matching, books[index])
			}
		}
		canonical, paths := feedPaths(genresDir, genre.String(), genreAliases(genre))
		pages = append(pages, Page{
			Paths:       paths,
			Title:       genre.String(),
			Description: fmt.Sprintf(plainFormat, genre),
			Books:       matching,
			Section:     genresDir,
			FeedPath:    canonical,
		})
	}
	return pages
}

// personFeeds groups books by the names returned by names, matching the
// server's case- and whitespace-insensitive comparison.
//
// Unlike genres, a name with no books gets no file at all — its absence is what
// produces the 404 the server used to return.
func personFeeds(
	books []audiobooks.Audiobook,
	dir, descriptionFormat string,
	names func(*audiobooks.Audiobook) []string,
) []Page {
	// Keyed by normalised name so spelling variants in the metadata collapse into
	// one feed, as they did when the server filtered at request time.
	grouped := map[string][]audiobooks.Audiobook{}
	// The canonical display name is the first spelling encountered, which is
	// stable because the book order is.
	canonical := map[string]string{}
	// Every spelling that appears in the metadata, not just the canonical one.
	// The library really does contain both "V.E. Schwab" and "V. E. Schwab", and
	// the server served a feed at each — so each needs its own aliases.
	variants := map[string][]string{}
	var order []string

	for index := range books {
		// One book can list several spellings that normalise to the same key, so
		// it would otherwise be appended to that feed once per spelling — a
		// duplicate item sharing a GUID with the original.
		counted := map[string]bool{}
		for _, name := range names(&books[index]) {
			key := normalise(name)
			if key == "" {
				continue
			}
			if _, known := canonical[key]; !known {
				canonical[key] = name
				order = append(order, key)
			}
			if !slices.Contains(variants[key], name) {
				variants[key] = append(variants[key], name)
			}
			if !counted[key] {
				counted[key] = true
				grouped[key] = append(grouped[key], books[index])
			}
		}
	}

	slices.Sort(order)
	pages := make([]Page, 0, len(order))
	for _, key := range order {
		name := canonical[key]
		var aliases []string
		for _, variant := range variants[key] {
			aliases = append(aliases, nameAliases(variant)...)
		}
		feedPath, paths := feedPaths(dir, name, dedupe(aliases))
		pages = append(pages, Page{
			Paths:       paths,
			Title:       name,
			Description: fmt.Sprintf(descriptionFormat, name),
			Books:       grouped[key],
			Section:     dir,
			FeedPath:    feedPath,
		})
	}
	return pages
}

// feedPaths turns alias names into output paths, dropping any that cannot
// safely be a directory name.
//
// The first return is the path the HTML pages link to: the display spelling
// where that can be a directory name, since that is the tidiest of the aliases,
// and otherwise whichever alias survived. A name with no usable spelling at all
// gets no feed, and so no link.
func feedPaths(dir, name string, aliases []string) (canonical string, all []string) {
	paths := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		if !safeSegment(alias) {
			continue
		}
		paths = append(paths, path.Join(podcastDir, dir, alias, feedName))
	}

	switch {
	case safeSegment(name):
		return path.Join(podcastDir, dir, name, feedName), paths
	case len(paths) > 0:
		return paths[0], paths
	default:
		return "", paths
	}
}

func hasAnyGenre(book *audiobooks.Audiobook, genres []audiobooks.Genre) bool {
	for _, genre := range genres {
		if slices.Contains(book.Genres, genre) {
			return true
		}
	}
	return false
}
