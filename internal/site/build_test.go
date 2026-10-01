package site

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/feed"
	"github.com/CallumKerson/Athenaeum/internal/scan"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current output")

const mediaRoot = "../../testdata"

// TestPlanGolden covers the whole pipeline at once — scanning, grouping, alias
// expansion, URL layout and feed rendering — by hashing every path the build
// would publish and its contents.
//
// It hashes the plan rather than walking the built tree on purpose: a
// case-insensitive filesystem folds `genre/Fantasy` and `genre/fantasy` into one
// file, so an on-disk golden tree would differ between macOS and Linux CI.
//
// Run with -update to regenerate.
func TestPlanGolden(t *testing.T) {
	manifest := renderManifest(t, testPlan(t).Feeds)
	goldenPath := filepath.Join("testdata", "golden_manifest.txt")

	if *update {
		require.NoError(t, os.WriteFile(goldenPath, []byte(manifest), 0o644))
		t.Log("golden manifest updated at", goldenPath)
		return
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)
	assert.Equal(t, string(expected), manifest)
}

// TestMainFeedGolden keeps one feed's full text under review, so a change to the
// rendered XML shows up as a readable diff rather than a changed hash. The other
// feed shapes are covered byte for byte in the feed package.
func TestMainFeedGolden(t *testing.T) {
	pages := testPlan(t).Feeds
	index := slices.IndexFunc(pages, func(page Page) bool {
		return slices.Contains(page.Paths, "podcast/feed.rss")
	})
	require.GreaterOrEqual(t, index, 0)

	rendered, err := testRenderer().Render(pages[index].Books, pages[index].Title, pages[index].Description)
	require.NoError(t, err)

	goldenPath := filepath.Join("testdata", "golden_feed.rss")
	if *update {
		require.NoError(t, os.WriteFile(goldenPath, rendered, 0o644))
		return
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)
	assert.Equal(t, string(expected), string(rendered))
}

// Every path the plan names must be readable from the built tree with the
// expected contents. Case-folded aliases resolve to the same file, which holds
// identical bytes either way, so this passes on both kinds of filesystem.
func TestBuildPublishesEveryPlannedPath(t *testing.T) {
	root := t.TempDir()
	content := testPlan(t)
	pages := content.Feeds

	_, err := Build(root, content, testRenderer(), true, testLogger())
	require.NoError(t, err)

	for index := range pages {
		expected, renderErr := testRenderer().Render(pages[index].Books, pages[index].Title, pages[index].Description)
		require.NoError(t, renderErr)
		for _, relPath := range pages[index].Paths {
			actual, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(relPath)))
			require.NoError(t, readErr, relPath)
			assert.Equal(t, string(expected), string(actual), relPath)
		}
	}

	assert.FileExists(t, filepath.Join(root, "index.html"))
	assert.FileExists(t, filepath.Join(root, "static", "itunes_image.jpg"))
	assert.FileExists(t, filepath.Join(root, "static", "itunes_image_small.jpg"))
	assert.FileExists(t, filepath.Join(root, MarkerName))
}

// A rebuild with no library changes must not touch a single file: stable
// modification times are what keep the serving proxy's ETags, and therefore
// subscribers' 304s, stable.
func TestBuildSkipsUnchangedFiles(t *testing.T) {
	root := t.TempDir()

	first, err := buildTestSite(t, root)
	require.NoError(t, err)
	assert.Positive(t, first.Written)

	second, err := buildTestSite(t, root)
	require.NoError(t, err)
	assert.Zero(t, second.Written)
	assert.Zero(t, second.Removed)
}

func TestBuildSweepsStaleFiles(t *testing.T) {
	root := t.TempDir()
	stale := filepath.Join(root, "podcast", "authors", "Someone Who Left", "feed.rss")

	content := testPlan(t)
	content.Feeds = append(content.Feeds, Page{
		Title: "Someone Who Left",
		Paths: []string{"podcast/authors/Someone Who Left/feed.rss"},
	})
	_, err := Build(root, content, testRenderer(), true, testLogger())
	require.NoError(t, err)
	require.FileExists(t, stale)

	result, err := buildTestSite(t, root)
	require.NoError(t, err)

	assert.Equal(t, 1, result.Removed)
	assert.NoFileExists(t, stale)
	assert.NoDirExists(t, filepath.Dir(stale), "emptied directories should be pruned")
	assert.DirExists(t, filepath.Join(root, "podcast", "authors"))
}

// Only files a previous build wrote are swept. Anything else in the site root —
// a robots.txt, or a symlink to the media library — belongs to someone else.
func TestBuildKeepsFilesItDidNotWrite(t *testing.T) {
	root := t.TempDir()
	_, err := buildTestSite(t, root)
	require.NoError(t, err)

	robots := filepath.Join(root, "robots.txt")
	require.NoError(t, os.WriteFile(robots, []byte("User-agent: *\n"), 0o644))
	media := filepath.Join(root, "media")
	require.NoError(t, os.Symlink(t.TempDir(), media))

	result, err := buildTestSite(t, root)
	require.NoError(t, err)

	assert.Zero(t, result.Removed)
	assert.FileExists(t, robots)
	_, err = os.Lstat(media)
	assert.NoError(t, err)
}

func TestBuildKeepsStaleFilesWhenSweepDisabled(t *testing.T) {
	root := t.TempDir()
	_, err := buildTestSite(t, root)
	require.NoError(t, err)

	_, err = Build(root, Content{}, testRenderer(), false, testLogger())
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(root, "podcast", "feed.rss"))

	// The skipped files are still remembered as the build's own, so the next
	// build that does sweep removes them.
	result, err := Build(root, Content{}, testRenderer(), true, testLogger())
	require.NoError(t, err)
	assert.Positive(t, result.Removed)
	assert.NoFileExists(t, filepath.Join(root, "podcast", "feed.rss"))
}

// On a case-insensitive filesystem, renaming an author from "Le Guin" to
// "LE GUIN" writes the new feed into the old file. On a case-sensitive
// filesystem a hard link stands in for that.
func TestSweepKeepsFileRenamedOnlyInCase(t *testing.T) {
	root := t.TempDir()
	oldPath := filepath.Join(root, "podcast", "authors", "Le Guin", "feed.rss")
	newPath := filepath.Join(root, "podcast", "authors", "LE GUIN", "feed.rss")
	require.NoError(t, os.MkdirAll(filepath.Dir(oldPath), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Dir(newPath), 0o755))
	require.NoError(t, os.WriteFile(newPath, []byte("feed"), 0o644))
	if _, err := os.Stat(oldPath); err != nil {
		require.NoError(t, os.Link(newPath, oldPath))
	}

	removed, err := sweepStale(root,
		map[string]bool{"podcast/authors/Le Guin/feed.rss": true},
		map[string]bool{"podcast/authors/LE GUIN/feed.rss": true},
		testLogger())
	require.NoError(t, err)

	assert.Zero(t, removed)
	assert.FileExists(t, oldPath)
}

// A case-only rename on a case-sensitive filesystem leaves two distinct files,
// and the old one is stale.
func TestSweepRemovesDistinctFileDifferingOnlyInCase(t *testing.T) {
	root := t.TempDir()
	oldPath := filepath.Join(root, "Le Guin.rss")
	newPath := filepath.Join(root, "LE GUIN.rss")
	require.NoError(t, os.WriteFile(newPath, []byte("new"), 0o644))
	if _, err := os.Stat(oldPath); err == nil {
		t.Skip("filesystem is case-insensitive")
	}
	require.NoError(t, os.WriteFile(oldPath, []byte("old"), 0o644))

	removed, err := sweepStale(root,
		map[string]bool{"Le Guin.rss": true},
		map[string]bool{"LE GUIN.rss": true},
		testLogger())
	require.NoError(t, err)

	assert.Equal(t, 1, removed)
	assert.NoFileExists(t, oldPath)
	assert.FileExists(t, newPath)
}

// The manifest is read back from disk, so a hand-edited one must not be able to
// send the sweep outside the site root.
func TestBuildIgnoresManifestPathsOutsideRoot(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "site")
	outside := filepath.Join(parent, "outside.txt")
	require.NoError(t, os.WriteFile(outside, []byte("mine"), 0o644))

	_, err := Build(root, Content{}, testRenderer(), true, testLogger())
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifestName), []byte("../outside.txt\n"), 0o644))

	_, err = Build(root, Content{}, testRenderer(), true, testLogger())
	require.NoError(t, err)
	assert.FileExists(t, outside)
}

// The sweep deletes files, so it must never take ownership of a directory it
// did not create.
func TestBuildRefusesDirectoryWithoutMarker(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "important.txt"), []byte("mine"), 0o644))

	_, err := Build(root, Content{}, testRenderer(), true, testLogger())

	require.ErrorIs(t, err, errNotSiteDir)
	assert.FileExists(t, filepath.Join(root, "important.txt"))
}

func TestBuildRejectsRootThatIsAFile(t *testing.T) {
	root := filepath.Join(t.TempDir(), "site")
	require.NoError(t, os.WriteFile(root, []byte("not a directory"), 0o644))

	_, err := Build(root, Content{}, testRenderer(), true, testLogger())

	require.Error(t, err)
}

func TestBuildAcceptsEmptyDirectory(t *testing.T) {
	_, err := Build(t.TempDir(), Content{}, testRenderer(), true, testLogger())
	require.NoError(t, err)
}

func testPlan(t *testing.T) Content {
	t.Helper()
	books, err := scan.Library(mediaRoot, scan.NewCache(), testLogger())
	require.NoError(t, err)
	require.NotEmpty(t, books)
	return Plan(books, []audiobooks.Genre{audiobooks.Erotica})
}

func buildTestSite(t *testing.T, root string) (Result, error) {
	t.Helper()
	return Build(root, testPlan(t), testRenderer(), true, testLogger())
}

// renderManifest lists every published path with the hash of its contents, one
// per line, sorted.
func renderManifest(t *testing.T, pages []Page) string {
	t.Helper()
	var lines []string
	for index := range pages {
		rendered, err := testRenderer().Render(pages[index].Books, pages[index].Title, pages[index].Description)
		require.NoError(t, err)
		digest := sha256.Sum256(rendered)
		for _, relPath := range pages[index].Paths {
			lines = append(lines, fmt.Sprintf("%x  %s", digest, relPath))
		}
	}
	slices.Sort(lines)
	return strings.Join(lines, "\n") + "\n"
}

func testRenderer() *feed.Renderer {
	return &feed.Renderer{
		Host:               "http://www.example-podcast.com/audiobooks",
		MediaPath:          "/media",
		ImageLink:          "http://www.example-podcast.com/audiobooks/static/itunes_image.jpg",
		Language:           "EN",
		Author:             "Athenaeum",
		Email:              "athenaeum@example.com",
		Explicit:           true,
		HandlePreUnixEpoch: true,
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
}
