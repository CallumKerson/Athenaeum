package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"

	"github.com/CallumKerson/podcasts"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

var errNoConfig = errors.New("no config")

func (s *Service) WriteFeedFromAudiobooks(ctx context.Context, books []audiobooks.Audiobook, feedOpts *FeedOpts, writer io.Writer) error {
	if feedOpts == nil {
		return errNoConfig
	}
	pod := &podcasts.Podcast{
		Title:       feedOpts.Title,
		Description: feedOpts.Description,
		Language:    feedOpts.Language,
		Link:        feedOpts.Link,
	}

	for bookIndex := range books {
		hostedFile, err := url.Parse(fmt.Sprintf("%s/%s%s", s.host, s.mediaPath, books[bookIndex].Path))
		if err != nil {
			return err
		}

		pod.AddItem(&podcasts.Item{
			Title:       books[bookIndex].Title,
			Description: &podcasts.CDATAText{Value: fmt.Sprintf("%s by %s", books[bookIndex].Title, books[bookIndex].GetAuthor())},
			PubDate:     podcasts.NewPubDate(books[bookIndex].ReleaseDate.AsTime(time.UTC).Add(8 * time.Hour)),
			Duration:    podcasts.NewDuration(books[bookIndex].Duration),
			GUID:        hostedFile.String(),
			Enclosure: &podcasts.Enclosure{
				URL:    hostedFile.String(),
				Length: fmt.Sprintf("%d", books[bookIndex].FileSize),
				Type:   books[bookIndex].MIMEType,
			},
			Subtitle:       books[bookIndex].GetAuthor(),
			ContentEncoded: &podcasts.CDATAText{Value: getHTMLSummary(&books[bookIndex])},
		})
	}

	feed, err := pod.Feed(podcasts.Block)
	if err != nil {
		return err
	}

	if feedOpts.Author != "" {
		err = feed.SetOptions(podcasts.Author(feedOpts.Author), podcasts.Owner(feedOpts.Author, feedOpts.Email))
		if err != nil {
			return err
		}
	}

	if feedOpts.Explicit {
		err = feed.SetOptions(podcasts.Explicit)
		if err != nil {
			return err
		}
	}

	if feedOpts.ImageLink != "" {
		err = feed.SetOptions(podcasts.Image(feedOpts.ImageLink))
		if err != nil {
			return err
		}
	}
	return feed.Write(writer)
}

func (s *Service) WriteFeed(ctx context.Context, feedOpts *FeedOpts, writer io.Writer) error {
	books, err := s.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	return s.WriteFeedFromAudiobooks(ctx, books, feedOpts, writer)
}

func (s *Service) GetFeed(ctx context.Context, feedOpts *FeedOpts) (string, error) {
	var buf bytes.Buffer
	if err := s.WriteFeed(ctx, feedOpts, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getHTMLSummary(book *audiobooks.Audiobook) string {
	var builder strings.Builder
	_, _ = builder.WriteString(fmt.Sprintf("<h1>%s</h1>", book.Title))
	if book.Subtitle != "" {
		_, _ = builder.WriteString(fmt.Sprintf("<h4>%s</h4>", book.Title))
	}
	_, _ = builder.WriteString(fmt.Sprintf("<h2>By %s</h2>", book.GetAuthor()))
	if book.Series != nil {
		_, _ = builder.WriteString(fmt.Sprintf("<h4>%s Book %v</h4>", book.Series.Title, book.Series.Sequence))
	}
	if book.GetNarrator() != "" {
		_, _ = builder.WriteString(fmt.Sprintf("<h4>Narrated by %s</h4>", book.GetNarrator()))
	}

	if book.Description != nil {
		switch book.Description.Format {
		case description.HTML:
			_, _ = builder.WriteString(book.Description.Text)
		case description.Markdown:
			md := []byte(book.Description.Text)
			_, _ = builder.WriteString(string(markdown.ToHTML(md, nil, nil)))
		case description.Plain, description.Undefined:
			lines := strings.Split(book.Description.Text, "\n")
			for _, line := range lines {
				_, _ = builder.WriteString(fmt.Sprintf("<p>%s</p>", line))
			}
		}
	}
	return builder.String()
}
