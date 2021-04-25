package htmlbuilder

import (
	"encoding/xml"
	"strings"
)

func H1(s string) string {
	var b strings.Builder
	b.WriteString(`<h1>`)
	_ = xml.EscapeText(&b, []byte(s))
	b.WriteString(`</h1>`)
	return b.String()
}

func H2(s string) string {
	var b strings.Builder
	b.WriteString(`<h2>`)
	_ = xml.EscapeText(&b, []byte(s))
	b.WriteString(`</h2>`)
	return b.String()
}

func H4(s string) string {
	var b strings.Builder
	b.WriteString(`<h4>`)
	_ = xml.EscapeText(&b, []byte(s))
	b.WriteString(`</h4>`)
	return b.String()
}

// P returns the string in paragraph tags. If there are newlines in the input,
// then the output will have each line in a separate paragraph tag
func P(s string) string {
	lines := strings.Split(s, "\n")
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(`<p>`)
		_ = xml.EscapeText(&b, []byte(line))
		b.WriteString(`</p>`)
	}
	return b.String()
}
