package site

import (
	"fmt"
	"path"
	"slices"
	"strings"
	"unicode"

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

	// The directory each section lives in, both under podcastDir for its feeds and
	// at the site root for its HTML pages, so /authors/ and /podcast/authors/ are
	// obviously the same set of names.
	authorsDir   = "authors"
	narratorsDir = "narrators"
	genresDir    = "genre"
	tagsDir      = "tags"
)

// Feed is one rendered feed together with every path it is published at. The
// aliases of a name all serve identical bytes, so they share a Feed and the
// feed is rendered once.
type Feed struct {
	Paths       []string
	Title       string
	Description string
	Books       []audiobooks.Audiobook
	// Link is the one of Paths that the HTML pages link to, or "" when no
	// spelling of the name can be a directory.
	Link string
}

// Name is one author, narrator, genre or tag: its feed, and the HTML page that
// lists its books.
type Name struct {
	Feed
	// Slug is the directory of the name's page within its section. Feed paths
	// keep the library's own spelling, because existing subscriptions point at
	// them, but the HTML tree is new, so its directories can be plain
	// lowercase-hyphen names that need no escaping and cannot fold together on a
	// case-insensitive filesystem.
	Slug string
}

// Section is one way of browsing the library, holding a Name per author,
// narrator, genre or tag.
type Section struct {
	Dir string
	// Names is sorted, and so is the order their slugs were assigned in.
	Names []Name
	// Listed maps each name exactly as a book lists it to its index in Names. It
	// is filled in while the books are grouped, so the HTML pages can link a
	// book to its names without repeating the grouping's normalisation.
	Listed map[string]int
}

// Content is everything a build publishes: the library the HTML pages list, the
// main feed, and a section for each way of browsing it.
type Content struct {
	Books     []audiobooks.Audiobook
	Main      Feed
	Authors   Section
	Narrators Section
	Genres    Section
	Tags      Section
}

// Plan works out the complete set of feeds and pages to publish for a library.
//
// Books are assumed to be sorted already; every feed preserves that order.
func Plan(books []audiobooks.Audiobook, excludedFromMainFeed []audiobooks.Genre) *Content {
	return &Content{
		Books: books,
		Main:  mainFeed(books, excludedFromMainFeed),
		Authors: personSection(books, authorsDir, byFormat, func(b *audiobooks.Audiobook) []string {
			return b.Authors
		}),
		Narrators: personSection(books, narratorsDir, readFormat, func(b *audiobooks.Audiobook) []string {
			return b.Narrators
		}),
		Genres: genreSection(books),
		Tags: personSection(books, tagsDir, plainFormat, func(b *audiobooks.Audiobook) []string {
			return b.Tags
		}),
	}
}

// sections lists the sections in the order the site presents them.
func sections(content *Content) []*Section {
	return []*Section{&content.Authors, &content.Narrators, &content.Genres, &content.Tags}
}

// feeds lists every feed in the plan, main feed first.
func feeds(content *Content) []*Feed {
	all := []*Feed{&content.Main}
	for _, section := range sections(content) {
		for index := range section.Names {
			all = append(all, &section.Names[index].Feed)
		}
	}
	return all
}

func mainFeed(books []audiobooks.Audiobook, excluded []audiobooks.Genre) Feed {
	included := make([]audiobooks.Audiobook, 0, len(books))
	for index := range books {
		if !hasAnyGenre(&books[index], excluded) {
			included = append(included, books[index])
		}
	}
	mainPath := path.Join(podcastDir, feedName)
	return Feed{
		Paths:       []string{mainPath},
		Title:       allTitle,
		Description: allSummary,
		Books:       included,
		Link:        mainPath,
	}
}

// genreSection holds every defined genre, including those with no books. The
// server returned an empty feed rather than a 404 for a valid but unused genre,
// and an existing subscription to one must not start 404ing.
func genreSection(books []audiobooks.Audiobook) Section {
	section := Section{Dir: genresDir, Listed: map[string]int{}}
	for _, genre := range audiobooks.AllGenres() {
		matching := []audiobooks.Audiobook{}
		for index := range books {
			if slices.Contains(books[index].Genres, genre) {
				matching = append(matching, books[index])
			}
		}
		link, paths := feedPaths(genresDir, genre.String(), genreAliases(genre))
		section.Listed[genre.String()] = len(section.Names)
		section.Names = append(section.Names, Name{Feed: Feed{
			Paths:       paths,
			Title:       genre.String(),
			Description: fmt.Sprintf(plainFormat, genre),
			Books:       matching,
			Link:        link,
		}})
	}
	assignSlugs(section.Names)
	return section
}

// personSection groups books by the names returned by names, matching the
// server's case- and whitespace-insensitive comparison.
//
// Unlike genres, a name with no books gets no file at all — its absence is what
// produces the 404 the server used to return.
func personSection(
	books []audiobooks.Audiobook,
	dir, descriptionFormat string,
	names func(*audiobooks.Audiobook) []string,
) Section {
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
	section := Section{Dir: dir, Names: make([]Name, 0, len(order)), Listed: map[string]int{}}
	for _, key := range order {
		name := canonical[key]
		var aliases []string
		for _, variant := range variants[key] {
			aliases = append(aliases, nameAliases(variant)...)
			section.Listed[variant] = len(section.Names)
		}
		link, paths := feedPaths(dir, name, dedupe(aliases))
		section.Names = append(section.Names, Name{Feed: Feed{
			Paths:       paths,
			Title:       name,
			Description: fmt.Sprintf(descriptionFormat, name),
			Books:       grouped[key],
			Link:        link,
		}})
	}
	assignSlugs(section.Names)
	return section
}

// assignSlugs gives every name in a section a unique page directory.
func assignSlugs(names []Name) {
	taken := map[string]bool{}
	for index := range names {
		slug := slugify(names[index].Title)
		// Distinct names can slugify identically — "Ann Leckie" and "Ann-Leckie" —
		// and must not end up sharing a page. Names are sorted, so which one keeps
		// the unsuffixed slug is stable between builds.
		unique := slug
		for suffix := 2; taken[unique]; suffix++ {
			unique = fmt.Sprintf("%s-%d", slug, suffix)
		}
		taken[unique] = true
		names[index].Slug = unique
	}
}

// slugify reduces a display name to a lowercase URL segment, replacing every run
// of other characters with a single hyphen.
func slugify(name string) string {
	var builder strings.Builder
	separated := false
	for _, char := range name {
		switch {
		case unicode.IsLetter(char) || unicode.IsDigit(char):
			if separated && builder.Len() > 0 {
				builder.WriteRune('-')
			}
			separated = false
			builder.WriteRune(unicode.ToLower(char))
		default:
			separated = true
		}
	}
	if builder.Len() == 0 {
		return "unnamed"
	}
	return builder.String()
}

// feedPaths turns alias names into output paths, dropping any that cannot
// safely be a directory name.
//
// The first return is the path the HTML pages link to: the display spelling
// where that can be a directory name, since that is the tidiest of the aliases,
// and otherwise whichever alias survived. A name with no usable spelling at all
// gets no feed, and so no link.
func feedPaths(dir, name string, aliases []string) (link string, all []string) {
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
