package service

import (
	"context"
	"io"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	allAudiobooksFeedTitle       = "Audiobooks"
	allAudiobooksFeedDescription = "Like movies for your mind!"
)

type AudiobooksClient interface {
	GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error)
}

type Service struct {
	Log             loggerrific.Logger
	host            string
	mediaPath       string
	feedExplicit    bool
	feedLanguage    string
	feedAuthor      string
	feedAuthorEmail string
	feedCopyright   string
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
