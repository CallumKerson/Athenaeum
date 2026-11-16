package site

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestTotalHoursGroupsThousands(t *testing.T) {
	tests := []struct {
		hours float64
		want  string
	}{
		{0, "0"},
		{12.4, "12"},
		{999.6, "1,000"},
		{11784, "11,784"},
		{1234567, "1,234,567"},
	}
	for _, test := range tests {
		books := []audiobooks.Audiobook{{Duration: time.Duration(test.hours * float64(time.Hour))}}
		assert.Equal(t, test.want, totalHours(books), test.hours)
	}
}

func TestHomeViewCountsTheLibrary(t *testing.T) {
	first := book(
		"A Wizard of Earthsea",
		[]audiobooks.Genre{audiobooks.Fantasy},
		[]string{"Ursula K. Le Guin"},
		nil,
		nil,
	)
	first.Duration = 30 * time.Hour
	second := book("Dune", []audiobooks.Genre{audiobooks.SciFi}, []string{"Frank Herbert"}, nil, nil)
	second.Duration = 20 * time.Hour
	third := book("The Dispossessed", []audiobooks.Genre{audiobooks.SciFi}, []string{"Ursula K. Le Guin"}, nil, nil)
	third.Duration = 10 * time.Hour
	content := Plan([]audiobooks.Audiobook{first, second, third}, nil)

	home := homeView(newLibrary(content, nil, nil)).Home

	assert.Equal(t, 3, home.Books)
	assert.Equal(t, "60", home.Hours)
	assert.Equal(t, 3, home.Days)
	assert.Empty(t, home.Mosaic, "no covers, no mosaic")
	assert.Len(t, home.Latest, 3)
	if assert.Len(t, home.Genres, 2) {
		assert.Equal(t, "Science Fiction", home.Genres[0].Name, "the biggest genre comes first")
		assert.Equal(t, 2, home.Genres[0].Count)
		assert.Equal(t, "30", home.Genres[0].Hours)
	}
}

func TestHomePageRendersTheHero(t *testing.T) {
	root := t.TempDir()
	_, err := Build(root, testPlan(t), testRenderer(t), testOptions(), testLogger())
	require.NoError(t, err)

	page := readPage(t, root, ".")
	assert.Contains(t, page, `<section class="hero">`)
	assert.Contains(t, page, `audiobooks, ready for your podcast app</h1>`)
	assert.Contains(t, page, `<ul class="books shelf">`)
	assert.Contains(t, page, `<a href="./books/">See every book</a>`)
	assert.Contains(t, page, `<a href="./genre/">Every genre</a>`)
}
