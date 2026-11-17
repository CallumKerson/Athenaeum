package description

import (
	"html"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
)

var (
	// blockBreak matches the tags that end a block of text, so each becomes a
	// paragraph break once the markup is gone.
	blockBreak = regexp.MustCompile(`(?i)<\s*(?:br\s*/?|/\s*(?:p|div|li|h[1-6]|blockquote))\s*>`)
	anyTag     = regexp.MustCompile(`<[^>]*>`)
	// hidden matches elements whose content is code rather than text.
	hidden = regexp.MustCompile(`(?is)<script\b.*?</script\s*>|<style\b.*?</style\s*>`)
)

// Paragraphs returns the description as plain-text paragraphs, whatever format
// it was written in. Markup is dropped rather than passed on, so a page can show
// the text escaped like any other value instead of trusting the library's HTML.
func (d *Description) Paragraphs() []string {
	if d == nil {
		return nil
	}

	text := d.Text
	switch d.Format {
	case Markdown:
		text = textOf(string(markdown.ToHTML([]byte(text), nil, nil)))
	case HTML:
		text = textOf(text)
	case Plain, Undefined:
	}

	var paragraphs []string
	for line := range strings.SplitSeq(text, "\n") {
		if line = strings.Join(strings.Fields(line), " "); line != "" {
			paragraphs = append(paragraphs, line)
		}
	}
	return paragraphs
}

// textOf strips HTML down to its text, with a line break wherever a block ended.
func textOf(markup string) string {
	// Source newlines inside a paragraph are only wrapping, not breaks.
	markup = strings.ReplaceAll(markup, "\n", " ")
	markup = hidden.ReplaceAllString(markup, "")
	markup = blockBreak.ReplaceAllString(markup, "\n")
	return html.UnescapeString(anyTag.ReplaceAllString(markup, ""))
}
