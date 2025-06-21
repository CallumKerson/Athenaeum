package description

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat_String(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{Plain, "Plain"},
		{Markdown, "Markdown"},
		{HTML, "HTML"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.format.String())
		})
	}
}

func TestParseFormat_ValidInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		// Exact matches
		{"plain", Plain},
		{"markdown", Markdown},
		{"html", HTML},

		// Case variations
		{"PLAIN", Plain},
		{"Plain", Plain},
		{"MARKDOWN", Markdown},
		{"Markdown", Markdown},
		{"HTML", HTML},
		{"Html", HTML},

		// With whitespace
		{" plain ", Plain},
		{"\tmarkdown\n", Markdown},
		{" html ", HTML},

		// Alternative spellings/aliases
		{"md", Markdown},
		{"MD", Markdown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseFormat(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseFormat_InvalidInputs(t *testing.T) {
	tests := []string{
		"",
		"unknown",
		"invalid",
		"xml",
		"json",
		"yaml",
		"123",
		"plain-text",
		"rich-text",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseFormat(input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not a valid textFormat")
		})
	}
}

func TestFormat_MarshalText(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{Plain, "Plain"},
		{Markdown, "Markdown"},
		{HTML, "HTML"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			data, err := tt.format.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestFormat_UnmarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"plain", Plain},
		{"markdown", Markdown},
		{"html", HTML},
		{"md", Markdown},         // alias
		{"PLAIN", Plain},         // uppercase
		{"HTML", HTML},           // uppercase
		{" markdown ", Markdown}, // with whitespace
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var format Format
			err := format.UnmarshalText([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, format)
		})
	}
}

func TestFormat_UnmarshalText_Error(t *testing.T) {
	var format Format
	err := format.UnmarshalText([]byte("invalid-format"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid textFormat")
}

func TestFormat_RoundTrip(t *testing.T) {
	// Test that marshalling and unmarshalling preserves the format
	allFormats := []Format{Plain, Markdown, HTML}

	for _, original := range allFormats {
		t.Run(original.String(), func(t *testing.T) {
			// Marshal to text
			data, err := original.MarshalText()
			require.NoError(t, err)

			// Unmarshal back
			var unmarshaled Format
			err = unmarshaled.UnmarshalText(data)
			require.NoError(t, err)

			// Should be equal
			assert.Equal(t, original, unmarshaled)
		})
	}
}

func TestFormat_ZeroValue(t *testing.T) {
	var format Format

	// Zero value should be Undefined (0)
	assert.Equal(t, Undefined, format)
	assert.Equal(t, "", format.String())

	// Zero value should marshal correctly
	data, err := format.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(data))
}
