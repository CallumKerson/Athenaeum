package service

import (
	"context"
	"fmt"
	"io"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

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
		Link:        s.getHost(),
		ImageLink:   s.feedImageLink,
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

func (s *Service) writeAudiobookFeed(ctx context.Context, title, description string,
	books []audiobooks.Audiobook, writer io.Writer) error {
	feedOtps := &FeedOpts{
		Title:       title,
		Description: description,
		Link:        s.getHost(),
		ImageLink:   s.feedImageLink,
		Explicit:    s.feedExplicit,
		Language:    s.feedLanguage,
		Author:      s.feedAuthor,
		Email:       s.feedAuthorEmail,
		Copyright:   s.feedCopyright,
	}
	return s.WriteFeedFromAudiobooks(ctx, books, feedOtps, writer)
}
