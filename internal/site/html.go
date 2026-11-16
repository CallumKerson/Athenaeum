package site

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"path"
	"slices"
	"strconv"
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
	// Current marks the navigation link to the section the page is in.
	Current bool
}

type entryView struct {
	Name     string
	Href     string
	FeedHref string
	// Search is the lowercased text the filter script matches against.
	Search string
	Count  int
	// Hours is the total listening time of the books, rounded to the hour.
	Hours string
}

type bookView struct {
	Title    string
	Subtitle string
	Series   string
	Date     string
	Duration string
	// Href links to the book's own page.
	Href string
	// Cover links to the book's thumbnail, or is "" when it has none.
	Cover string
	// SeriesHref links to the page of the book's series, if it has one.
	SeriesHref string
	Search     string
	Authors    []linkView
	Narrators  []linkView
	Genres     []linkView
	Tags       []linkView
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
	Entries     []entryView
	Books       []bookView
	// Book is set only on a book's own page.
	Book *bookPage
	// Home is set only on the home page.
	Home *homePage
	// Series is set only on a series page: its books in reading order.
	Series []seriesEntry
}

// bookPage is everything a book's own page shows beyond its list entry.
type bookPage struct {
	bookView
	Paragraphs []string
	AudioURL   string
	MIMEType   string
	// Feeds are the authors' feeds, for subscribing from the book.
	Feeds []linkView
	// SeriesBooks lists the whole series in reading order, this book included.
	// It is not called Series because the embedded bookView already has a Series:
	// the line naming the series and this book's place in it.
	SeriesBooks []seriesEntry
	SeriesTitle string
}

type seriesEntry struct {
	bookView
	Position string
	Current  bool
}

// library is what every page is built from: the plan, and what the build
// learned while publishing it.
type library struct {
	*Content
	// covers holds the published thumbnail of each ImagePath that has one.
	covers map[string]string
	// audioURL is the public URL a book's audio is played from.
	audioURL func(*audiobooks.Audiobook) string
	shown    []shownSection
	// authorPages finds an author by the slug of their page.
	authorPages map[string]*Name
	// seriesOf finds the series a book, by its Path, belongs to.
	seriesOf map[string]*Series
}

func newLibrary(content *Content, covers map[string]string, audioURL func(*audiobooks.Audiobook) string) *library {
	lib := &library{
		Content:     content,
		covers:      covers,
		audioURL:    audioURL,
		shown:       shownSections(content),
		authorPages: make(map[string]*Name, len(content.Authors.Names)),
	}
	for index := range content.Authors.Names {
		lib.authorPages[content.Authors.Names[index].Slug] = &content.Authors.Names[index]
	}
	lib.seriesOf = map[string]*Series{}
	for index := range content.Series {
		series := &content.Series[index]
		for member := range series.Books {
			lib.seriesOf[series.Books[member].Path] = series
		}
	}
	return lib
}

// buildHTML renders the browsable tree that sits alongside the feeds: a home
// page, the searchable list of every book, a page per book, an index per
// section, and a page for each author, narrator, genre and tag linking to its
// feed.
func buildHTML(lib *library) ([]htmlPage, error) {
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

	if err := add(tpls.home, indexName, homeView(lib)); err != nil {
		return nil, err
	}
	if err := add(tpls.books, path.Join(booksDir, indexName), allBooksView(lib)); err != nil {
		return nil, err
	}
	for index := range lib.Books {
		book := &lib.Books[index]
		relPath := path.Join(booksDir, lib.BookSlugs[book.Path], indexName)
		if err := add(tpls.book, relPath, bookPageView(lib, book)); err != nil {
			return nil, err
		}
	}
	if len(lib.Series) > 0 {
		if err := add(tpls.entries, path.Join(seriesDir, indexName), seriesIndexView(lib)); err != nil {
			return nil, err
		}
	}
	for index := range lib.Series {
		series := &lib.Series[index]
		if err := add(
			tpls.series,
			path.Join(seriesDir, series.Slug, indexName),
			seriesPageView(lib, series),
		); err != nil {
			return nil, err
		}
	}

	for _, section := range lib.shown {
		if err := add(
			tpls.entries,
			path.Join(section.Dir, indexName),
			sectionView(lib, section),
		); err != nil {
			return nil, err
		}
		for index := range section.Names {
			name := &section.Names[index]
			if err := add(
				tpls.books,
				path.Join(section.Dir, name.Slug, indexName),
				nameView(lib, section, name),
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
	book    *template.Template
	books   *template.Template
	series  *template.Template
}

func parseHTMLTemplates() (*htmlTemplates, error) {
	var parsed htmlTemplates
	for _, target := range []struct {
		into **template.Template
		name string
	}{
		{&parsed.home, "home.gohtml"},
		{&parsed.entries, "entries.gohtml"},
		{&parsed.book, "book.gohtml"},
		{&parsed.books, "books.gohtml"},
		{&parsed.series, "series.gohtml"},
	} {
		tpl, err := template.ParseFS(templates.Templates, "layout.gohtml", target.name)
		if err != nil {
			return nil, err
		}
		*target.into = tpl
	}
	return &parsed, nil
}

// How many covers the home page mosaic shows, and how many books its shelf of
// latest releases holds.
const (
	mosaicSize = 15
	shelfSize  = 16
)

// homePage is what the home page shows beyond the navigation.
type homePage struct {
	Books int
	Hours string
	Days  int
	// BooksHref and GenresHref link to the full lists the home page samples.
	BooksHref  string
	GenresHref string
	// Mosaic is a spread of covers from across the library.
	Mosaic []bookView
	Latest []bookView
	Genres []genreTile
}

// genreTile is a genre on the home page, shown by a few of its covers.
type genreTile struct {
	entryView
	Covers []string
}

// homeView describes the home page. It shows the books in the main feed, which
// its subscribe button is for, so genres kept out of that feed stay off the
// front page too.
func homeView(lib *library) *pageView {
	root := rootFromRoot
	featured := lib.Main.Books
	listening := totalDuration(featured)
	home := &homePage{
		Books:      len(featured),
		Hours:      totalHours(featured),
		Days:       int(math.Round(listening.Hours() / 24)),
		Mosaic:     mosaic(lib, root, featured),
		BooksHref:  dirHref(root, booksDir),
		GenresHref: dirHref(root, lib.Genres.Dir),
	}
	latest := newestFirst(featured)
	for index := range latest[:min(shelfSize, len(latest))] {
		home.Latest = append(home.Latest, describeBook(lib, root, &latest[index]))
	}
	for index := range lib.Genres.Names {
		genre := &lib.Genres.Names[index]
		books := inMainFeed(lib, genre.Books)
		if len(books) == 0 {
			continue
		}
		tile := genreTile{entryView: entryView{
			Name:  genre.Title,
			Href:  dirHref(root, lib.Genres.Dir, genre.Slug),
			Count: len(books),
			Hours: totalHours(books),
		}}
		newest := newestFirst(books)
		for index := range newest {
			if cover, ok := lib.covers[newest[index].ImagePath]; ok && len(tile.Covers) < 3 {
				tile.Covers = append(tile.Covers, hrefFor(root, cover))
			}
		}
		home.Genres = append(home.Genres, tile)
	}
	slices.SortStableFunc(home.Genres, func(a, b genreTile) int { return b.Count - a.Count })

	return &pageView{
		Title:    allTitle,
		Heading:  allTitle,
		Blurb:    allSummary,
		Root:     root,
		FeedHref: feedHref(root, &lib.Main),
		Nav:      navLinks(lib, root, ""),
		Home:     home,
	}
}

// mosaic picks covers spread evenly across the library in path order, so the
// home page shows its range rather than only the latest books, and the choice
// only changes when the library does.
func mosaic(lib *library, root string, books []audiobooks.Audiobook) []bookView {
	var withCovers []*audiobooks.Audiobook
	for index := range books {
		if _, ok := lib.covers[books[index].ImagePath]; ok {
			withCovers = append(withCovers, &books[index])
		}
	}
	slices.SortFunc(withCovers, func(a, b *audiobooks.Audiobook) int { return strings.Compare(a.Path, b.Path) })

	count := min(mosaicSize, len(withCovers))
	views := make([]bookView, 0, count)
	for index := range count {
		views = append(views, describeBook(lib, root, withCovers[index*len(withCovers)/count]))
	}
	return views
}

// inMainFeed keeps the books that the main feed includes.
func inMainFeed(lib *library, books []audiobooks.Audiobook) []audiobooks.Audiobook {
	included := make(map[string]bool, len(lib.Main.Books))
	for index := range lib.Main.Books {
		included[lib.Main.Books[index].Path] = true
	}
	var kept []audiobooks.Audiobook
	for index := range books {
		if included[books[index].Path] {
			kept = append(kept, books[index])
		}
	}
	return kept
}

func totalDuration(books []audiobooks.Audiobook) time.Duration {
	var total time.Duration
	for index := range books {
		total += books[index].Duration
	}
	return total
}

// totalHours renders the listening time of some books in whole hours, with a
// thousands separator: "11,784".
func totalHours(books []audiobooks.Audiobook) string {
	hours := strconv.Itoa(int(math.Round(totalDuration(books).Hours())))
	var builder strings.Builder
	for index, digit := range hours {
		if index > 0 && (len(hours)-index)%3 == 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune(digit)
	}
	return builder.String()
}

func allBooksView(lib *library) *pageView {
	return &pageView{
		Title:       pageTitle(booksTitle),
		Heading:     booksTitle,
		Blurb:       fmt.Sprintf("%d books in the library, most recently published first.", len(lib.Books)),
		Root:        rootFromDir,
		Nav:         navLinks(lib, rootFromDir, booksDir),
		Search:      true,
		SearchLabel: "Search by title, author, narrator, series, genre or tag",
		Books:       bookViews(lib, rootFromDir, lib.Books),
	}
}

func sectionView(lib *library, section shownSection) *pageView {
	view := &pageView{
		Title:       pageTitle(section.Title),
		Heading:     section.Title,
		Blurb:       section.Blurb,
		Root:        rootFromDir,
		Nav:         navLinks(lib, rootFromDir, section.Dir),
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
			Hours:    totalHours(name.Books),
		})
	}
	return view
}

func nameView(lib *library, section shownSection, name *Name) *pageView {
	return &pageView{
		Title:    pageTitle(name.Title),
		Heading:  name.Title,
		Blurb:    name.Description,
		Root:     rootFromSubDir,
		FeedHref: feedHref(rootFromSubDir, &name.Feed),
		Nav:      navLinks(lib, rootFromSubDir, section.Dir),
		Books:    bookViews(lib, rootFromSubDir, name.Books),
	}
}

// navLinks links every top-level section from a page whose prefix back to the
// site root is root, marking current, the directory of the page's own section,
// or "" for the home page.
func navLinks(lib *library, root, current string) []linkView {
	nav := make([]linkView, 0, len(lib.shown)+2)
	nav = append(nav, linkView{Name: booksTitle, Href: dirHref(root, booksDir), Current: current == booksDir})
	if len(lib.Series) > 0 {
		nav = append(nav, linkView{Name: seriesTitle, Href: dirHref(root, seriesDir), Current: current == seriesDir})
	}
	for _, section := range lib.shown {
		nav = append(nav, linkView{
			Name:    section.Title,
			Href:    dirHref(root, section.Dir),
			Current: current == section.Dir,
		})
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

// bookViews describes books for a page whose prefix back to the site root is
// root.
func bookViews(lib *library, root string, books []audiobooks.Audiobook) []bookView {
	sorted := newestFirst(books)
	views := make([]bookView, 0, len(sorted))
	for index := range sorted {
		views = append(views, describeBook(lib, root, &sorted[index]))
	}
	return views
}

func describeBook(lib *library, root string, book *audiobooks.Audiobook) bookView {
	view := bookView{
		Title:     book.Title,
		Subtitle:  book.Subtitle,
		Date:      formatReleaseDate(book),
		Duration:  formatDuration(book.Duration),
		Href:      dirHref(root, booksDir, lib.BookSlugs[book.Path]),
		Authors:   links(root, &lib.Authors, book.Authors),
		Narrators: links(root, &lib.Narrators, book.Narrators),
		Genres:    links(root, &lib.Genres, genreNames(book.Genres)),
		Tags:      links(root, &lib.Tags, book.Tags),
	}
	if cover, ok := lib.covers[book.ImagePath]; ok {
		view.Cover = hrefFor(root, cover)
	}
	if series := lib.seriesOf[book.Path]; series != nil {
		view.SeriesHref = dirHref(root, seriesDir, series.Slug)
	}
	if book.Series != nil {
		view.Series = book.Series.Title
		if sequence := book.Series.Sequence; sequence != nil {
			view.Series = fmt.Sprintf("%s %s %s", book.Series.Title, strings.ToLower(sequence.Noun()), sequence)
		}
	}
	view.Search = searchText(book, view.Series, view.Date)
	return view
}

// moreByLimit caps the "more by" list, which for a prolific author would
// otherwise repeat their whole page.
const moreByLimit = 12

// bookPageView describes a book's own page. series holds every series in
// reading order, from seriesInOrder.
func bookPageView(lib *library, book *audiobooks.Audiobook) *pageView {
	root := rootFromSubDir
	page := &bookPage{
		bookView:   describeBook(lib, root, book),
		Paragraphs: book.Description.Paragraphs(),
		MIMEType:   book.MIMEType,
	}
	if lib.audioURL != nil {
		page.AudioURL = lib.audioURL(book)
	}
	page.Feeds = authorFeeds(lib, root, book)
	series := lib.seriesOf[book.Path]
	if series != nil && len(series.Books) > 1 {
		page.SeriesTitle = series.Title
		page.SeriesBooks = seriesEntries(lib, root, series, book.Path)
	}

	var others []audiobooks.Audiobook
	for index := range lib.Books {
		other := &lib.Books[index]
		inSameSeries := series != nil && lib.seriesOf[other.Path] == series
		if other.Path != book.Path && !inSameSeries && sameAuthor(lib.Authors.Slugs, book, other) {
			others = append(others, *other)
		}
	}
	// The authors' other books outside this series fill the page's book list.
	moreBy := bookViews(lib, root, others)
	if len(moreBy) > moreByLimit {
		moreBy = moreBy[:moreByLimit]
	}

	// No Heading: the page puts the title beside the cover, under the series.
	return &pageView{
		Title: pageTitle(book.Title),
		Root:  root,
		Nav:   navLinks(lib, root, booksDir),
		Book:  page,
		Books: moreBy,
	}
}

func authorFeeds(lib *library, root string, book *audiobooks.Audiobook) []linkView {
	var feeds []linkView
	for _, author := range book.Authors {
		if name := lib.authorPages[lib.Authors.Slugs[author]]; name != nil && name.Link != "" {
			feeds = append(feeds, linkView{Name: name.Title, Href: hrefFor(root, name.Link)})
		}
	}
	return feeds
}

func position(book *audiobooks.Audiobook) string {
	if book.Series == nil || book.Series.Sequence == nil {
		return ""
	}
	return book.Series.Sequence.String()
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
