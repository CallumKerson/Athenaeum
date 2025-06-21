package audiobook

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	nonExistentDirectory = "/non/existent/directory"
	testMediaRoot        = "/test/media"
)

func TestScanAll_WithTestData(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)

	books, err := ScanAll(context.TODO(), mediaRoot, logger)

	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestScanAll_NonExistentDirectory(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := nonExistentDirectory

	books, err := ScanAll(context.TODO(), mediaRoot, logger)

	assert.Error(t, err)
	assert.Nil(t, books)
}

func TestScanAll_EmptyDirectory(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := t.TempDir()

	books, err := ScanAll(context.TODO(), mediaRoot, logger)

	assert.NoError(t, err)
	assert.Empty(t, books)
}

func TestScanForUpdates_NoChanges(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)
	existing := testbooks.Audiobooks

	books, changed, err := ScanForUpdates(context.TODO(), mediaRoot, existing, logger)

	assert.NoError(t, err)
	assert.False(t, changed)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestScanForUpdates_WithChanges(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)
	existing := []audiobooks.Audiobook{testbooks.Audiobooks[0]} // Only one book

	books, changed, err := ScanForUpdates(context.TODO(), mediaRoot, existing, logger)

	assert.NoError(t, err)
	assert.True(t, changed)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestScanForUpdates_EmptyExisting(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)
	existing := []audiobooks.Audiobook{}

	books, changed, err := ScanForUpdates(context.TODO(), mediaRoot, existing, logger)

	assert.NoError(t, err)
	assert.True(t, changed)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestScanForUpdates_NonExistentDirectory(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := nonExistentDirectory
	existing := []audiobooks.Audiobook{}

	books, changed, err := ScanForUpdates(context.TODO(), mediaRoot, existing, logger)

	assert.Error(t, err)
	assert.False(t, changed)
	assert.Empty(t, books)
}

func TestGetM4BPathFromTOMLPath(t *testing.T) {
	tests := []struct {
		name     string
		tomlPath string
		expected string
	}{
		{
			name:     "basic conversion",
			tomlPath: "/path/to/book.toml",
			expected: "/path/to/book.m4b",
		},
		{
			name:     "nested path",
			tomlPath: "/media/Author/Book/book.toml",
			expected: "/media/Author/Book/book.m4b",
		},
		{
			name:     "complex filename",
			tomlPath: "/path/my-great-book_v2.toml",
			expected: "/path/my-great-book_v2.m4b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getM4BPathFromTOMLPath(tt.tomlPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTOMLPathFromM4BPath(t *testing.T) {
	tests := []struct {
		name     string
		m4bPath  string
		expected string
	}{
		{
			name:     "basic conversion",
			m4bPath:  "/path/to/book.m4b",
			expected: "/path/to/book.toml",
		},
		{
			name:     "nested path",
			m4bPath:  "/media/Author/Book/book.m4b",
			expected: "/media/Author/Book/book.toml",
		},
		{
			name:     "complex filename",
			m4bPath:  "/path/my-great-book_v2.m4b",
			expected: "/path/my-great-book_v2.toml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTOMLPathFromM4BPath(tt.m4bPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListToMap(t *testing.T) {
	mediaRoot := testMediaRoot
	books := []audiobooks.Audiobook{
		{Path: "book1.m4b", Title: "Book 1"},
		{Path: "book2.m4b", Title: "Book 2"},
	}

	result := listToMap(mediaRoot, books)

	assert.Len(t, result, 2)
	// The map uses the full path as key (root + relative path)
	assert.Contains(t, result, "/test/media/book1.m4b")
	assert.Contains(t, result, "/test/media/book2.m4b")
	assert.Equal(t, books[0], result["/test/media/book1.m4b"])
	assert.Equal(t, books[1], result["/test/media/book2.m4b"])
}

func TestListToMap_EmptyList(t *testing.T) {
	mediaRoot := testMediaRoot
	books := []audiobooks.Audiobook{}

	result := listToMap(mediaRoot, books)

	assert.Empty(t, result)
}

func TestListToMap_DifferentMediaRoot(t *testing.T) {
	mediaRoot := "/different/media"
	books := []audiobooks.Audiobook{
		{Path: "book1.m4b", Title: "Book 1"},
		{Path: "book2.m4b", Title: "Book 2"},
	}

	result := listToMap(mediaRoot, books)

	// All books should be included with mediaRoot prepended
	assert.Len(t, result, 2)
	assert.Contains(t, result, "/different/media/book1.m4b")
	assert.Contains(t, result, "/different/media/book2.m4b")
}

func TestParseM4BInfo_WithValidTestFile(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)

	// Use one of the test m4b files
	m4bPath := filepath.Join(mediaRoot, "Amal El-Mohtar and Max Gladstone", "This Is How You Lose the Time War", "This Is How You Lose the Time War.m4b")
	book := &audiobooks.Audiobook{}

	err := parseM4BInfo(m4bPath, mediaRoot, book, logger)

	assert.NoError(t, err)
	assert.NotEmpty(t, book.Path)
	assert.NotZero(t, book.Duration)
	assert.NotZero(t, book.FileSize)
	assert.Equal(t, "audio/mp4a-latm", book.MIMEType)
}

func TestParseM4BInfo_NonExistentFile(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := testMediaRoot
	m4bPath := "/non/existent/file.m4b"
	book := &audiobooks.Audiobook{}

	err := parseM4BInfo(m4bPath, mediaRoot, book, logger)

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestTrackM4BParseTime(t *testing.T) {
	logger := logrus.NewLogger()
	start := time.Now().Add(-100 * time.Millisecond)

	// This function should not panic or error
	assert.NotPanics(t, func() {
		trackM4BParseTime(start, "/test/path.m4b", logger)
	})
}

func TestGetAudiobook_WithValidTestFiles(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := dataloader.GetRootTestdata(t)

	m4bPath := filepath.Join(mediaRoot, "Amal El-Mohtar and Max Gladstone", "This Is How You Lose the Time War", "This Is How You Lose the Time War.m4b")
	tomlPath := filepath.Join(mediaRoot, "Amal El-Mohtar and Max Gladstone", "This Is How You Lose the Time War", "This Is How You Lose the Time War.toml")

	book, err := getAudiobook(m4bPath, tomlPath, mediaRoot, logger)

	assert.NoError(t, err)
	assert.NotEmpty(t, book.Title)
	assert.NotEmpty(t, book.Path)
	assert.NotZero(t, book.Duration)
	assert.NotZero(t, book.FileSize)
}

func TestGetAudiobook_Integration(t *testing.T) {
	// This is more of an integration test that would require actual test files
	// For now, we test the error path
	logger := logrus.NewLogger()
	mediaRoot := testMediaRoot

	book, err := getAudiobook("/non/existent/file.m4b", "/non/existent/file.toml", mediaRoot, logger)

	assert.Error(t, err)
	assert.Equal(t, audiobooks.Audiobook{}, book)
}

func TestScanForUpdates_BookRemoved(t *testing.T) {
	logger := logrus.NewLogger()
	mediaRoot := t.TempDir()

	// Existing books include one that no longer exists in the media root
	existing := []audiobooks.Audiobook{
		{Path: "removed-book.m4b", Title: "Removed Book"},
	}

	books, changed, err := ScanForUpdates(context.TODO(), mediaRoot, existing, logger)

	assert.NoError(t, err)
	assert.True(t, changed) // Should detect changes because a book was removed
	assert.Empty(t, books)  // No books should remain
}
