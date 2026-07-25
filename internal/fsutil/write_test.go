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
