package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	// appName is the per-application directory under the XDG base directories.
	appName = "athenaeum"
	// configName is the config file inside the XDG config directory.
	configName = "athenaeum.toml"
	// cacheName holds the m4b durations learned by previous builds.
	cacheName = "m4b.json"
	// legacyConfigPath is checked only so an existing install gets told where its
	// config moved to. Delete once the migration is old news.
	legacyConfigPath = ".athenaeum/config.yaml"
)

// BuildConfig is everything `athenaeum build` needs. The key names match the
// server's old YAML config, so an existing config translates across unchanged.
type BuildConfig struct {
	Host                   string
	Media                  Media
	Site                   Site
	Podcast                Podcast
	ThirdParty             ThirdParty
	ExclusionsFromMainFeed ExclusionsFromMainFeed
}

// Media is where the audiobook library lives.
type Media struct {
	Root string
}

// Site is where the generated tree is written.
type Site struct {
	Root string
}

// Podcast holds the channel-level iTunes metadata shared by every feed.
type Podcast struct {
	Explicit     bool
	Language     string
	Author       string
	Email        string
	PreUnixEpoch PreUnixEpoch
}

// PreUnixEpoch controls whether release dates before 1970 are clamped up to the
// epoch, since podcast clients vary in how they treat negative timestamps.
type PreUnixEpoch struct {
	Handle bool
}

type ThirdParty struct {
	NotifyOvercast bool
}

// ExclusionsFromMainFeed names genres that appear in their own feed but not in
// the main one.
type ExclusionsFromMainFeed struct {
	Genres []string
}

func (e ExclusionsFromMainFeed) GetGenres() ([]audiobooks.Genre, error) {
	genres := make([]audiobooks.Genre, 0, len(e.Genres))
	for _, genreName := range e.Genres {
		genre, err := audiobooks.ParseGenre(genreName)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

// LoadBuildConfig reads the TOML config file, falling back to defaults for
// anything it does not set. A missing config file is not an error — every value
// can also come from a flag.
func LoadBuildConfig(pathToConfigFile string, out io.Writer) (*BuildConfig, error) {
	cfg := &BuildConfig{
		Host: "http://localhost:8080",
		Podcast: Podcast{
			Explicit:     true,
			Language:     "EN",
			PreUnixEpoch: PreUnixEpoch{Handle: true},
		},
	}

	if pathToConfigFile == "" {
		defaultPath, err := DefaultConfigPath()
		if err != nil {
			return nil, err
		}
		pathToConfigFile = defaultPath
	}

	contents, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		fmt.Fprintln(out, "No config file at", pathToConfigFile, "- using defaults and flags")
		noteLegacyConfig(out)
		return cfg, nil
	}

	if err = toml.Unmarshal(contents, cfg); err != nil {
		return nil, fmt.Errorf("reading %s: %w", pathToConfigFile, err)
	}
	return cfg, nil
}

// noteLegacyConfig points at the old YAML config if it is still lying around,
// so the move to XDG paths does not look like the config silently vanished.
func noteLegacyConfig(out io.Writer) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	legacy := filepath.Join(home, legacyConfigPath)
	if _, err = os.Stat(legacy); err == nil {
		fmt.Fprintln(out, "Found the old server config at", legacy+".",
			"`athenaeum build` reads TOML from the XDG path instead.")
	}
}

// DefaultConfigPath is $XDG_CONFIG_HOME/athenaeum/athenaeum.toml.
func DefaultConfigPath() (string, error) {
	return xdgPath("XDG_CONFIG_HOME", ".config", configName)
}

// DefaultCachePath is $XDG_CACHE_HOME/athenaeum/m4b.json.
func DefaultCachePath() (string, error) {
	return xdgPath("XDG_CACHE_HOME", ".cache", cacheName)
}

// xdgPath resolves an XDG base directory by hand rather than through
// os.UserConfigDir/os.UserCacheDir, because those ignore the XDG variables on
// macOS and return the ~/Library equivalents instead.
func xdgPath(envVar, homeFallback, name string) (string, error) {
	base := os.Getenv(envVar)
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, homeFallback)
	}
	return filepath.Join(base, appName, name), nil
}
