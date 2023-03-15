package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"

	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	defaultPort = 8080
)

type Config struct {
	Host    string
	Port    int
	DB      DB
	Media   Media
	Podcast Podcast
	Log     Log
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

type Podcast struct {
	Copyright string
	Explicit  bool
	Language  string
	Author    string
	Email     string
}

func (c *Config) GetLogLevel() string {
	return c.Log.Level
}

func (c *Config) GetHTTPHandlerOpts() []transportHttp.HandlerOption {
	return []transportHttp.HandlerOption{transportHttp.WithMediaConfig(c.Media.Root, c.Media.HostPath), transportHttp.WithVersion(Version)}
}

func InitConfig(cfg *Config, pathToConfigFile string, out io.Writer) error {
	viper.SetDefault("Podcast.Copyright", "None")
	viper.SetDefault("Podcast.Explicit", true)
	viper.SetDefault("Podcast.Language", "EN")
	viper.SetDefault("Media.HostPath", "/media")
	viper.SetDefault("Media.Root", "/srv/media")
	viper.SetDefault("DB.Root", "/usr/local/athenaeum")
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
		fmt.Fprintln(out, "No valid config path found from environment variable ATHENAEUM_CONFIG_PATH,",
			"reading config from environment variables only")
	} else {
		viper.SetConfigFile(pathToConfigFile)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintln(out, "Cannot read config from file:", err)
			return err
		}
	}

	err := viper.Unmarshal(cfg)
	return err
}
