package service

import (
	"context"
	"fmt"
	"io"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	allAudiobooksFeedTitle       = "Audiobooks"
	allAudiobooksFeedDescription = "Like movies in your mind!"
	genreFeedDescriptionFormat   = "%s Audiobooks"
)

type AudiobooksClient interface {
	GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error)
	GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error)
}

type Service struct {
	Log                     loggerrific.Logger
	host                    string
	fedImageLink            string
	mediaPath               string
	feedExplicit            bool
	feedLanguage            string
	feedAuthor              string
	feedAuthorEmail         string
	feedCopyright           string
	handlePreUnixEpochDates bool
	AudiobooksClient
}

func (s *Service) WriteAllAudiobooksFeed(ctx context.Context, writer io.Writer) error {
	books, err := s.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	allAudiobooksFeedOpts := &FeedOpts{
		Title:       allAudiobooksFeedTitle,
		Description: allAudiobooksFeedDescription,
		Link:        s.host,
		ImageLink:   s.fedImageLink,
		Explicit:    s.feedExplicit,
		Language:    s.feedLanguage,
		Author:      s.feedAuthor,
		Email:       s.feedAuthorEmail,
		Copyright:   s.feedCopyright,
	}
	return s.WriteFeedFromAudiobooks(ctx, books, allAudiobooksFeedOpts, writer)
}

func (s *Service) WriteGenreAudiobookFeed(ctx context.Context, genre audiobooks.Genre, writer io.Writer) error {
	books, err := s.GetAudiobooksByGenre(ctx, genre)
	if err != nil {
		return err
	}
	allAudiobooksFeedOpts := &FeedOpts{
		Title:       genre.String(),
		Description: fmt.Sprintf(genreFeedDescriptionFormat, genre),
		Link:        s.host,
		ImageLink:   s.fedImageLink,
		Explicit:    s.feedExplicit,
		Language:    s.feedLanguage,
		Author:      s.feedAuthor,
		Email:       s.feedAuthorEmail,
		Copyright:   s.feedCopyright,
	}
	return s.WriteFeedFromAudiobooks(ctx, books, allAudiobooksFeedOpts, writer)
}

func (s *Service) IsReady(ctx context.Context) bool {
	return true
}

func New(audiobooksClient AudiobooksClient, logger loggerrific.Logger, opts ...Option) *Service {
	svc := &Service{
		Log:              logger,
		AudiobooksClient: audiobooksClient,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}
