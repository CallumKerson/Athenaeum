package scan

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const mediaRoot = "../../testdata"

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestLibrary(t *testing.T) {
	books, err := Library(mediaRoot, NewCache(), testLogger())
	require.NoError(t, err)
	require.Len(t, books, 2)

	// Sorted by release date, so the 1968 book comes before the 2019 one.
	assert.Equal(t, "A Wizard of Earthsea", books[0].Title)
	assert.Equal(t, "This Is How You Lose the Time War", books[1].Title)

	earthsea := books[0]
	// Paths are relative to the media root and keep their leading slash: they are
	// concatenated straight onto the media URL prefix to form each item's GUID.
	assert.Equal(t, "/Ursula K Le Guin/Earthsea/1 A Wizard of Earthsea/A Wizard of Earthsea.m4b", earthsea.Path)
	assert.Equal(t, "/Ursula K Le Guin/Earthsea/1 A Wizard of Earthsea/A Wizard of Earthsea.png", earthsea.ImagePath)
	assert.Equal(t, "audio/mp4a-latm", earthsea.MIMEType)
	assert.Positive(t, earthsea.FileSize)
	assert.Positive(t, earthsea.Duration)
	assert.Equal(t, []audiobooks.Genre{audiobooks.Childrens, audiobooks.Fantasy}, earthsea.Genres)
	assert.Equal(t, []string{"Ursula K. Le Guin"}, earthsea.Authors)
}

// A trailing separator on the media root must not leave the relative paths
// without their leading slash, which would corrupt every GUID.
func TestLibraryToleratesTrailingSeparator(t *testing.T) {
	withSlash, err := Library(mediaRoot+"/", NewCache(), testLogger())
	require.NoError(t, err)
	withoutSlash, err := Library(mediaRoot, NewCache(), testLogger())
	require.NoError(t, err)

	assert.Equal(t, withoutSlash, withSlash)
}

// One book missing its metadata should cost that book, not the whole build.
func TestLibrarySkipsM4BWithoutMetadata(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "orphan.m4b"), []byte("not really an m4b"), 0o644))

	books, err := Library(root, NewCache(), testLogger())

	require.NoError(t, err)
	assert.Empty(t, books)
}

func TestLibraryUsesAndFillsCache(t *testing.T) {
	cache := NewCache()

	_, err := Library(mediaRoot, cache, testLogger())
	require.NoError(t, err)
	assert.Equal(t, 2, cache.Misses)
	assert.Zero(t, cache.Hits)

	warm := *cache
	warm.Hits, warm.Misses = 0, 0
	books, err := Library(mediaRoot, &warm, testLogger())
	require.NoError(t, err)
	assert.Equal(t, 2, warm.Hits)
	assert.Zero(t, warm.Misses)
	assert.Len(t, books, 2)
}

// Metadata edits do not touch the .m4b, so they must not be cached away.
func TestLibraryPicksUpMetadataEditsWithWarmCache(t *testing.T) {
	root := t.TempDir()
	m4bPath := filepath.Join(root, "Book.m4b")
	tomlPath := filepath.Join(root, "Book.toml")
	require.NoError(t, copyFile(
		filepath.Join(mediaRoot, "Ursula K Le Guin", "Earthsea", "1 A Wizard of Earthsea", "A Wizard of Earthsea.m4b"),
		m4bPath))
	require.NoError(t, os.WriteFile(tomlPath, []byte("Title = \"Before\"\n"), 0o644))

	cache := NewCache()
	books, err := Library(root, cache, testLogger())
	require.NoError(t, err)
	require.Len(t, books, 1)
	assert.Equal(t, "Before", books[0].Title)

	require.NoError(t, os.WriteFile(tomlPath, []byte("Title = \"After\"\n"), 0o644))

	books, err = Library(root, cache, testLogger())
	require.NoError(t, err)
	require.Len(t, books, 1)
	assert.Equal(t, "After", books[0].Title)
}

func TestCacheLookup(t *testing.T) {
	modTime := time.Date(2026, time.July, 25, 12, 0, 0, 0, time.UTC)
	cache := NewCache()
	cache.Store("/Book.m4b", 100, modTime, 5*time.Hour)

	duration, found := cache.Lookup("/Book.m4b", 100, modTime)
	assert.True(t, found)
	assert.Equal(t, 5*time.Hour, duration)

	_, found = cache.Lookup("/Book.m4b", 101, modTime)
	assert.False(t, found, "a changed size should miss")

	_, found = cache.Lookup("/Book.m4b", 100, modTime.Add(time.Second))
	assert.False(t, found, "a changed modification time should miss")

	_, found = cache.Lookup("/Other.m4b", 100, modTime)
	assert.False(t, found)
}

// Without pruning the cache would grow forever as books are removed.
func TestCacheRoundTripAndPrune(t *testing.T) {
	modTime := time.Date(2026, time.July, 25, 12, 0, 0, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "nested", "m4b.json")

	cache := NewCache()
	cache.Store("/Kept.m4b", 1, modTime, time.Hour)
	cache.Store("/Removed.m4b", 2, modTime, time.Hour)
	cache.Prune(map[string]bool{"/Kept.m4b": true})
	require.NoError(t, cache.Save(path))

	loaded, err := LoadCache(path)
	require.NoError(t, err)

	_, found := loaded.Lookup("/Kept.m4b", 1, modTime)
	assert.True(t, found)
	_, found = loaded.Lookup("/Removed.m4b", 2, modTime)
	assert.False(t, found)
}

// A missing cache just means a slower build, never a failed one.
func TestLoadCacheMissingFile(t *testing.T) {
	cache, err := LoadCache(filepath.Join(t.TempDir(), "absent.json"))
	require.NoError(t, err)
	assert.NotNil(t, cache)
}

func TestLoadCacheCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.json")
	require.NoError(t, os.WriteFile(path, []byte("{not json"), 0o644))

	cache, err := LoadCache(path)

	require.Error(t, err)
	assert.NotNil(t, cache, "a usable empty cache should still come back")
}

func copyFile(from, to string) error {
	contents, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	return os.WriteFile(to, contents, 0o644)
}
