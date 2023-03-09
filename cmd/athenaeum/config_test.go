package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"
)

func TestConfig_FromEnvironment(t *testing.T) {
	// given
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_HOST":       "http://localhost:8081",
		"ATHENAEUM_MEDIA_ROOT": "/mount/audiobooks",
	})
	t.Cleanup(envVarCleanup)
	viper.Reset()

	// when
	config, err := NewConfig(defaultPort, tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8081", config.Host)
	assert.Equal(t, "/mount/audiobooks", config.Media.Root)
}

func TestConfig_FromFile(t *testing.T) {
	// given
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_CONFIG_PATH": configFilePath,
	})
	t.Cleanup(envVarCleanup)
	viper.Reset()

	configYAML := `---
Host: "http://localhost:8082"
Podcast:
  Opts:
    Title: "Test Audiobooks"
    Description: "A Test Audiobook feed"
    Explicit: False
`
	err := os.WriteFile(configFilePath, []byte(configYAML), 0644)
	assert.NoError(t, err)

	// when
	config, err := NewConfig(defaultPort, tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8082", config.Host)
	assert.Equal(t, "Test Audiobooks", config.Podcast.Opts.Title)
	assert.Equal(t, "A Test Audiobook feed", config.Podcast.Opts.Description)
	assert.Equal(t, false, config.Podcast.Opts.Explicit)
	assert.Equal(t, "http://localhost:8082", config.Podcast.Opts.Link)
}

func TestConfig_EnvironmentOverridesFile(t *testing.T) {
	// given
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_CONFIG_PATH": configFilePath,
		"ATHENAEUM_HOST":        "http://127.0.0.1",
	})
	t.Cleanup(envVarCleanup)
	viper.Reset()

	configYAML := `---
Host: "http://localhost:8083"
`
	err := os.WriteFile(configFilePath, []byte(configYAML), 0644)
	assert.NoError(t, err)

	// when
	config, err := NewConfig(defaultPort, tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1", config.Host)
}

func TestConfig_DefaultsOnly(t *testing.T) {
	// given
	viper.Reset()

	// when
	config, err := NewConfig(defaultPort, tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "Audiobooks", config.Podcast.Opts.Title)
	assert.Equal(t, "Like movies for your mind!", config.Podcast.Opts.Description)
	assert.Equal(t, "None", config.Podcast.Opts.Copyright)
	assert.Equal(t, true, config.Podcast.Opts.Explicit)
	assert.Equal(t, "EN", config.Podcast.Opts.Language)
	assert.Equal(t, "media", config.Media.HostPath)
	assert.Equal(t, "/srv/media", config.Media.Root)
	assert.Equal(t, "http://localhost:8080", config.Host)
}

func TestConfig_BadFile(t *testing.T) {
	// given
	configFilePath := filepath.Join(t.TempDir(), "not-a-file.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_CONFIG_PATH": configFilePath,
	})
	t.Cleanup(envVarCleanup)
	viper.Reset()

	// when
	_, err := NewConfig(defaultPort, tlogger.NewTLogger(t))

	// then
	if assert.Error(t, err) {
		var expectedErr *fs.PathError
		assert.ErrorAs(t, err, &expectedErr)
	}
}

func envVarSetter(t *testing.T, envs map[string]string) (closer func()) {
	originalEnvVars := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvVars[name] = originalValue
		}
		err := os.Setenv(name, value)
		assert.NoError(t, err)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvVars[name]
			if has {
				t.Log("Setting env", name, "to", origValue)
				err := os.Setenv(name, origValue)
				assert.NoError(t, err)
			} else {
				t.Log("Unsetting env", name)
				err := os.Unsetenv(name)
				assert.NoError(t, err)
			}
		}
	}
}