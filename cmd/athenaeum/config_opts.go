package main

import (
	"fmt"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
)

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
