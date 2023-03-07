package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/podcasts"
)

type Config struct {
	Host    string
	Podcast podcasts.FeedOpts
}

func New(logger loggerrific.Logger) (Config, error) {
	viper.SetDefault("Podcast.Title", "Audiobooks")
	viper.SetDefault("Podcast.Description", "Like movies for your mind!")
	viper.SetDefault("Podcast.Copyright", "None")
	viper.SetDefault("Podcast.Explicit", true)
	viper.SetDefault("Podcast.Language", "EN")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("athenaeum")
	_ = viper.BindEnv("Config.Path")

	_ = viper.BindEnv("Host")

	_ = viper.BindEnv("Podcast.Title")
	_ = viper.BindEnv("Podcast.Description")
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
			return Config{}, err
		}
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return cfg, err
}
