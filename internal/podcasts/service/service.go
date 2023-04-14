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
	return s.writeAudiobookFeed(ctx, allAudiobooksFeedTitle, allAudiobooksFeedDescription, books, writer)
}

func (s *Service) WriteGenreAudiobookFeed(ctx context.Context, genre audiobooks.Genre, writer io.Writer) error {
	books, err := s.GetAudiobooksByGenre(ctx, genre)
	if err != nil {
		return err
	}
	genreFeedOpts := &FeedOpts{
		Title:       genre.String(),
		Description: fmt.Sprintf(descriptionFormat, genre),
		Link:        s.host,
		ImageLink:   s.fedImageLink,
		Explicit:    s.feedExplicit,
		Language:    s.feedLanguage,
		Author:      s.feedAuthor,
		Email:       s.feedAuthorEmail,
		Copyright:   s.feedCopyright,
	}
	return s.WriteFeedFromAudiobooks(ctx, books, genreFeedOpts, writer)
}

func (s *Service) WriteAuthorAudiobookFeed(ctx context.Context, author string, writer io.Writer) (bool, error) {
	books, err := s.GetAudiobooksByAuthor(ctx, author)
	if err != nil {
		return false, err
	}
	if len(books) < 1 {
		return false, nil
	}
	return true, s.writeAudiobookFeed(ctx, author, fmt.Sprintf(authorDescriptionFormat, author), books, writer)
}

func (s *Service) WriteNarratorAudiobookFeed(ctx context.Context, narrator string, writer io.Writer) (bool, error) {
	books, err := s.GetAudiobooksByNarrator(ctx, narrator)
	if err != nil {
		return false, err
	}
	if len(books) < 1 {
		return false, nil
	}
	return true, s.writeAudiobookFeed(ctx, narrator, fmt.Sprintf(narratorDescriptionFormat, narrator), books, writer)
}

func (s *Service) WriteTagAudiobookFeed(ctx context.Context, tag string, writer io.Writer) (bool, error) {
	books, err := s.GetAudiobooksByTag(ctx, tag)
	if err != nil {
		return false, err
	}
	if len(books) < 1 {
		return false, nil
	}
	return true, s.writeAudiobookFeed(ctx, tag, fmt.Sprintf(descriptionFormat, tag), books, writer)
}

func (s *Service) IsReady(ctx context.Context) bool {
	return true
}

func (s *Service) UpdateFeeds(ctx context.Context) error {
	return s.UpdateAudiobooks(ctx)
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

func (s *Service) writeAudiobookFeed(ctx context.Context, title, description string,
	books []audiobooks.Audiobook, writer io.Writer) error {
	feedOtps := &FeedOpts{
		Title:       title,
		Description: description,
		Link:        s.host,
		ImageLink:   s.fedImageLink,
		Explicit:    s.feedExplicit,
		Language:    s.feedLanguage,
		Author:      s.feedAuthor,
		Email:       s.feedAuthorEmail,
		Copyright:   s.feedCopyright,
	}
	return s.WriteFeedFromAudiobooks(ctx, books, feedOtps, writer)
}
