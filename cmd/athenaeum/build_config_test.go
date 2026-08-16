package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestLoadBuildConfigReadsTOML(t *testing.T) {
	path := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(path, []byte(`
Host = "https://books.example.com"

[Media]
Root = "/srv/audiobooks"

[Site]
Root = "/srv/www"

[Podcast]
Author = "Someone"
Email = "someone@example.com"
Language = "FR"
Explicit = false

[ThirdParty]
NotifyOvercast = true

[ExclusionsFromMainFeed]
Genres = ["Erotica"]
`), 0o600))

	var out bytes.Buffer
	cfg, err := LoadBuildConfig(path, &out)
	require.NoError(t, err)

	assert.Equal(t, "https://books.example.com", cfg.Host)
	assert.Equal(t, "/srv/audiobooks", cfg.Media.Root)
	assert.Equal(t, "/srv/www", cfg.Site.Root)
	assert.Equal(t, "Someone", cfg.Podcast.Author)
	assert.Equal(t, "FR", cfg.Podcast.Language)
	assert.False(t, cfg.Podcast.Explicit)
	assert.True(t, cfg.ThirdParty.NotifyOvercast)
	assert.Equal(t, []string{"Erotica"}, cfg.ExclusionsFromMainFeed.Genres)
	assert.Empty(t, out.String())
}

// A config file that sets nothing must not clear the defaults, since Explicit
// and PreUnixEpoch default to true and TOML has no way to say "unset".
func TestLoadBuildConfigKeepsDefaultsForAbsentKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(path, []byte("Host = \"https://books.example.com\"\n"), 0o600))

	cfg, err := LoadBuildConfig(path, &bytes.Buffer{})
	require.NoError(t, err)

	assert.True(t, cfg.Podcast.Explicit)
	assert.True(t, cfg.Podcast.PreUnixEpoch.Handle)
	assert.Equal(t, "EN", cfg.Podcast.Language)
}

func TestLoadBuildConfigMissingFileUsesDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nothing-here.toml")

	var out bytes.Buffer
	cfg, err := LoadBuildConfig(path, &out)
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8080", cfg.Host)
	assert.True(t, cfg.Podcast.Explicit)
	assert.Contains(t, out.String(), "No config file at")
}

func TestLoadBuildConfigInvalidTOML(t *testing.T) {
	path := filepath.Join(t.TempDir(), configName)
	require.NoError(t, os.WriteFile(path, []byte("Host = \n"), 0o600))

	_, err := LoadBuildConfig(path, &bytes.Buffer{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), path)
}

// The message exists so that the move from the server's YAML config to the XDG
// TOML one does not look like the config silently vanished.
func TestLoadBuildConfigNotesLegacyConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".athenaeum"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(home, legacyConfigPath), []byte("host: x\n"), 0o600))

	var out bytes.Buffer
	_, err := LoadBuildConfig(filepath.Join(t.TempDir(), "absent.toml"), &out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), ".athenaeum/config.yaml")
}

func TestLoadBuildConfigWithoutLegacyConfigSaysNothingExtra(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	var out bytes.Buffer
	_, err := LoadBuildConfig(filepath.Join(t.TempDir(), "absent.toml"), &out)
	require.NoError(t, err)

	assert.NotContains(t, out.String(), "config.yaml")
}

// os.UserConfigDir returns ~/Library/Application Support on macOS and ignores
// XDG_CONFIG_HOME entirely, which is why these paths are resolved by hand.
func TestDefaultPathsPreferXDGVariables(t *testing.T) {
	configHome, cacheHome := t.TempDir(), t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	t.Setenv("XDG_CACHE_HOME", cacheHome)

	configPath, err := DefaultConfigPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(configHome, appName, configName), configPath)

	cachePath, err := DefaultCachePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(cacheHome, appName, cacheName), cachePath)
}

func TestDefaultPathsFallBackToHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")

	configPath, err := DefaultConfigPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".config", appName, configName), configPath)

	cachePath, err := DefaultCachePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".cache", appName, cacheName), cachePath)
}

func TestGetGenres(t *testing.T) {
	genres, err := ExclusionsFromMainFeed{Genres: []string{"Erotica", "sci-fi"}}.GetGenres()

	require.NoError(t, err)
	assert.Equal(t, []audiobooks.Genre{audiobooks.Erotica, audiobooks.SciFi}, genres)
}

func TestGetGenresEmpty(t *testing.T) {
	genres, err := ExclusionsFromMainFeed{}.GetGenres()

	require.NoError(t, err)
	assert.Empty(t, genres)
}

func TestGetGenresUnknownGenre(t *testing.T) {
	_, err := ExclusionsFromMainFeed{Genres: []string{"Erotica", "spycraft"}}.GetGenres()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "spycraft")
}
