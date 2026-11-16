package cover

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stripes is a PNG three times wider than it is tall: red, green and blue thirds.
func stripes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	for column := range 300 {
		for row := range 100 {
			img.Set(column, row, []color.Color{
				color.RGBA{R: 255, A: 255},
				color.RGBA{G: 255, A: 255},
				color.RGBA{B: 255, A: 255},
			}[column/100])
		}
	}
	return encodePNG(t, img)
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func decode(t *testing.T, contents []byte) image.Image {
	t.Helper()
	img, err := jpeg.Decode(bytes.NewReader(contents))
	require.NoError(t, err)
	return img
}

// near reports whether two colours match within JPEG's lossiness.
func near(t *testing.T, want color.RGBA, got color.Color) {
	t.Helper()
	r, g, b, _ := got.RGBA()
	const tolerance = 40
	assert.InDelta(t, want.R, r>>8, tolerance, "red")
	assert.InDelta(t, want.G, g>>8, tolerance, "green")
	assert.InDelta(t, want.B, b>>8, tolerance, "blue")
}

func TestThumbnailCropsToTheCentredSquare(t *testing.T) {
	thumbnail, err := Thumbnail(stripes(t))
	require.NoError(t, err)

	img := decode(t, thumbnail)
	assert.Equal(t, image.Rect(0, 0, Size, Size), img.Bounds())
	// Only the green middle third survives the crop, right out to the edges.
	green := color.RGBA{G: 255, A: 255}
	near(t, green, img.At(Size/2, Size/2))
	near(t, green, img.At(4, Size/2))
	near(t, green, img.At(Size-5, Size/2))
}

// A rebuild must write identical bytes, or every thumbnail would change on
// every build and break the stable mtimes the site relies on.
func TestThumbnailIsDeterministic(t *testing.T) {
	first, err := Thumbnail(stripes(t))
	require.NoError(t, err)
	second, err := Thumbnail(stripes(t))
	require.NoError(t, err)
	assert.Equal(t, first, second)
}

func TestThumbnailFlattensTransparencyOntoWhite(t *testing.T) {
	thumbnail, err := Thumbnail(encodePNG(t, image.NewNRGBA(image.Rect(0, 0, 50, 50))))
	require.NoError(t, err)
	near(t, color.RGBA{R: 255, G: 255, B: 255, A: 255}, decode(t, thumbnail).At(Size/2, Size/2))
}

func TestThumbnailRejectsWhatIsNotAnImage(t *testing.T) {
	_, err := Thumbnail([]byte("not an image"))
	assert.Error(t, err)
}

// library writes a cover into a fresh media root and returns the root and the
// cover's media-root-relative path.
func library(t *testing.T) (mediaRoot, imagePath string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "Author", "Book")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Book.png"), stripes(t), 0o644))
	return root, "/Author/Book/Book.png"
}

func cached(t *testing.T, dir string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, "*.jpg"))
	require.NoError(t, err)
	return matches
}

func TestSourceReusesTheCachedThumbnailUntilTheCoverChanges(t *testing.T) {
	mediaRoot, imagePath := library(t)
	cacheDir := filepath.Join(t.TempDir(), "covers")
	source := NewSource(mediaRoot, cacheDir)

	made, err := source.Thumbnail(imagePath)
	require.NoError(t, err)
	entries := cached(t, cacheDir)
	require.Len(t, entries, 1)

	// Whatever is in the cache is what comes back, so the cover was not decoded.
	require.NoError(t, os.WriteFile(entries[0], []byte("from the cache"), 0o644))
	again, err := source.Thumbnail(imagePath)
	require.NoError(t, err)
	assert.Equal(t, []byte("from the cache"), again)

	// A cover with a new modification time is a new cover.
	later := time.Now().Add(time.Hour)
	require.NoError(t, os.Chtimes(filepath.Join(mediaRoot, imagePath), later, later))
	remade, err := source.Thumbnail(imagePath)
	require.NoError(t, err)
	assert.Equal(t, made, remade)
	assert.Len(t, cached(t, cacheDir), 2)
}

func TestSourceWithoutACacheDirCachesNothing(t *testing.T) {
	mediaRoot, imagePath := library(t)
	source := NewSource(mediaRoot, "")

	thumbnail, err := source.Thumbnail(imagePath)
	require.NoError(t, err)
	assert.NotEmpty(t, thumbnail)
	assert.NoError(t, source.Prune())
}

func TestSourceReportsAMissingCover(t *testing.T) {
	mediaRoot, _ := library(t)
	_, err := NewSource(mediaRoot, t.TempDir()).Thumbnail("/Author/Book/Missing.png")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestPruneKeepsOnlyTheThumbnailsThisBuildUsed(t *testing.T) {
	mediaRoot, imagePath := library(t)
	cacheDir := t.TempDir()
	stale := filepath.Join(cacheDir, "stale.jpg")
	unrelated := filepath.Join(cacheDir, "notes.txt")
	require.NoError(t, os.WriteFile(stale, []byte("old"), 0o644))
	require.NoError(t, os.WriteFile(unrelated, []byte("keep"), 0o644))

	source := NewSource(mediaRoot, cacheDir)
	_, err := source.Thumbnail(imagePath)
	require.NoError(t, err)
	require.NoError(t, source.Prune())

	assert.NoFileExists(t, stale)
	assert.FileExists(t, unrelated)
	assert.Len(t, cached(t, cacheDir), 1)
}

func TestPruneToleratesACacheDirThatWasNeverMade(t *testing.T) {
	mediaRoot, _ := library(t)
	assert.NoError(t, NewSource(mediaRoot, filepath.Join(t.TempDir(), "never")).Prune())
}
