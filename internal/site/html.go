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
	"unicode"

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

// browseSection is one way of browsing the library. Dir is shared with the feed
// tree, so that /authors/ and the feeds under /podcast/authors/ are obviously
// the same set of names.
type browseSection struct {
	Dir   string
	Title string
	Blurb string
}

var browseSections = []browseSection{
	{Dir: authorsDir, Title: "Authors", Blurb: "Everyone who wrote something in the library."},
	{Dir: narratorsDir, Title: "Narrators", Blurb: "Everyone who read something in the library."},
	{
		Dir:   genresDir,
		Title: "Genres",
		Blurb: "Every genre Athenaeum knows, including those nothing is filed under yet.",
	},
	{Dir: tagsDir, Title: "Tags", Blurb: "Whatever the library's own metadata tags books with."},
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
func buildHTML(content Content) ([]htmlPage, error) {
	builder := newHTMLBuilder(content)
	tpls, err := parseHTMLTemplates()
	if err != nil {
		return nil, err
	}

	pages := []htmlPage{}
	add := func(tpl *template.Template, relPath string, view *pageView) error {
		var buf bytes.Buffer
		if err := tpl.ExecuteTemplate(&buf, "layout", view); err != nil {
			return fmt.Errorf("rendering %s: %w", relPath, err)
		}
		pages = append(pages, htmlPage{Path: relPath, Contents: buf.Bytes()})
		return nil
	}

	if err := add(tpls.home, indexName, builder.home()); err != nil {
		return nil, err
	}
	if err := add(tpls.books, path.Join(booksDir, indexName), builder.allBooks()); err != nil {
		return nil, err
	}

	for _, section := range browseSections {
		if err := add(tpls.entries, path.Join(section.Dir, indexName), builder.sectionIndex(section)); err != nil {
			return nil, err
		}
		for index := range builder.sections[section.Dir] {
			feedPage := &builder.sections[section.Dir][index]
			relPath := path.Join(section.Dir, builder.slugOf(section.Dir, feedPage.Title), indexName)
			if err := add(tpls.books, relPath, builder.feedPage(feedPage)); err != nil {
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

type htmlBuilder struct {
	books []audiobooks.Audiobook
	main  *Page
	// Feed pages by section directory, in plan order — which is sorted by name.
	sections map[string][]Page
	// Slug per name, keyed by section directory and normalised name.
	slugs map[string]string
}

func newHTMLBuilder(content Content) *htmlBuilder {
	builder := &htmlBuilder{
		books:    content.Books,
		sections: map[string][]Page{},
		slugs:    map[string]string{},
	}
	for index := range content.Feeds {
		feedPage := content.Feeds[index]
		if feedPage.Section == "" {
			builder.main = &content.Feeds[index]
			continue
		}
		builder.sections[feedPage.Section] = append(builder.sections[feedPage.Section], feedPage)
	}
	for _, section := range browseSections {
		builder.assignSlugs(section.Dir)
	}
	return builder
}

// assignSlugs gives every name in a section a unique URL segment.
//
// Feed paths keep the library's own spelling, because existing subscriptions
// point at them. The HTML tree is new, so its directories can be plain
// lowercase-hyphen names, which need no escaping and cannot fold together on a
// case-insensitive filesystem.
func (b *htmlBuilder) assignSlugs(dir string) {
	taken := map[string]bool{}
	for _, feedPage := range b.sections[dir] {
		slug := slugify(feedPage.Title)
		// Distinct names can slugify identically — "Ann Leckie" and "Ann-Leckie" —
		// and must not end up sharing a page. Plan order is sorted, so which name
		// keeps the unsuffixed slug is stable between builds.
		unique := slug
		for suffix := 2; taken[unique]; suffix++ {
			unique = fmt.Sprintf("%s-%d", slug, suffix)
		}
		taken[unique] = true
		b.slugs[slugKey(dir, feedPage.Title)] = unique
	}
}

func (b *htmlBuilder) slugOf(dir, name string) string {
	return b.slugs[slugKey(dir, name)]
}

func (b *htmlBuilder) home() *pageView {
	view := &pageView{
		Title:   allTitle,
		Heading: allTitle + " 🎧📚",
		Blurb:   allSummary,
		Root:    rootFromRoot,
		Nav:     b.nav(rootFromRoot),
		Sections: []entryView{{
			Name:  booksTitle,
			Href:  dirHref(rootFromRoot, booksDir),
			Note:  "Search the whole library.",
			Count: len(b.books),
		}},
	}
	if b.main != nil {
		view.FeedHref = hrefFor(rootFromRoot, b.main.FeedPath)
	}
	for _, section := range browseSections {
		view.Sections = append(view.Sections, entryView{
			Name:  section.Title,
			Href:  dirHref(rootFromRoot, section.Dir),
			Note:  section.Blurb,
			Count: len(b.sections[section.Dir]),
		})
	}
	return view
}

func (b *htmlBuilder) allBooks() *pageView {
	return &pageView{
		Title:       pageTitle(booksTitle),
		Heading:     booksTitle,
		Blurb:       fmt.Sprintf("%d books in the library, most recently published first.", len(b.books)),
		Root:        rootFromDir,
		Nav:         b.nav(rootFromDir),
		Search:      true,
		SearchLabel: "Search by title, author, narrator, series, genre or tag",
		Books:       b.bookViews(rootFromDir, b.books),
	}
}

func (b *htmlBuilder) sectionIndex(section browseSection) *pageView {
	pages := b.sections[section.Dir]
	view := &pageView{
		Title:       pageTitle(section.Title),
		Heading:     section.Title,
		Blurb:       section.Blurb,
		Root:        rootFromDir,
		Nav:         b.nav(rootFromDir),
		Search:      true,
		SearchLabel: "Search " + strings.ToLower(section.Title),
		Entries:     make([]entryView, 0, len(pages)),
	}
	for index := range pages {
		entry := entryView{
			Name:   pages[index].Title,
			Href:   dirHref(rootFromDir, section.Dir, b.slugOf(section.Dir, pages[index].Title)),
			Search: strings.ToLower(pages[index].Title),
			Count:  len(pages[index].Books),
		}
		if pages[index].FeedPath != "" {
			entry.FeedHref = hrefFor(rootFromDir, pages[index].FeedPath)
		}
		view.Entries = append(view.Entries, entry)
	}
	return view
}

func (b *htmlBuilder) feedPage(feedPage *Page) *pageView {
	view := &pageView{
		Title:   pageTitle(feedPage.Title),
		Heading: feedPage.Title,
		Blurb:   feedPage.Description,
		Root:    rootFromSubDir,
		Nav:     b.nav(rootFromSubDir),
		Books:   b.bookViews(rootFromSubDir, feedPage.Books),
	}
	if feedPage.FeedPath != "" {
		view.FeedHref = hrefFor(rootFromSubDir, feedPage.FeedPath)
	}
	return view
}

func (b *htmlBuilder) nav(root string) []linkView {
	nav := make([]linkView, 0, len(browseSections)+1)
	nav = append(nav, linkView{Name: booksTitle, Href: dirHref(root, booksDir)})
	for _, section := range browseSections {
		nav = append(nav, linkView{Name: section.Title, Href: dirHref(root, section.Dir)})
	}
	return nav
}

func (b *htmlBuilder) bookViews(root string, books []audiobooks.Audiobook) []bookView {
	sorted := newestFirst(books)
	views := make([]bookView, 0, len(sorted))
	for index := range sorted {
		book := &sorted[index]
		view := bookView{
			Title:     book.Title,
			Subtitle:  book.Subtitle,
			Date:      formatReleaseDate(book),
			Duration:  formatDuration(book.Duration),
			Authors:   b.links(root, authorsDir, book.Authors),
			Narrators: b.links(root, narratorsDir, book.Narrators),
			Genres:    b.links(root, genresDir, genreNames(book.Genres)),
			Tags:      b.links(root, tagsDir, book.Tags),
		}
		if book.Series != nil {
			view.Series = fmt.Sprintf("%s %s %s",
				book.Series.Title, strings.ToLower(book.Series.Sequence.Noun()), book.Series.Sequence)
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

// links turns names into links to their pages. A name with no page — only
// possible if it is not in the plan at all — is still shown, as plain text.
func (b *htmlBuilder) links(root, dir string, names []string) []linkView {
	links := make([]linkView, 0, len(names))
	for _, name := range names {
		link := linkView{Name: name}
		if slug := b.slugOf(dir, name); slug != "" {
			link.Href = dirHref(root, dir, slug)
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

func slugKey(dir, name string) string {
	return dir + "\x00" + normalise(name)
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
