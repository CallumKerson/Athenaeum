// Package cover makes the small square thumbnails the HTML pages show in place
// of a book's full-size artwork, which can be thousands of pixels across.
package cover

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg" // Covers are JPEG or PNG; see scan.coverExtensions.
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/draw"

	"github.com/CallumKerson/Athenaeum/internal/fsutil"
)

// Size is the width and height of every thumbnail in pixels: twice the largest
// size a page shows a cover at, so it stays sharp on high-density screens.
const Size = 480

const quality = 80

var errEmptyImage = errors.New("image has no pixels")

// cacheVersion is part of every cache key. Bump it whenever Thumbnail's output
// changes, so cached thumbnails made the old way are not reused.
const cacheVersion = 1

// Thumbnail decodes an image, crops it to a centred square and scales it to
// Size, returning it as a JPEG. The output depends only on the input bytes, so
// rebuilding an unchanged cover writes an identical file.
func Thumbnail(source []byte) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	side := min(bounds.Dx(), bounds.Dy())
	if side == 0 {
		return nil, errEmptyImage
	}
	crop := image.Rect(0, 0, side, side).Add(image.Point{
		X: bounds.Min.X + (bounds.Dx()-side)/2,
		Y: bounds.Min.Y + (bounds.Dy()-side)/2,
	})

	dst := image.NewRGBA(image.Rect(0, 0, Size, Size))
	// JPEG has no transparency, so a transparent PNG is flattened onto white
	// rather than onto the black an unfilled RGBA would give it.
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, crop, draw.Over, nil)

	var out bytes.Buffer
	if err := jpeg.Encode(&out, dst, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// Source makes thumbnails for the covers in a media library, keeping each one
// in a cache directory so that only new or changed covers are decoded again.
// It is safe for concurrent use.
type Source struct {
	mediaRoot string
	// cacheDir is "" when caching is off, and every thumbnail is made afresh.
	cacheDir string

	mu   sync.Mutex
	used map[string]bool
}

// NewSource reads covers from mediaRoot and caches their thumbnails in
// cacheDir, or caches nothing when cacheDir is "".
func NewSource(mediaRoot, cacheDir string) *Source {
	return &Source{mediaRoot: mediaRoot, cacheDir: cacheDir, used: map[string]bool{}}
}

// Thumbnail returns the thumbnail of the cover at imagePath, a path relative to
// the media root like Audiobook.ImagePath.
func (s *Source) Thumbnail(imagePath string) ([]byte, error) {
	fullPath := filepath.Join(s.mediaRoot, filepath.FromSlash(imagePath))
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	var cachePath string
	if s.cacheDir != "" {
		name := cacheName(imagePath, info)
		s.mu.Lock()
		s.used[name] = true
		s.mu.Unlock()
		cachePath = filepath.Join(s.cacheDir, name)
		if cached, err := os.ReadFile(cachePath); err == nil {
			return cached, nil
		}
	}

	source, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	thumbnail, err := Thumbnail(source)
	if err != nil {
		return nil, fmt.Errorf("making thumbnail of %s: %w", imagePath, err)
	}
	if cachePath != "" {
		// A cache that cannot be written only costs time on the next build.
		if err := os.MkdirAll(s.cacheDir, 0o755); err == nil {
			_ = fsutil.WriteAtomic(cachePath, thumbnail)
		}
	}
	return thumbnail, nil
}

// Prune deletes every cached thumbnail that no call to Thumbnail asked for, so
// replaced and removed covers do not pile up in the cache. Only call it once
// every thumbnail the build needs has been made.
func (s *Source) Prune() error {
	if s.cacheDir == "" {
		return nil
	}
	entries, err := os.ReadDir(s.cacheDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, entry := range entries {
		if entry.IsDir() || s.used[entry.Name()] || !strings.HasSuffix(entry.Name(), ".jpg") {
			continue
		}
		if err := os.Remove(filepath.Join(s.cacheDir, entry.Name())); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

// cacheName identifies one version of one cover. Size and modification time
// stand in for the contents, as they do for the duration cache, so an unchanged
// cover is never read.
func cacheName(imagePath string, info os.FileInfo) string {
	key := strings.Join([]string{
		strconv.Itoa(cacheVersion),
		strconv.Itoa(Size),
		imagePath,
		strconv.FormatInt(info.Size(), 10),
		strconv.FormatInt(info.ModTime().UnixNano(), 10),
	}, "\x00")
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:]) + ".jpg"
}
