package main

import (
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Address            string
	LogLevel           string
	LogFormat          string
	staticDir          string
	templateDir        string
	GracefulTimeout    time.Duration
	DBLocation         string
	libraryHome        string
	APIUser            string
	APIPassword        string
	host               string
	podcastAuthorName  string
	podcastAuthorEmail string
}

func (c *config) StaticDir() string {
	return c.staticDir
}

func (c *config) TemplateDir() string {
	return c.templateDir
}

func (c *config) LibraryHome() string {
	return c.libraryHome
}

func (c *config) Host() string {
	return c.host
}

func (c *config) PodcastAuthorEmail() string {
	return c.podcastAuthorEmail
}

func (c *config) PodcastAuthorName() string {
	return c.podcastAuthorName
}

func loadConfig() config {
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("ADDRESS", ":8030")
	viper.SetDefault("STATIC_DIR", "./resources/static/")
	viper.SetDefault("TEMPLATE_DIR", "./resources/templates/")
	viper.SetDefault("GRACEFUL_TIMEOUT", 20*time.Second)
	viper.SetDefault("DB_LOCATION", "/Users/ckerson/Music/exampledb/")
	viper.SetDefault("LIBRARY_HOME", "/Users/ckerson/Music/samplebooks/")
	viper.SetDefault("API_USERNAME", "librarian")
	viper.SetDefault("API_PASSWORD", "a9ba69d1-fa46-485c-b410-e4055467210f")
	viper.SetDefault("HOST", "http://localhost"+viper.GetString("ADDRESS"))
	viper.SetDefault("PODCAST_AUTHOR_NAME", "Podcast Author")
	viper.SetDefault("PODCAST_AUTHOR_EMIAL", "demo@example.demo")

	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	return config{
		// application configuration
		Address:            viper.GetString("ADDRESS"),
		staticDir:          viper.GetString("STATIC_DIR"),
		templateDir:        viper.GetString("TEMPLATE_DIR"),
		LogLevel:           viper.GetString("LOG_LEVEL"),
		LogFormat:          viper.GetString("LOG_FORMAT"),
		GracefulTimeout:    viper.GetDuration("GRACEFUL_TIMEOUT"),
		DBLocation:         viper.GetString("DB_LOCATION"),
		libraryHome:        viper.GetString("LIBRARY_HOME"),
		APIUser:            viper.GetString("API_USERNAME"),
		APIPassword:        viper.GetString("API_PASSWORD"),
		host:               viper.GetString("HOST"),
		podcastAuthorName:  viper.GetString("PODCAST_AUTHOR_NAME"),
		podcastAuthorEmail: viper.GetString("PODCAST_AUTHOR_EMAIL"),
	}
}
