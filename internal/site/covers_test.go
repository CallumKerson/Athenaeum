package site

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var errCannotDecode = errors.New("cannot decode")

// fakeCovers returns the image path as the thumbnail, and fails for any path
// listed in broken.
type fakeCovers struct {
	broken map[string]bool
}

func (f fakeCovers) Thumbnail(imagePath string) ([]byte, error) {
	if f.broken[imagePath] {
		return nil, errCannotDecode
	}
	return []byte(imagePath), nil
}

func withCover(title string) audiobooks.Audiobook {
	audiobook := book(title, nil, []string{"Author"}, nil, nil)
	audiobook.ImagePath = "/" + title + "/" + title + ".jpg"
	return audiobook
}

var coverPathPattern = regexp.MustCompile(`^covers/[0-9a-f]{16}\.jpg$`)

func TestPlanNamesAThumbnailForEveryCover(t *testing.T) {
	content := Plan(
		[]audiobooks.Audiobook{withCover("One"), withCover("Two"), book("Bare", nil, []string{"Author"}, nil, nil)},
		Exclusions{},
	)

	require.Len(t, content.Covers, 2, "a book without a cover gets no thumbnail")
	assert.Regexp(t, coverPathPattern, content.Covers["/One/One.jpg"])
	assert.Regexp(t, coverPathPattern, content.Covers["/Two/Two.jpg"])
	assert.NotEqual(t, content.Covers["/One/One.jpg"], content.Covers["/Two/Two.jpg"])

	again := Plan([]audiobooks.Audiobook{withCover("Two"), withCover("One")}, Exclusions{})
	assert.Equal(t, content.Covers["/One/One.jpg"], again.Covers["/One/One.jpg"], "stable between builds")
}

func TestBuildPublishesCoversAndShowsThem(t *testing.T) {
	root := t.TempDir()
	content := Plan(
		[]audiobooks.Audiobook{withCover("One"), book("Bare", nil, []string{"Author"}, nil, nil)},
		Exclusions{},
	)

	_, err := Build(root, content, testRenderer(t), Options{Covers: fakeCovers{}, Sweep: true}, testLogger())
	require.NoError(t, err)

	thumbnail, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(content.Covers["/One/One.jpg"])))
	require.NoError(t, err)
	assert.Equal(t, "/One/One.jpg", string(thumbnail))

	page := readPage(t, root, "books")
	assert.Contains(t, page, `src="../`+content.Covers["/One/One.jpg"]+`"`)
	assert.Equal(t, 1, strings.Count(page, "<img"), "the book without a cover shows an empty square instead")
}

func TestBuildLeavesOutACoverThatCannotBeRead(t *testing.T) {
	root := t.TempDir()
	content := Plan([]audiobooks.Audiobook{withCover("One"), withCover("Broken")}, Exclusions{})
	covers := fakeCovers{broken: map[string]bool{"/Broken/Broken.jpg": true}}

	_, err := Build(root, content, testRenderer(t), Options{Covers: covers, Sweep: true}, testLogger())
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(root, filepath.FromSlash(content.Covers["/One/One.jpg"])))
	assert.NoFileExists(t, filepath.Join(root, filepath.FromSlash(content.Covers["/Broken/Broken.jpg"])))
	page := readPage(t, root, "books")
	assert.NotContains(t, page, content.Covers["/Broken/Broken.jpg"])
}

func TestBuildWithoutCoversPublishesNone(t *testing.T) {
	root := t.TempDir()
	_, err := Build(
		root,
		Plan([]audiobooks.Audiobook{withCover("One")}, Exclusions{}),
		testRenderer(t),
		Options{Sweep: true},
		testLogger(),
	)
	require.NoError(t, err)
	assert.NoDirExists(t, filepath.Join(root, coversDir))
	assert.NotContains(t, readPage(t, root, "books"), "<img")
}

func TestBuildSweepsTheThumbnailOfARemovedCover(t *testing.T) {
	root := t.TempDir()
	content := Plan([]audiobooks.Audiobook{withCover("One")}, Exclusions{})
	opts := Options{Covers: fakeCovers{}, Sweep: true}
	_, err := Build(root, content, testRenderer(t), opts, testLogger())
	require.NoError(t, err)
	thumbnail := filepath.Join(root, filepath.FromSlash(content.Covers["/One/One.jpg"]))
	require.FileExists(t, thumbnail)

	_, err = Build(root, Plan([]audiobooks.Audiobook{book("One", nil, []string{"Author"}, nil, nil)}, Exclusions{}),
		testRenderer(t), opts, testLogger())
	require.NoError(t, err)
	assert.NoFileExists(t, thumbnail)
}

func readPage(t *testing.T, root, dir string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, dir, indexName))
	require.NoError(t, err)
	return string(contents)
}
