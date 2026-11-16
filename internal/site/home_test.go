package site

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
