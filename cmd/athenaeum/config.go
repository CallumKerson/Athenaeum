package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultPort      = 8080
	defaultConfigDir = ".athenaeum"
)

type Config struct {
	Host       string
	Port       int
	DB         DB
	Media      Media
	Podcast    Podcast
	ThirdParty ThirdParty
	Log        Log
}

type Log struct {
	Level string
}

type DB struct {
	Root string
}

type Media struct {
	Root     string
	HostPath string
}

type ThirdParty struct {
	// Deprecated
	UpdateOvercast bool
	NotifyOvercast bool
}

type Podcast struct {
	Copyright    string
	Explicit     bool
	Language     string
	Author       string
	Email        string
	ImagePath    string
	PreUnixEpoch PreUnixEpoch
}

type PreUnixEpoch struct {
	Handle bool
}

func (c *Config) GetLogLevel() string {
	return c.Log.Level
}

func InitConfig(cfg *Config, pathToConfigFile string, out io.Writer) error {
	viper.SetDefault("Podcast.Copyright", "None")
	viper.SetDefault("Podcast.Explicit", true)
	viper.SetDefault("Podcast.Language", "EN")
	viper.SetDefault("Podcast.PreUnixEpoch.Handle", true)
	viper.SetDefault("Media.HostPath", "/media")
	viper.SetDefault("Media.Root", "/srv/media")
	viper.SetDefault("ThirdParty.UpdateOvercast", false)
	viper.SetDefault("ThirdParty.NotifyOvercast", false)
	viper.SetDefault("Port", defaultPort)
	viper.SetDefault("Host", fmt.Sprintf("http://localhost:%d", defaultPort))
	viper.SetDefault("Log.Level", "INFO")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("athenaeum")
	_ = viper.BindEnv("Config.Path")

	_ = viper.BindEnv("Host")
	_ = viper.BindEnv("Log.Level")
	_ = viper.BindEnv("DB.Root")
	_ = viper.BindEnv("Media.Root")

	_ = viper.BindEnv("Podcast.Explicit")
	_ = viper.BindEnv("Podcast.Language")
	_ = viper.BindEnv("Podcast.Author")
	_ = viper.BindEnv("Podcast.Email")
	_ = viper.BindEnv("Podcast.Copyright")

	viper.AutomaticEnv()

	if pathToConfigFile == "" {
		pathToConfigFile = viper.GetString("Config.Path")
	}

	if pathToConfigFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		viper.AddConfigPath(filepath.Join(home, defaultConfigDir))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	} else {
		viper.SetConfigFile(pathToConfigFile)
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(out, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(out, "No valid config path provided by flag or found from environment variable ATHENAEUM_CONFIG_PATH,",
			"reading config from environment variables only")
	}
	if viper.GetString("DB.Root") == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		viper.SetDefault("DB.Root", filepath.Join(home, defaultConfigDir, "data"))
	}

	err := viper.Unmarshal(cfg)
	return err
}
