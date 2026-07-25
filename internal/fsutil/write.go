// Package fsutil holds the handful of filesystem helpers shared by the scanner
// and the site builder.
package fsutil

import (
	"bytes"
	"os"
	"path/filepath"
)

// WriteAtomic writes via a temporary file in the same directory followed by a
// rename, so a reader never observes a partially written file.
func WriteAtomic(path string, contents []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err = tmp.Write(contents); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// os.CreateTemp makes the file 0600, which is wrong for something a web
	// server has to read — it only works while the server runs as the same user.
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// WriteIfChanged writes contents only when they differ from what is already on
// disk, and reports whether it wrote. Leaving unchanged files alone keeps their
// modification times stable, which keeps the serving proxy's ETags stable, which
// keeps subscribers getting 304s across rebuilds.
func WriteIfChanged(path string, contents []byte) (bool, error) {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, contents) {
		return false, nil
	}
	return true, WriteAtomic(path, contents)
}
