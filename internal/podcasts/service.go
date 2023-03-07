package podcasts

import (
	"context"
	"errors"
	"time"

	"github.com/CallumKerson/loggerrific"
	"github.com/CallumKerson/podcasts"
)

type Service struct {
	Log loggerrific.Logger
	Cfg *FeedOpts
}

var errNoConfig = errors.New("no config")

func (s *Service) GetFeed(ctx context.Context) (string, error) {
	if s.Cfg == nil {
		return "", errNoConfig
	}
	pod := &podcasts.Podcast{
		Title:       s.Cfg.Title,
		Description: s.Cfg.Description,
		Language:    s.Cfg.Language,
		Link:        "http://www.example-podcast.com/my-podcast",
	}

	time1, err := time.Parse("2 Jan 2006 03:04PM", "10 Nov 2022 08:00AM")
	if err != nil {
		return "", err
	}

	pod.AddItem(&podcasts.Item{
		Title:    "Episode 1",
		GUID:     "http://www.example-podcast.com/my-podcast/1/episode-one",
		PubDate:  podcasts.NewPubDate(time1),
		Duration: podcasts.NewDuration(time.Second * 230),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/1/episode.mp3",
			Length: "12312",
			Type:   "MP3",
		},
	})

	time2, err := time.Parse("2 Jan 2006 03:04PM", "10 Nov 2022 11:00AM")
	if err != nil {
		return "", err
	}

	pod.AddItem(&podcasts.Item{
		Title:    "Episode 2",
		GUID:     "http://www.example-podcast.com/my-podcast/2/episode-two",
		PubDate:  podcasts.NewPubDate(time2),
		Duration: podcasts.NewDuration(time.Second * 320),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/2/episode.mp3",
			Length: "46732",
			Type:   "MP3",
		},
	})

	feed, err := pod.Feed(
		podcasts.Author(s.Cfg.Author),
		podcasts.Block,
		podcasts.Owner(s.Cfg.Author, s.Cfg.Email),
	)
	if err != nil {
		return "", nil
	}

	if s.Cfg.Explicit {
		err = feed.SetOptions(podcasts.Explicit)
		if err != nil {
			return "", nil
		}
	}

	return feed.XML()
}

func (s *Service) IsReady(ctx context.Context) (bool, error) {
	return true, nil
}

func NewService(opts *FeedOpts, logger loggerrific.Logger) *Service {
	return &Service{Log: logger, Cfg: opts}
}
