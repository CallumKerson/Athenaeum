package main

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/site"
)

// mediaRoot is the pair of test audiobooks in the repository root.
const mediaRoot = "../../testdata"

// isolateEnv points the XDG variables at scratch directories so that a test can
// never read — or write — the developer's real config and duration cache.
func isolateEnv(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
}

func TestBuildCommandWritesSite(t *testing.T) {
	isolateEnv(t)
	out := t.TempDir()

	var stdout, stderr bytes.Buffer
	cmd := NewRootCommand()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{
		"build",
		"--media-root", mediaRoot,
		"--out", out,
		"--host", "https://books.example.com",
		"--cache", filepath.Join(t.TempDir(), "m4b.json"),
	})

	require.NoError(t, cmd.ExecuteContext(context.Background()))

	// The feed and the pages a subscriber and a browser respectively land on.
	assert.FileExists(t, filepath.Join(out, "podcast", "feed.rss"))
	assert.FileExists(t, filepath.Join(out, "index.html"))
	assert.FileExists(t, filepath.Join(out, "books", "index.html"))
	assert.FileExists(t, filepath.Join(out, site.MarkerName))

	feed, err := os.ReadFile(filepath.Join(out, "podcast", "feed.rss"))
	require.NoError(t, err)
	assert.Contains(t, string(feed), "A Wizard of Earthsea")
	// The host reaches the enclosure URLs, which are also the GUIDs.
	assert.Contains(t, string(feed), "https://books.example.com/media/")

	assert.Contains(t, stderr.String(), "scanned library")
}

// The whole point of the generator is that running it twice is a no-op, so that
// unchanged files keep their mtimes and the web server keeps serving 304s.
func TestBuildCommandIsIdempotent(t *testing.T) {
	isolateEnv(t)
	out := t.TempDir()
	cachePath := filepath.Join(t.TempDir(), "m4b.json")

	run := func() string {
		var stderr bytes.Buffer
		cmd := NewRootCommand()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&stderr)
		cmd.SetArgs([]string{
			"build", "--media-root", mediaRoot, "--out", out,
			"--host", "https://books.example.com", "--cache", cachePath,
		})
		require.NoError(t, cmd.ExecuteContext(context.Background()))
		return stderr.String()
	}

	run()
	feedPath := filepath.Join(out, "podcast", "feed.rss")
	first, err := os.Stat(feedPath)
	require.NoError(t, err)

	second := run()

	after, err := os.Stat(feedPath)
	require.NoError(t, err)
	assert.Equal(t, first.ModTime(), after.ModTime(), "unchanged feed was rewritten")
	assert.Contains(t, second, "written=0")
	// The second run should have read every duration from the cache.
	assert.Contains(t, second, "parsed=0")
}

func TestBuildCommandRejectsForeignOutputDirectory(t *testing.T) {
	isolateEnv(t)
	out := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(out, "important.txt"), []byte("mine"), 0o600))

	cmd := NewRootCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"build", "--media-root", mediaRoot, "--out", out,
		"--cache", filepath.Join(t.TempDir(), "m4b.json"),
	})

	err := cmd.ExecuteContext(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), site.MarkerName)
	assert.FileExists(t, filepath.Join(out, "important.txt"))
}

func TestBuildCommandRejectsUnknownExcludedGenre(t *testing.T) {
	isolateEnv(t)
	configPath := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(configPath, []byte(
		"[ExclusionsFromMainFeed]\nGenres = [\"spycraft\"]\n"), 0o600))

	cmd := NewRootCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"build", "--config", configPath,
		"--media-root", mediaRoot, "--out", t.TempDir(),
	})

	err := cmd.ExecuteContext(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "spycraft")
}

func TestResolveConfigFlagsBeatFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(configPath, []byte(`
Host = "https://from-file.example.com"

[Media]
Root = "/from/file/media"

[Site]
Root = "/from/file/site"
`), 0o600))

	cfg, err := resolveConfig(&buildFlags{
		configPath: configPath,
		mediaRoot:  "/from/flag/media",
		siteRoot:   "/from/flag/site",
		host:       "https://from-flag.example.com",
	}, &bytes.Buffer{})

	require.NoError(t, err)
	assert.Equal(t, "/from/flag/media", cfg.Media.Root)
	assert.Equal(t, "/from/flag/site", cfg.Site.Root)
	assert.Equal(t, "https://from-flag.example.com", cfg.Host)
}

func TestResolveConfigKeepsFileValuesWhenFlagsAreEmpty(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(configPath, []byte(`
[Media]
Root = "/from/file/media"

[Site]
Root = "/from/file/site"
`), 0o600))

	cfg, err := resolveConfig(&buildFlags{configPath: configPath}, &bytes.Buffer{})

	require.NoError(t, err)
	assert.Equal(t, "/from/file/media", cfg.Media.Root)
	assert.Equal(t, "/from/file/site", cfg.Site.Root)
}

func TestResolveConfigRequiresRoots(t *testing.T) {
	absent := filepath.Join(t.TempDir(), "absent.toml")

	_, err := resolveConfig(&buildFlags{configPath: absent}, &bytes.Buffer{})
	require.ErrorIs(t, err, errNoMediaRoot)

	_, err = resolveConfig(&buildFlags{configPath: absent, mediaRoot: "/media"}, &bytes.Buffer{})
	require.ErrorIs(t, err, errNoSiteRoot)
}

func TestResolveConfigPropagatesLoadError(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(configPath, []byte("Host = \n"), 0o600))

	_, err := resolveConfig(&buildFlags{configPath: configPath}, &bytes.Buffer{})

	require.Error(t, err)
}

func TestScanLibraryUsesCacheOnSecondRun(t *testing.T) {
	isolateEnv(t)
	cachePath := filepath.Join(t.TempDir(), "m4b.json")
	cfg := &BuildConfig{Media: Media{Root: mediaRoot}}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	books, cache, err := scanLibrary(cfg, &buildFlags{cachePath: cachePath}, logger)
	require.NoError(t, err)
	require.NotEmpty(t, books)
	assert.Zero(t, cache.Hits)
	assert.Positive(t, cache.Misses)

	_, warm, err := scanLibrary(cfg, &buildFlags{cachePath: cachePath}, logger)
	require.NoError(t, err)
	assert.Equal(t, len(books), warm.Hits)
	assert.Zero(t, warm.Misses)
}

// --no-cache must neither read nor write the cache file.
func TestScanLibraryNoCacheLeavesCacheFileAlone(t *testing.T) {
	isolateEnv(t)
	cachePath := filepath.Join(t.TempDir(), "m4b.json")
	cfg := &BuildConfig{Media: Media{Root: mediaRoot}}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	_, cache, err := scanLibrary(cfg, &buildFlags{cachePath: cachePath, noCache: true}, logger)

	require.NoError(t, err)
	assert.Zero(t, cache.Hits)
	assert.NoFileExists(t, cachePath)
}

// A corrupt cache costs a re-read, never a failed build.
func TestScanLibraryToleratesUnreadableCache(t *testing.T) {
	isolateEnv(t)
	cachePath := filepath.Join(t.TempDir(), "m4b.json")
	require.NoError(t, os.WriteFile(cachePath, []byte("{not json"), 0o600))

	var logged bytes.Buffer
	books, cache, err := scanLibrary(
		&BuildConfig{Media: Media{Root: mediaRoot}},
		&buildFlags{cachePath: cachePath},
		slog.New(slog.NewTextHandler(&logged, nil)),
	)

	require.NoError(t, err)
	assert.NotEmpty(t, books)
	assert.Positive(t, cache.Misses)
	assert.Contains(t, logged.String(), "could not read cache")
}

func TestScanLibraryFallsBackToDefaultCachePath(t *testing.T) {
	isolateEnv(t)

	_, _, err := scanLibrary(
		&BuildConfig{Media: Media{Root: mediaRoot}},
		&buildFlags{},
		slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)),
	)

	require.NoError(t, err)
	expected, err := DefaultCachePath()
	require.NoError(t, err)
	assert.FileExists(t, expected)
}

func TestScanLibraryMissingMediaRoot(t *testing.T) {
	isolateEnv(t)

	_, _, err := scanLibrary(
		&BuildConfig{Media: Media{Root: filepath.Join(t.TempDir(), "absent")}},
		&buildFlags{cachePath: filepath.Join(t.TempDir(), "m4b.json")},
		slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)),
	)

	require.Error(t, err)
}

func TestRendererForCarriesConfigThrough(t *testing.T) {
	renderer := rendererFor(&BuildConfig{
		Host: "https://books.example.com",
		Podcast: Podcast{
			Explicit:     true,
			Language:     "FR",
			Author:       "Someone",
			Email:        "someone@example.com",
			PreUnixEpoch: PreUnixEpoch{Handle: true},
		},
	})

	assert.Equal(t, "https://books.example.com", renderer.Host)
	assert.Equal(t, mediaPath, renderer.MediaPath)
	assert.Equal(t, "https://books.example.com/static/itunes_image.jpg", renderer.ImageLink)
	assert.Equal(t, "FR", renderer.Language)
	assert.Equal(t, "Someone", renderer.Author)
	assert.Equal(t, "someone@example.com", renderer.Email)
	assert.True(t, renderer.Explicit)
	assert.True(t, renderer.HandlePreUnixEpoch)
}

func TestBuildLoggerLevel(t *testing.T) {
	var quiet, verbose bytes.Buffer

	buildLogger(&quiet, false).Debug("hidden")
	buildLogger(&verbose, true).Debug("shown")

	assert.Empty(t, quiet.String())
	assert.Contains(t, verbose.String(), "shown")
}

// notifyOvercast talks to a hard-coded endpoint, so the reachable case here is
// the request never leaving: a cancelled context fails before any network use.
func TestNotifyOvercastCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := notifyOvercast(ctx, "https://books.example.com")

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestVersionCommand(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version"})

	require.NoError(t, cmd.Execute())

	assert.Contains(t, out.String(), "version: ")
	assert.Contains(t, out.String(), "commit:  ")
	assert.Contains(t, out.String(), "built at:")
}

func TestRootCommandHasSubcommands(t *testing.T) {
	subcommands := NewRootCommand().Commands()
	names := make([]string, 0, len(subcommands))
	for _, sub := range subcommands {
		names = append(names, sub.Name())
	}

	assert.Contains(t, names, "build")
	assert.Contains(t, names, "version")
}
