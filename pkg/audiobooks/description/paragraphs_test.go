package description

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParagraphs(t *testing.T) {
	tests := []struct {
		name        string
		description *Description
		want        []string
	}{
		{
			name:        "nil",
			description: nil,
			want:        nil,
		},
		{
			name:        "plain text splits on lines and drops blank ones",
			description: &Description{Text: "First line.\n\n  Second   line.  \n", Format: Plain},
			want:        []string{"First line.", "Second line."},
		},
		{
			name:        "undefined format reads as plain",
			description: &Description{Text: "One\nTwo"},
			want:        []string{"One", "Two"},
		},
		{
			name: "HTML keeps paragraphs and drops markup",
			description: &Description{
				Text:   "<p>Ged was called <i>Sparrowhawk</i>.</p>\n<p>He loosed a\nshadow &amp; more.<br/>Then</p>",
				Format: HTML,
			},
			want: []string{"Ged was called Sparrowhawk.", "He loosed a shadow & more.", "Then"},
		},
		{
			name: "HTML drops scripts and keeps escaped markup as text",
			description: &Description{
				Text:   `<p>Before</p><script>alert("x")</script><p>&lt;b&gt;After</p>`,
				Format: HTML,
			},
			want: []string{"Before", "<b>After"},
		},
		{
			name:        "markdown renders to paragraphs",
			description: &Description{Text: "A *bold* start.\n\nSecond paragraph.", Format: Markdown},
			want:        []string{"A bold start.", "Second paragraph."},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.description.Paragraphs())
		})
	}
}
