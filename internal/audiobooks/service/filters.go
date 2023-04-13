package service

import "github.com/CallumKerson/Athenaeum/pkg/audiobooks"

type Filter func(a *audiobooks.Audiobook) bool

func AuthorFilter(name string) Filter {
	return func(a *audiobooks.Audiobook) bool {
		if a != nil && contains(a.Authors, name) {
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
		if a != nil && contains(a.Narrators, name) {
			return true
		}
		return false
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
