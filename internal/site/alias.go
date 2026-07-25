package site

import (
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// nameAliases returns every directory name a feed for name should be written
// under.
//
// The server this replaces matched authors, narrators and tags case- and
// whitespace-insensitively, so `/authors/Adrian%20Tchaikovsky/`,
// `/authors/adrian tchaikovsky/` and `/authors/adriantchaikovsky/` all resolved
// to the same feed. A static tree has no matching logic, so each accepted
// spelling gets its own file.
func nameAliases(name string) []string {
	return dedupe([]string{name, strings.ToLower(name), normalise(name)})
}

// genreAliases returns the directory names for a genre feed: the canonical
// display name plus every synonym ParseGenre accepts, so `sci-fi`, `scifi`,
// `sciencefiction` and `Science Fiction` all keep working.
func genreAliases(genre audiobooks.Genre) []string {
	return dedupe(append([]string{genre.String()}, genre.Aliases()...))
}

// normalise mirrors the server's filter comparison: lowercase, whitespace
// removed.
func normalise(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, str)
}

// safeSegment reports whether a name can be used as a single path segment.
// Anything containing a separator, or that would escape the output root, is
// rejected rather than sanitised — a book's metadata should not be able to
// direct writes outside the site directory.
func safeSegment(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	if strings.ContainsRune(name, '/') || strings.ContainsRune(name, filepath.Separator) {
		return false
	}
	return !strings.ContainsRune(name, 0)
}

func dedupe(values []string) []string {
	seen := map[string]bool{}
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		unique = append(unique, value)
	}
	slices.Sort(unique)
	return unique
}
