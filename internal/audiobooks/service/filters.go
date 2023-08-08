package service

import (
	"strings"
	"unicode"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type Filter func(a *audiobooks.Audiobook) bool

func AuthorFilter(name string) Filter {
	return func(a *audiobooks.Audiobook) bool {
		if a != nil && containsIgnoringCaseAndWhitespace(a.Authors, name) {
			return true
		}
		return false
	}
}

func GenreFilter(genre audiobooks.Genre) Filter {
	return func(a *audiobooks.Audiobook) bool {
		if a != nil && contains(a.Genres, genre) {
			return true
		}
		return false
	}
}

func NarratorFilter(name string) Filter {
	return func(a *audiobooks.Audiobook) bool {
		if a != nil && containsIgnoringCaseAndWhitespace(a.Narrators, name) {
			return true
		}
		return false
	}
}

func TagFilter(tag string) Filter {
	return func(a *audiobooks.Audiobook) bool {
		if a != nil && containsIgnoringCaseAndWhitespace(a.Tags, tag) {
			return true
		}
		return false
	}
}

func NotFilter(filter Filter) Filter {
	return func(a *audiobooks.Audiobook) bool {
		return !filter(a)
	}
}

func AndFilter(filters ...Filter) Filter {
	return func(a *audiobooks.Audiobook) bool {
		fulfilledFilters := 0
		for filterIndex := range filters {
			if filters[filterIndex](a) {
				fulfilledFilters++
			}
		}
		return fulfilledFilters == len(filters)
	}
}

func contains[K comparable](slice []K, item K) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func containsIgnoringCaseAndWhitespace(slice []string, item string) bool {
	for _, v := range slice {
		if normaliseString(v) == normaliseString(item) {
			return true
		}
	}
	return false
}

func normaliseString(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, str)
}
