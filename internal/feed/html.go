package feed

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

// summaryHTML builds the rich item body that clients show in place of the plain
// description: headings for the book, its author, series and narrator, followed
// by the description converted from whatever format it was written in.
func summaryHTML(book *audiobooks.Audiobook) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "<h1>%s</h1>", book.Title)
	if book.Subtitle != "" {
		// Reproduces the original: the title, not the subtitle, is repeated here.
		fmt.Fprintf(&builder, "<h4>%s</h4>", book.Title)
	}
	fmt.Fprintf(&builder, "<h2>By %s</h2>", book.GetAuthor())
	if book.Series != nil {
		fmt.Fprintf(&builder, "<h4>%s %s %s</h4>", book.Series.Title, book.Series.Sequence.Noun(), book.Series.Sequence)
	}
	if book.GetNarrator() != "" {
		fmt.Fprintf(&builder, "<h4>Narrated by %s</h4>", book.GetNarrator())
	}

	if book.Description != nil {
		switch book.Description.Format {
		case description.HTML:
			builder.WriteString(book.Description.Text)
		case description.Markdown:
			builder.Write(markdown.ToHTML([]byte(book.Description.Text), nil, nil))
		case description.Plain, description.Undefined:
			for line := range strings.SplitSeq(book.Description.Text, "\n") {
				fmt.Fprintf(&builder, "<p>%s</p>", line)
			}
		}
	}

	return builder.String()
}
