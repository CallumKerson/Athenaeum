package fsutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAtomicCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "file.txt")

	require.NoError(t, WriteAtomic(path, []byte("hello")))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(contents))
}

// os.CreateTemp defaults to 0600, which would make the published files
// unreadable by a web server running as any other user.
func TestWriteAtomicIsWorldReadable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, WriteAtomic(path, []byte("hello")))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
}

func TestWriteAtomicLeavesNoTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, WriteAtomic(filepath.Join(dir, "file.txt"), []byte("hello")))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "file.txt", entries[0].Name())
}

// Leaving unchanged files alone is what keeps their modification times, and so
// the serving proxy's ETags, stable across rebuilds.
func TestWriteIfChanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")

	changed, err := WriteIfChanged(path, []byte("first"))
	require.NoError(t, err)
	assert.True(t, changed)

	before, err := os.Stat(path)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	changed, err = WriteIfChanged(path, []byte("first"))
	require.NoError(t, err)
	assert.False(t, changed)

	after, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, before.ModTime(), after.ModTime())

	changed, err = WriteIfChanged(path, []byte("second"))
	require.NoError(t, err)
	assert.True(t, changed)

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "second", string(contents))
}

func TestWriteAtomicParentIsNotADirectory(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("a file, not a directory"), 0o600))

	err := WriteAtomic(filepath.Join(blocker, "nested", "file.txt"), []byte("hello"))

	require.Error(t, err)
}

// The temp file lands beside its target, so a directory that cannot be written
// to fails before anything is renamed into place.
func TestWriteAtomicUnwritableDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	dir := t.TempDir()
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	err := WriteAtomic(filepath.Join(dir, "file.txt"), []byte("hello"))

	require.Error(t, err)
	assert.NoFileExists(t, filepath.Join(dir, "file.txt"))
}

// A write that fails must not leave the previous contents replaced, nor leave a
// half-written temp file behind for the sweep to trip over.
func TestWriteAtomicFailureLeavesExistingFileIntact(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	require.NoError(t, WriteAtomic(path, []byte("original")))
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	require.Error(t, WriteAtomic(path, []byte("replacement")))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "original", string(contents))

	require.NoError(t, os.Chmod(dir, 0o700))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "a temp file was left behind")
}
