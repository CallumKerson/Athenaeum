package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

type Config struct {
	Host    string
	DB      DB
	Media   Media
	Podcast Podcast
}

type DB struct {
	Root string
}

type Media struct {
	Root     string
	HostPath string
}

type Podcast struct {
	Root      string
	Copyright string
	Explicit  bool
	Language  string
	Author    string
	Email     string
}

func (c *Config) GetMediaHost() string {
	return fmt.Sprintf("%s/%s", c.Host, c.Media.HostPath)
}

func (c *Config) GetMediaServiceOpts() []mediaService.Option {
	return []mediaService.Option{mediaService.WithPathToMediaRoot(c.Media.Root)}
}

func (c *Config) GetBoltDBOps() []bolt.Option {
	return []bolt.Option{bolt.WithDBDefaults(), bolt.WithPathToDBDirectory(c.DB.Root)}
}

func (c *Config) GetPodcastServiceOpts() []podcastService.Option {
	return []podcastService.Option{podcastService.WithHost(c.Host),
		podcastService.WithMediaPath(c.Media.HostPath),
		podcastService.WithPodcastFeedInfo(c.Podcast.Explicit, c.Podcast.Language, c.Podcast.Author, c.Podcast.Email, c.Podcast.Copyright)}
}

func (c *Config) GetHTTPHandlerOpts() []transportHttp.HandlerOption {
	return []transportHttp.HandlerOption{transportHttp.WithMediaConfig(c.Media.Root, c.Media.HostPath), transportHttp.WithVersion(Version)}
}

func NewConfig(port int, logger loggerrific.Logger) (*Config, error) {
	viper.SetDefault("Podcast.Copyright", "None")
	viper.SetDefault("Podcast.Explicit", true)
	viper.SetDefault("Podcast.Language", "EN")
	viper.SetDefault("Podcast.Root", "/srv/podcasts")
	viper.SetDefault("Media.HostPath", "/media")
	viper.SetDefault("Media.Root", "/srv/media")
	viper.SetDefault("DB.Root", "/usr/local/athenaeum")
	viper.SetDefault("Host", fmt.Sprintf("http://localhost:%d", port))

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("athenaeum")
	_ = viper.BindEnv("Config.Path")

	_ = viper.BindEnv("Host")
	_ = viper.BindEnv("DB.Root")
	_ = viper.BindEnv("Media.Root")

	_ = viper.BindEnv("Podcast.Explicit")
	_ = viper.BindEnv("Podcast.Language")
	_ = viper.BindEnv("Podcast.Author")
	_ = viper.BindEnv("Podcast.Email")
	_ = viper.BindEnv("Podcast.Copyright")

	viper.AutomaticEnv()

	pathToConfig := viper.GetString("Config.Path")

	if !filepath.IsAbs(pathToConfig) {
		logger.Infoln("No valid config path found from environment variable ATHENAEUM_CONFIG_PATH,",
			"reading config from environment variables only")
	} else {
		viper.SetConfigFile(pathToConfig)
		err := viper.ReadInConfig()
		if err != nil {
			logger.WithError(err).Errorln("Cannot read config from file")
			return nil, err
		}
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}
