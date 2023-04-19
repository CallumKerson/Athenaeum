package service

import (
	"context"
	"strings"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	allAudiobooksFeedTitle       = "Audiobooks"
	allAudiobooksFeedDescription = "Like movies in your mind!"
	descriptionFormat            = "%s Audiobooks"
	authorDescriptionFormat      = "Audiobooks by %s"
	narratorDescriptionFormat    = "Audiobooks Narrated by %s"
)

type AudiobooksClient interface {
	GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error)
	GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error)
	GetAudiobooksByAuthor(ctx context.Context, author string) (books []audiobooks.Audiobook, err error)
	GetAudiobooksByNarrator(ctx context.Context, narrator string) (books []audiobooks.Audiobook, err error)
	GetAudiobooksByTag(ctx context.Context, tag string) (books []audiobooks.Audiobook, err error)
	UpdateAudiobooks(ctx context.Context) error
}

type Service struct {
	log                     loggerrific.Logger
	host                    string
	feedImageLink           string
	mediaPath               string
	feedExplicit            bool
	feedLanguage            string
	feedAuthor              string
	feedAuthorEmail         string
	feedCopyright           string
	handlePreUnixEpochDates bool
	AudiobooksClient
}

func (s *Service) IsReady(ctx context.Context) bool {
	return true
}

func (s *Service) UpdateFeeds(ctx context.Context) error {
	return s.UpdateAudiobooks(ctx)
}

func New(audiobooksClient AudiobooksClient, host, mediaPath string, opts ...Option) *Service {
	svc := &Service{
		log:              noOpLogger.New(),
		AudiobooksClient: audiobooksClient,
		host:             host,
		mediaPath:        mediaPath,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *Service) getMediaPath() string {
	return strings.Trim(s.mediaPath, "/")
}

func (s *Service) getHost() string {
	return strings.Trim(s.host, "/")
}
