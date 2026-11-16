package site

import (
	"cmp"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	statsDir   = "stats"
	statsTitle = "Stats"
	statsBlurb = "The whole library in numbers, excluded genres included."
	// topNarrators is how many narrators the stats page ranks.
	topNarrators = 10
)

// statsPage is what the stats page shows.
type statsPage struct {
	Totals    []statTotal
	Genres    []statBar
	Narrators []statBar
	Decades   []statBar
	Longest   *bookView
}

// statTotal is one headline number.
type statTotal struct {
	Value string
	Label string
}

// statBar is one row of a bar chart. Share is the bar's length as a whole
// percentage of the longest bar, so it can go straight into the page's CSS.
type statBar struct {
	Label string
	Href  string
	Value string
	Note  string
	Share int
}

func statsView(lib *library) *pageView {
	root := rootFromDir
	stats := &statsPage{
		Totals: []statTotal{
			{strconv.Itoa(len(lib.Books)), plural(len(lib.Books), "book")},
			{totalHours(lib.Books), "hours of listening"},
			{strconv.Itoa(len(lib.Authors.Names)), plural(len(lib.Authors.Names), "author")},
			{strconv.Itoa(len(lib.Narrators.Names)), plural(len(lib.Narrators.Names), "narrator")},
			{strconv.Itoa(len(lib.Series)), "series"},
		},
		Genres:    namesByHours(root, &lib.Genres, len(lib.Genres.Names)),
		Narrators: namesByHours(root, &lib.Narrators, topNarrators),
		Decades:   decades(lib.Books),
	}
	if longest := longestBook(lib.Books); longest != nil {
		view := describeBook(lib, root, longest)
		stats.Longest = &view
	}
	return &pageView{
		Title:   pageTitle(statsTitle),
		Heading: statsTitle,
		Blurb:   statsBlurb,
		Root:    root,
		Nav:     navLinks(lib, root, statsDir),
		Stats:   stats,
	}
}

// namesByHours ranks a section's names by their total listening time, keeping
// the first limit with any books. Ties keep the section's own order.
func namesByHours(root string, section *Section, limit int) []statBar {
	type ranked struct {
		name    *Name
		minutes float64
	}
	var names []ranked
	for index := range section.Names {
		name := &section.Names[index]
		if len(name.Books) > 0 {
			names = append(names, ranked{name, totalDuration(name.Books).Minutes()})
		}
	}
	slices.SortStableFunc(names, func(a, b ranked) int { return cmp.Compare(b.minutes, a.minutes) })
	names = names[:min(limit, len(names))]

	bars := make([]statBar, 0, len(names))
	for _, entry := range names {
		bars = append(bars, statBar{
			Label: entry.name.Title,
			Href:  dirHref(root, section.Dir, entry.name.Slug),
			Value: hoursLabel(totalDuration(entry.name.Books)),
			Note:  strconv.Itoa(len(entry.name.Books)) + " " + plural(len(entry.name.Books), "book"),
			Share: share(entry.minutes, names[0].minutes),
		})
	}
	return bars
}

// decades counts the books first published in each decade, oldest first,
// leaving out books with no release date.
func decades(books []audiobooks.Audiobook) []statBar {
	counts := map[int]int{}
	for index := range books {
		if released := releaseTime(&books[index]); !released.IsZero() {
			counts[released.Year()/10*10]++
		}
	}
	starts := make([]int, 0, len(counts))
	most := 0
	for start, count := range counts {
		starts = append(starts, start)
		most = max(most, count)
	}
	slices.Sort(starts)

	bars := make([]statBar, 0, len(starts))
	for _, start := range starts {
		bars = append(bars, statBar{
			Label: strconv.Itoa(start) + "s",
			Value: strconv.Itoa(counts[start]) + " " + plural(counts[start], "book"),
			Share: share(float64(counts[start]), float64(most)),
		})
	}
	return bars
}

// longestBook finds the book with the most listening time, taking the first
// in path order on a tie so the choice is stable.
func longestBook(books []audiobooks.Audiobook) *audiobooks.Audiobook {
	var longest *audiobooks.Audiobook
	for index := range books {
		book := &books[index]
		if longest == nil || book.Duration > longest.Duration ||
			(book.Duration == longest.Duration && book.Path < longest.Path) {
			longest = book
		}
	}
	if longest == nil || longest.Duration <= time.Duration(0) {
		return nil
	}
	return longest
}

// hoursLabel says how long something is in whole hours, or in minutes when
// that would round to none.
func hoursLabel(listening time.Duration) string {
	hours := int(math.Round(listening.Hours()))
	if hours == 0 {
		minutes := int(math.Round(listening.Minutes()))
		return strconv.Itoa(minutes) + " " + plural(minutes, "minute")
	}
	return totalHoursOf(listening) + " " + plural(hours, "hour")
}

// share is value as a whole percentage of most, at least 1 so that every bar
// with anything in it shows.
func share(value, most float64) int {
	if most <= 0 || value <= 0 {
		return 0
	}
	return max(1, int(math.Round(value/most*100)))
}
