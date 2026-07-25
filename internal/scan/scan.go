// Package scan walks a media root and turns the `.m4b` + `.toml` file pairs it
// finds into audiobooks.
package scan

import (
	"cmp"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/alfg/mp4"
	"github.com/pelletier/go-toml/v2"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// Every book in the library is an MPEG-4 audio file, so the type is constant
// rather than sniffed.
const mimeType = "audio/mp4a-latm"

// coverExtensions are probed in order; the first that exists wins.
var coverExtensions = []string{".jpg", ".jpeg", ".png"}

// Library walks mediaRoot and returns every audiobook under it, sorted by
// release date then title.
//
// An audiobook is an `.m4b` file with a sibling `.toml` of the same base name.
// A file missing its companion, or one that cannot be parsed, is logged and
// skipped: one malformed book should not take the whole library offline.
//
// The cache is consulted for durations and updated in place; entries for files
// that no longer exist are pruned.
func Library(mediaRoot string, cache *Cache, logger *slog.Logger) ([]audiobooks.Audiobook, error) {
	var books []audiobooks.Audiobook
	seen := map[string]bool{}

	err := filepath.WalkDir(mediaRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".m4b" {
			return nil
		}

		tomlPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".toml"
		if _, statErr := os.Stat(tomlPath); statErr != nil {
			logger.Warn("audiobook has no .toml metadata file, skipping", "m4b", path)
			return nil
		}

		// WalkDir normalises the paths it emits, so they are not always literally
		// prefixed by mediaRoot as written — a root of "./books" yields
		// "books/x.m4b". Trimming the string would leave the root in the path, and
		// from there in every enclosure URL and GUID.
		relPath, relErr := filepath.Rel(mediaRoot, path)
		if relErr != nil {
			logger.Warn("audiobook is not under the media root, skipping", "m4b", path, "error", relErr)
			return nil
		}
		relPath = "/" + filepath.ToSlash(relPath)
		seen[relPath] = true

		book, bookErr := readAudiobook(path, tomlPath, relPath, cache)
		if bookErr != nil {
			logger.Warn("could not read audiobook, skipping", "m4b", path, "error", bookErr)
			return nil
		}
		books = append(books, book)
		return nil
	})
	if err != nil {
		return nil, err
	}

	cache.Prune(seen)
	sortByReleaseThenTitle(books)
	return books, nil
}

// readAudiobook fills an audiobook from its metadata file, then overwrites the
// fields derived from the media file itself. The TOML is parsed first so that
// hand-written values can never shadow the real duration or file size.
func readAudiobook(m4bPath, tomlPath, relPath string, cache *Cache) (audiobooks.Audiobook, error) {
	var book audiobooks.Audiobook

	tomlFile, err := os.Open(tomlPath)
	if err != nil {
		return book, err
	}
	defer tomlFile.Close()
	if err := toml.NewDecoder(tomlFile).Decode(&book); err != nil {
		return book, err
	}

	info, err := os.Stat(m4bPath)
	if err != nil {
		return book, err
	}

	duration, cached := cache.Lookup(relPath, info.Size(), info.ModTime())
	if !cached {
		if duration, err = readDuration(m4bPath, info.Size()); err != nil {
			return book, err
		}
		cache.Store(relPath, info.Size(), info.ModTime(), duration)
	}

	book.Path = relPath
	book.FileSize = uint64(info.Size()) //nolint:gosec // file sizes are never negative
	book.Duration = duration
	book.MIMEType = mimeType
	// The cover sits beside the book, so its path is the book's with the
	// extension swapped — derived from relPath so it cannot disagree with it.
	if coverPath := findCover(m4bPath); coverPath != "" {
		book.ImagePath = strings.TrimSuffix(relPath, filepath.Ext(relPath)) + filepath.Ext(coverPath)
	}

	return book, nil
}

// readDuration parses the file's moov box. This is the only expensive part of
// scanning a book, and the only reason the cache exists.
func readDuration(m4bPath string, size int64) (time.Duration, error) {
	file, err := os.Open(m4bPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	info, err := mp4.OpenFromReader(file, size)
	if err != nil {
		return 0, err
	}
	if info.Moov == nil || info.Moov.Mvhd == nil || info.Moov.Mvhd.Timescale == 0 {
		return 0, nil
	}
	return time.Duration(
		(float32(info.Moov.Mvhd.Duration) / float32(info.Moov.Mvhd.Timescale)) * float32(time.Second),
	), nil
}

func findCover(m4bPath string) string {
	basePath := strings.TrimSuffix(m4bPath, filepath.Ext(m4bPath))
	for _, ext := range coverExtensions {
		coverPath := basePath + ext
		if _, err := os.Stat(coverPath); err == nil {
			return coverPath
		}
	}
	return ""
}

// sortByReleaseThenTitle gives the build a deterministic item order. Books with
// no release date sort first, matching how they are treated in feeds.
func sortByReleaseThenTitle(books []audiobooks.Audiobook) {
	slices.SortStableFunc(books, func(a, b audiobooks.Audiobook) int {
		if order := releaseTime(&a).Compare(releaseTime(&b)); order != 0 {
			return order
		}
		return cmp.Compare(a.Title, b.Title)
	})
}

func releaseTime(book *audiobooks.Audiobook) time.Time {
	if book.ReleaseDate == nil {
		return time.Time{}
	}
	return book.ReleaseDate.AsTime(time.UTC)
}
