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
		{name: "match author ignoring case", filter: AuthorFilter("Max gladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match author ignoring whitespace", filter: AuthorFilter("MaxGladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match author ignoring case and whitespace", filter: AuthorFilter("maxgladstone"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match author", filter: AuthorFilter("Max Gladstone"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match genre", filter: GenreFilter(audiobooks.SciFi), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring case", filter: NarratorFilter("Emily woo zeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring whitespace", filter: NarratorFilter("EmilyWooZeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match narrator ignoring case and whitespace", filter: NarratorFilter("emilywoozeller"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match narrator", filter: NarratorFilter("Emily Woo Zeller"), book: testbooks.Audiobooks[1], expectedMatch: false},
		{name: "match tag", filter: TagFilter("Hugo Awards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring case", filter: TagFilter("Hugo awards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring whitespace", filter: TagFilter("HugoAwards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "match tag ignoring case and whitespace", filter: TagFilter("hugoawards"), book: testbooks.Audiobooks[0], expectedMatch: true},
		{name: "not match tag", filter: TagFilter("Nebula Awards"), book: testbooks.Audiobooks[1], expectedMatch: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedMatch, testCase.filter(&testCase.book))
		})
	}
}
