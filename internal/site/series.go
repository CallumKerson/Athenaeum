package site

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	seriesDir   = "series"
	seriesTitle = "Series"
	seriesBlurb = "Every series in the library, with its books in reading order."
)

// Series is the books that share a series title and an author, in reading
// order. Different authors can use the same title, so one title can make more
// than one Series. It has a page but no feed, so adding series published no
// new feed paths.
type Series struct {
	Title string
	// Slug is the directory of the series' page under /series/.
	Slug  string
	Books []audiobooks.Audiobook
}

// seriesGroups groups books into series. Within a title, books join a series
// when they share an author with any book already in it, comparing authors by
// their slugs so that spelling variants still match. Titles sort by normalised
// name and groups within a title by their first path, so which series keeps
// the plain slug is stable between builds.
func seriesGroups(books []audiobooks.Audiobook, authorSlugs map[string]string) []Series {
	byTitle := map[string][]*audiobooks.Audiobook{}
	for index := range books {
		book := &books[index]
		if book.Series != nil && strings.TrimSpace(book.Series.Title) != "" {
			byTitle[book.Series.Title] = append(byTitle[book.Series.Title], book)
		}
	}
	titles := make([]string, 0, len(byTitle))
	for title := range byTitle {
		titles = append(titles, title)
	}
	slices.SortFunc(titles, func(a, b string) int {
		return cmp.Or(strings.Compare(normalise(a), normalise(b)), strings.Compare(a, b))
	})

	var groups []Series
	taken := map[string]bool{}
	for _, title := range titles {
		for _, members := range byAuthor(byTitle[title], authorSlugs) {
			sortReadingOrder(members)
			group := Series{Title: title, Slug: claimSlug(taken, slugify(title))}
			for _, member := range members {
				group.Books = append(group.Books, *member)
			}
			groups = append(groups, group)
		}
	}
	return groups
}

// byAuthor splits books with one series title into groups linked by shared
// authors, ordered by each group's first path.
func byAuthor(books []*audiobooks.Audiobook, authorSlugs map[string]string) [][]*audiobooks.Audiobook {
	slices.SortFunc(books, func(a, b *audiobooks.Audiobook) int { return strings.Compare(a.Path, b.Path) })

	var groups [][]*audiobooks.Audiobook
	for _, book := range books {
		var joined []*audiobooks.Audiobook
		kept := groups[:0:0]
		for _, group := range groups {
			if slices.ContainsFunc(group, func(other *audiobooks.Audiobook) bool {
				return sameAuthor(authorSlugs, book, other)
			}) {
				joined = append(joined, group...)
				continue
			}
			kept = append(kept, group)
		}
		joined = append(joined, book)
		kept = append(kept, joined)
		groups = kept
	}
	slices.SortFunc(groups, func(a, b []*audiobooks.Audiobook) int {
		return strings.Compare(firstPath(a), firstPath(b))
	})
	return groups
}

func firstPath(books []*audiobooks.Audiobook) string {
	first := books[0].Path
	for _, book := range books[1:] {
		first = min(first, book.Path)
	}
	return first
}

func sameAuthor(authorSlugs map[string]string, first, second *audiobooks.Audiobook) bool {
	for _, mine := range first.Authors {
		for _, theirs := range second.Authors {
			slug, ok := authorSlugs[mine]
			if ok && slug == authorSlugs[theirs] {
				return true
			}
		}
	}
	return false
}

// sortReadingOrder puts numbered books in order of their position, then the
// unnumbered ones by title.
func sortReadingOrder(books []*audiobooks.Audiobook) {
	slices.SortStableFunc(books, func(a, b *audiobooks.Audiobook) int {
		return cmp.Or(
			comparePositions(a.Series.Sequence, b.Series.Sequence),
			strings.Compare(a.Title, b.Title),
			strings.Compare(a.Path, b.Path),
		)
	})
}

func comparePositions(first, second *audiobooks.Sequence) int {
	switch {
	case first == nil && second == nil:
		return 0
	case first == nil:
		return 1
	case second == nil:
		return -1
	default:
		return first.First.Cmp(second.First)
	}
}

// seriesEntries describes a series' books for a page whose prefix back to the
// site root is root, marking the book at current, if any.
func seriesEntries(lib *library, root string, series *Series, current string) []seriesEntry {
	entries := make([]seriesEntry, 0, len(series.Books))
	for index := range series.Books {
		book := &series.Books[index]
		entries = append(entries, seriesEntry{
			bookView: describeBook(lib, root, book),
			Position: position(book),
			Current:  book.Path == current,
		})
	}
	return entries
}

func seriesIndexView(lib *library) *pageView {
	root := rootFromDir
	view := &pageView{
		Title:       pageTitle(seriesTitle),
		Heading:     seriesTitle,
		Blurb:       seriesBlurb,
		Root:        root,
		Nav:         navLinks(lib, root, seriesDir),
		Search:      true,
		SearchLabel: "Search series",
		Entries:     make([]entryView, 0, len(lib.Series)),
	}
	for index := range lib.Series {
		series := &lib.Series[index]
		view.Entries = append(view.Entries, entryView{
			Name:   series.Title,
			Href:   dirHref(root, seriesDir, series.Slug),
			Search: strings.ToLower(series.Title + " " + strings.Join(seriesAuthors(series), " ")),
			Count:  len(series.Books),
			Hours:  totalHours(series.Books),
		})
	}
	return view
}

func seriesPageView(lib *library, series *Series) *pageView {
	root := rootFromSubDir
	blurb := fmt.Sprintf(
		"%d %s, %s hours.",
		len(series.Books),
		plural(len(series.Books), "book"),
		totalHours(series.Books),
	)
	if authors := seriesAuthors(series); len(authors) > 0 {
		blurb = "By " + strings.Join(authors, ", ") + ". " + blurb
	}
	return &pageView{
		Title:   pageTitle(series.Title),
		Heading: series.Title,
		Blurb:   blurb,
		Root:    root,
		Nav:     navLinks(lib, root, seriesDir),
		Series:  seriesEntries(lib, root, series, ""),
	}
}

func plural(count int, noun string) string {
	if count == 1 {
		return noun
	}
	return noun + "s"
}

// seriesAuthors lists everyone who wrote a book in the series, in the order
// they first appear.
func seriesAuthors(series *Series) []string {
	var authors []string
	for index := range series.Books {
		for _, author := range series.Books[index].Authors {
			if !slices.Contains(authors, author) {
				authors = append(authors, author)
			}
		}
	}
	return authors
}
