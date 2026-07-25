package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
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

// Site is where the generated tree is written.
type Site struct {
	Root string
}

// BuildConfig is everything `athenaeum build` needs. It reuses the server's
// config types so the two read the same keys during the migration.
type BuildConfig struct {
	Host                   string
	Media                  Media
	Site                   Site
	Podcast                Podcast
	ThirdParty             ThirdParty
	ExclusionsFromMainFeed ExclusionsFromMainFeed
}

// LoadBuildConfig reads the TOML config file, falling back to defaults for
// anything it does not set. A missing config file is not an error — every value
// can also come from a flag.
func LoadBuildConfig(pathToConfigFile string, out io.Writer) (*BuildConfig, error) {
	cfg := &BuildConfig{
		Host: "http://localhost:8080",
		Podcast: Podcast{
			Copyright:    "None",
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
