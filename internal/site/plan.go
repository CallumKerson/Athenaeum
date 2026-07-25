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
)

// Page is one rendered feed together with every path it is published at. The
// aliases of a name all serve identical bytes, so they share a Page and the
// feed is rendered once.
type Page struct {
	Paths       []string
	Title       string
	Description string
	Books       []audiobooks.Audiobook
}

// Plan works out the complete set of feeds to publish for a library.
//
// Books are assumed to be sorted already; every feed preserves that order.
func Plan(books []audiobooks.Audiobook, excludedFromMainFeed []audiobooks.Genre) []Page {
	pages := []Page{mainFeed(books, excludedFromMainFeed)}
	pages = append(pages, genreFeeds(books)...)
	pages = append(pages, personFeeds(books, "authors", byFormat, func(b *audiobooks.Audiobook) []string {
		return b.Authors
	})...)
	pages = append(pages, personFeeds(books, "narrators", readFormat, func(b *audiobooks.Audiobook) []string {
		return b.Narrators
	})...)
	pages = append(pages, personFeeds(books, "tags", plainFormat, func(b *audiobooks.Audiobook) []string {
		return b.Tags
	})...)
	return pages
}

func mainFeed(books []audiobooks.Audiobook, excluded []audiobooks.Genre) Page {
	included := make([]audiobooks.Audiobook, 0, len(books))
	for index := range books {
		if !hasAnyGenre(&books[index], excluded) {
			included = append(included, books[index])
		}
	}
	return Page{
		Paths:       []string{path.Join(podcastDir, feedName)},
		Title:       allTitle,
		Description: allSummary,
		Books:       included,
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
		pages = append(pages, Page{
			Paths:       feedPaths("genre", genreAliases(genre)),
			Title:       genre.String(),
			Description: fmt.Sprintf(plainFormat, genre),
			Books:       matching,
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
			grouped[key] = append(grouped[key], books[index])
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
		pages = append(pages, Page{
			Paths:       feedPaths(dir, dedupe(aliases)),
			Title:       name,
			Description: fmt.Sprintf(descriptionFormat, name),
			Books:       grouped[key],
		})
	}
	return pages
}

// feedPaths turns alias names into output paths, dropping any that cannot
// safely be a directory name.
func feedPaths(dir string, aliases []string) []string {
	paths := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		if !safeSegment(alias) {
			continue
		}
		paths = append(paths, path.Join(podcastDir, dir, alias, feedName))
	}
	return paths
}

func hasAnyGenre(book *audiobooks.Audiobook, genres []audiobooks.Genre) bool {
	for _, genre := range genres {
		if slices.Contains(book.Genres, genre) {
			return true
		}
	}
	return false
}
