package main

import (
	"fmt"
	"strings"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	"github.com/CallumKerson/Athenaeum/internal/memcache"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
	"github.com/CallumKerson/Athenaeum/pkg/client"
)

func (c *Config) GetMediaServiceOpts() []mediaService.Option {
	return []mediaService.Option{mediaService.WithPathToMediaRoot(c.Media.Root)}
}

func (c *Config) GetBoltDBOps() []bolt.Option {
	return []bolt.Option{bolt.WithDBDefaults(), bolt.WithPathToDBDirectory(c.DB.Root)}
}

func (c *Config) GetPodcastServiceOpts() []podcastService.Option {
	return []podcastService.Option{podcastService.WithHost(c.Host),
		podcastService.WithMediaPath(transportHttp.MediaPath),
		podcastService.WithPodcastFeedInfo(c.Podcast.Explicit, c.Podcast.Language, c.Podcast.Author, c.Podcast.Email, c.Podcast.Copyright,
			fmt.Sprintf("%s%s/itunes_image.jpg", c.Host, transportHttp.StaticPath)),
		podcastService.WithHandlePreUnixEpoch(c.Podcast.PreUnixEpoch.Handle)}
}

func (c *Config) GetClientOpts() []client.Option {
	return []client.Option{client.WithHost(c.Host), client.WithVersion(Version)}
}

func (c *Config) GetMemcacheOpts() []memcache.Option {
	ttl, _ := c.Cache.GetTTL()
	return []memcache.Option{memcache.WithTTL(ttl), memcache.WithCapacity(c.Cache.Length)}
}

func (c *Config) GetLogger() *logrus.Logger {
	log := logrus.NewLogger()
	level := strings.ToLower(c.Log.Level)
	switch level {
	case "debug":
		log.SetLevelDebug()
	case "warn":
		log.SetLevelWarn()
	case "error":
		log.SetLevelError()
	default:
		log.SetLevelInfo()
	}
	return log
}
