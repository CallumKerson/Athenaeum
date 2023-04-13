package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestFilers(t *testing.T) {
	tests := []struct {
		name          string
		filter        Filter
		book          audiobooks.Audiobook
		expectedMatch bool
	}{
		{name: "match author", filter: AuthorFilter("Max Gladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match author", filter: AuthorFilter("Max Gladstone"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[1], expectedMatch: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedMatch, testCase.filter(&testCase.book))
		})
	}
}
