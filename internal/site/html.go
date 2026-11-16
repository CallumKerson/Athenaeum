package site

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/templates"
)

// The searchable list of the whole library, and what it is called.
const (
	booksDir   = "books"
	booksTitle = "All books"
)

// The relative prefix back to the site root, by how deep the page is. Links are
// relative so that the site works wherever it is mounted — the feeds are served
// from a host that may itself include a path.
const (
	rootFromRoot   = "./"
	rootFromDir    = "../"
	rootFromSubDir = "../../"
)

// shownSection is a section of the plan together with how the site introduces
// it. Keeping the text beside the section, rather than in a lookup keyed by
// directory, means no section can be shown without its heading.
type shownSection struct {
	*Section
	Title string
	Blurb string
}

// shownSections lists the plan's sections in the order the site presents them.
// TestEverySectionIsShown checks it against sections, so a section added to the
// plan cannot quietly go without pages.
func shownSections(content *Content) []shownSection {
	return []shownSection{
		{&content.Authors, "Authors", "Everyone who wrote something in the library."},
		{&content.Narrators, "Narrators", "Everyone who read something in the library."},
		{&content.Genres, "Genres", "Every genre Athenaeum knows, including those nothing is filed under yet."},
		{&content.Tags, "Tags", "Whatever the library's own metadata tags books with."},
	}
}

// htmlPage is a rendered page and the path it is published at.
type htmlPage struct {
	Path     string
	Contents []byte
}

type linkView struct {
	Name string
	Href string
}

type entryView struct {
	Name     string
	Href     string
	FeedHref string
	Note     string
	// Search is the lowercased text the filter script matches against.
	Search string
	Count  int
}

type bookView struct {
	Title     string
	Subtitle  string
	Series    string
	Date      string
	Duration  string
	Search    string
	Authors   []linkView
	Narrators []linkView
	Genres    []linkView
	Tags      []linkView
}

// pageView is the data every template receives.
type pageView struct {
	Title       string
	Heading     string
	Blurb       string
	Root        string
	FeedHref    string
	SearchLabel string
	Search      bool
	Nav         []linkView
	Sections    []entryView
	Entries     []entryView
	Books       []bookView
}

// buildHTML renders the browsable tree that sits alongside the feeds: a home
// page, the searchable list of every book, an index per section, and a page for
// each author, narrator, genre and tag linking to its feed.
func buildHTML(content *Content) ([]htmlPage, error) {
	tpls, err := parseHTMLTemplates()
	if err != nil {
		return nil, err
	}
	shown := shownSections(content)

	pages := []htmlPage{}
	add := func(tpl *template.Template, relPath string, view *pageView) error {
		var buf bytes.Buffer
		if err := tpl.ExecuteTemplate(&buf, "layout", view); err != nil {
			return fmt.Errorf("rendering %s: %w", relPath, err)
		}
		pages = append(pages, htmlPage{Path: relPath, Contents: buf.Bytes()})
		return nil
	}

	if err := add(tpls.home, indexName, homeView(content, shown)); err != nil {
		return nil, err
	}
	if err := add(tpls.books, path.Join(booksDir, indexName), allBooksView(content, shown)); err != nil {
		return nil, err
	}

	for _, section := range shown {
		if err := add(
			tpls.entries,
			path.Join(section.Dir, indexName),
			sectionView(shown, section),
		); err != nil {
			return nil, err
		}
		for index := range section.Names {
			name := &section.Names[index]
			if err := add(
				tpls.books,
				path.Join(section.Dir, name.Slug, indexName),
				nameView(content, shown, name),
			); err != nil {
				return nil, err
			}
		}
	}

	return pages, nil
}

// htmlTemplates holds one parsed set per page shape. Each page template defines
// "content", so they cannot share a set, and there are enough author pages that
// reparsing per page would be wasteful.
type htmlTemplates struct {
	home    *template.Template
	entries *template.Template
	books   *template.Template
}

func parseHTMLTemplates() (*htmlTemplates, error) {
	var parsed htmlTemplates
	for _, target := range []struct {
		into **template.Template
		name string
	}{
		{&parsed.home, "home.gohtml"},
		{&parsed.entries, "entries.gohtml"},
		{&parsed.books, "books.gohtml"},
	} {
		tpl, err := template.ParseFS(templates.Templates, "layout.gohtml", target.name)
		if err != nil {
			return nil, err
		}
		*target.into = tpl
	}
	return &parsed, nil
}

func homeView(content *Content, shown []shownSection) *pageView {
	view := &pageView{
		Title:    allTitle,
		Heading:  allTitle + " 🎧📚",
		Blurb:    allSummary,
		Root:     rootFromRoot,
		FeedHref: feedHref(rootFromRoot, &content.Main),
		Nav:      navLinks(shown, rootFromRoot),
		Sections: []entryView{{
			Name:  booksTitle,
			Href:  dirHref(rootFromRoot, booksDir),
			Note:  "Search the whole library.",
			Count: len(content.Books),
		}},
	}
	for _, section := range shown {
		view.Sections = append(view.Sections, entryView{
			Name:  section.Title,
			Href:  dirHref(rootFromRoot, section.Dir),
			Note:  section.Blurb,
			Count: len(section.Names),
		})
	}
	return view
}

func allBooksView(content *Content, shown []shownSection) *pageView {
	return &pageView{
		Title:       pageTitle(booksTitle),
		Heading:     booksTitle,
		Blurb:       fmt.Sprintf("%d books in the library, most recently published first.", len(content.Books)),
		Root:        rootFromDir,
		Nav:         navLinks(shown, rootFromDir),
		Search:      true,
		SearchLabel: "Search by title, author, narrator, series, genre or tag",
		Books:       bookViews(content, rootFromDir, content.Books),
	}
}

func sectionView(shown []shownSection, section shownSection) *pageView {
	view := &pageView{
		Title:       pageTitle(section.Title),
		Heading:     section.Title,
		Blurb:       section.Blurb,
		Root:        rootFromDir,
		Nav:         navLinks(shown, rootFromDir),
		Search:      true,
		SearchLabel: "Search " + strings.ToLower(section.Title),
		Entries:     make([]entryView, 0, len(section.Names)),
	}
	for index := range section.Names {
		name := &section.Names[index]
		view.Entries = append(view.Entries, entryView{
			Name:     name.Title,
			Href:     dirHref(rootFromDir, section.Dir, name.Slug),
			FeedHref: feedHref(rootFromDir, &name.Feed),
			Search:   strings.ToLower(name.Title),
			Count:    len(name.Books),
		})
	}
	return view
}

func nameView(content *Content, shown []shownSection, name *Name) *pageView {
	return &pageView{
		Title:    pageTitle(name.Title),
		Heading:  name.Title,
		Blurb:    name.Description,
		Root:     rootFromSubDir,
		FeedHref: feedHref(rootFromSubDir, &name.Feed),
		Nav:      navLinks(shown, rootFromSubDir),
		Books:    bookViews(content, rootFromSubDir, name.Books),
	}
}

func navLinks(shown []shownSection, root string) []linkView {
	nav := make([]linkView, 0, len(shown)+1)
	nav = append(nav, linkView{Name: booksTitle, Href: dirHref(root, booksDir)})
	for _, section := range shown {
		nav = append(nav, linkView{Name: section.Title, Href: dirHref(root, section.Dir)})
	}
	return nav
}

// feedHref is empty for a name with no feed, which the templates leave out.
func feedHref(root string, feed *Feed) string {
	if feed.Link == "" {
		return ""
	}
	return hrefFor(root, feed.Link)
}

func bookViews(content *Content, root string, books []audiobooks.Audiobook) []bookView {
	sorted := newestFirst(books)
	views := make([]bookView, 0, len(sorted))
	for index := range sorted {
		book := &sorted[index]
		view := bookView{
			Title:     book.Title,
			Subtitle:  book.Subtitle,
			Date:      formatReleaseDate(book),
			Duration:  formatDuration(book.Duration),
			Authors:   links(root, &content.Authors, book.Authors),
			Narrators: links(root, &content.Narrators, book.Narrators),
			Genres:    links(root, &content.Genres, genreNames(book.Genres)),
			Tags:      links(root, &content.Tags, book.Tags),
		}
		if book.Series != nil {
			view.Series = book.Series.Title
			if sequence := book.Series.Sequence; sequence != nil {
				view.Series = fmt.Sprintf("%s %s %s", book.Series.Title, strings.ToLower(sequence.Noun()), sequence)
			}
		}
		view.Search = searchText(book, view.Series, view.Date)
		views = append(views, view)
	}
	return views
}

// newestFirst copies the books into reverse publication order, so that a page
// leads with what came out most recently.
//
// The copy is what makes this safe: these slices are the ones the feeds are
// rendered from, and a feed's item order is oldest first — reordering it in place
// would rewrite every feed in the site.
func newestFirst(books []audiobooks.Audiobook) []audiobooks.Audiobook {
	sorted := slices.Clone(books)
	// Stable, so books sharing a release date keep the order the scan sorted them
	// into, which is by title. A book with no date has the zero time and therefore
	// lands at the end.
	slices.SortStableFunc(sorted, func(a, b audiobooks.Audiobook) int {
		return releaseTime(&b).Compare(releaseTime(&a))
	})
	return sorted
}

func releaseTime(book *audiobooks.Audiobook) time.Time {
	if book.ReleaseDate == nil {
		return time.Time{}
	}
	return book.ReleaseDate.AsTime(time.UTC)
}

// formatReleaseDate renders the date a feed gives the book as its pubDate. The
// feed offsets it by eight hours so that clients in western timezones show the
// right day; the day itself is what appears here.
func formatReleaseDate(book *audiobooks.Audiobook) string {
	if book.ReleaseDate == nil {
		return ""
	}
	return releaseTime(book).Format("2 Jan 2006")
}

// links turns names into links to their pages. A name with no page — one that
// is blank once whitespace is removed, so the plan gave it no feed — is still
// shown, as plain text.
func links(root string, section *Section, names []string) []linkView {
	links := make([]linkView, 0, len(names))
	for _, name := range names {
		link := linkView{Name: name}
		if slug, ok := section.Slugs[name]; ok {
			link.Href = dirHref(root, section.Dir, slug)
		}
		links = append(links, link)
	}
	return links
}

// searchText is the haystack the filter script matches a query against: every
// word about the book that a reader might type, lowercased.
func searchText(book *audiobooks.Audiobook, series, date string) string {
	fields := make([]string, 0, 4+len(book.Authors)+len(book.Narrators)+len(book.Genres)+len(book.Tags))
	fields = append(fields, book.Title, book.Subtitle, series, date)
	fields = append(fields, book.Authors...)
	fields = append(fields, book.Narrators...)
	fields = append(fields, genreNames(book.Genres)...)
	fields = append(fields, book.Tags...)

	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			parts = append(parts, field)
		}
	}
	return strings.ToLower(strings.Join(parts, " "))
}

// pageTitle names a page in the browser's tab and history, where it needs to say
// which site it belongs to.
func pageTitle(heading string) string {
	return heading + " · " + allTitle
}

func genreNames(genres []audiobooks.Genre) []string {
	names := make([]string, 0, len(genres))
	for _, genre := range genres {
		if name := genre.String(); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// formatDuration renders a listening time. A book shorter than a minute is a
// placeholder or a failed scan rather than a real duration, so it says nothing.
func formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return ""
	}
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%dh %02dm", hours, minutes)
}

// hrefFor builds a relative link from a page whose prefix back to the root is
// root. The path is escaped by url.URL rather than pasted into a string, for the
// same reason feed.mediaURL does it: a name containing ?, # or % must stay data
// instead of becoming part of the URL grammar.
func hrefFor(root string, segments ...string) string {
	link := url.URL{Path: path.Join(segments...)}
	return root + link.EscapedPath()
}

// dirHref links to a directory, which the web server serves the index.html of.
func dirHref(root string, segments ...string) string {
	return hrefFor(root, segments...) + "/"
}
