package config

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/podcasts"
)

type Config struct {
	Host    string
	Media   Media
	Podcast Podcast
}

type Media struct {
	Root     string
	HostPath string
}

type Podcast struct {
	Root string
	Opts podcasts.FeedOpts
}

func (c *Config) GetMediaHost() string {
	return path.Join(c.Host, c.Media.HostPath)
}

func New(logger loggerrific.Logger) (*Config, error) {
	viper.SetDefault("Podcast.Opts.Title", "Audiobooks")
	viper.SetDefault("Podcast.Opts.Description", "Like movies for your mind!")
	viper.SetDefault("Podcast.Opts.Copyright", "None")
	viper.SetDefault("Podcast.Opts.Explicit", true)
	viper.SetDefault("Podcast.Opts.Language", "EN")
	viper.SetDefault("Podcast.Root", "/srv/podcasts")
	viper.SetDefault("Media.HostPath", "media")
	viper.SetDefault("Media.Root", "/srv/media")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("athenaeum")
	_ = viper.BindEnv("Config.Path")

	_ = viper.BindEnv("Host")
	_ = viper.BindEnv("Media.Root")

	_ = viper.BindEnv("Podcast.Opts.Title")
	_ = viper.BindEnv("Podcast.Opts.Description")
	_ = viper.BindEnv("Podcast.Opts.Explicit")
	_ = viper.BindEnv("Podcast.Opts.Language")
	_ = viper.BindEnv("Podcast.Opts.Author")
	_ = viper.BindEnv("Podcast.Opts.Email")
	_ = viper.BindEnv("Podcast.Opts.Copyright")

	viper.RegisterAlias("Podcast.Opts.Link", "Host")

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
