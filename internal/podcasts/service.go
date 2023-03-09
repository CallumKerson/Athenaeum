package podcasts

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"

	"github.com/CallumKerson/loggerrific"
	"github.com/CallumKerson/podcasts"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

type AudiobooksClient interface {
	GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error)
}

type Service struct {
	Log        loggerrific.Logger
	Cfg        *FeedOpts
	HostPrefix string
	AudiobooksClient
}

var errNoConfig = errors.New("no config")

func (s *Service) WriteFeedFromAudiobooks(ctx context.Context, books []audiobooks.Audiobook, writer io.Writer) error {
	if s.Cfg == nil {
		return errNoConfig
	}
	pod := &podcasts.Podcast{
		Title:       s.Cfg.Title,
		Description: s.Cfg.Description,
		Language:    s.Cfg.Language,
		Link:        s.Cfg.Link,
	}

	for bookIndex := range books {
		hostedFile := fmt.Sprintf("%s%s", s.HostPrefix, books[bookIndex].Path)
		pod.AddItem(&podcasts.Item{
			Title:       books[bookIndex].Title,
			Description: &podcasts.CDATAText{Value: fmt.Sprintf("%s by %s", books[bookIndex].Title, books[bookIndex].GetAuthor())},
			PubDate:     podcasts.NewPubDate(books[bookIndex].ReleaseDate.AsTime(time.UTC).Add(8 * time.Hour)),
			Duration:    podcasts.NewDuration(books[bookIndex].Duration),
			GUID:        hostedFile,
			Enclosure: &podcasts.Enclosure{
				URL:    hostedFile,
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

	if s.Cfg.Author != "" {
		err = feed.SetOptions(podcasts.Author(s.Cfg.Author), podcasts.Owner(s.Cfg.Author, s.Cfg.Email))
		if err != nil {
			return err
		}
	}

	if s.Cfg.Explicit {
		err = feed.SetOptions(podcasts.Explicit)
		if err != nil {
			return err
		}
	}
	return feed.Write(writer)
}

func (s *Service) WriteFeed(ctx context.Context, writer io.Writer) error {
	books, err := s.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	return s.WriteFeedFromAudiobooks(ctx, books, writer)
}

func (s *Service) GetFeed(ctx context.Context) (string, error) {
	var buf bytes.Buffer
	if err := s.WriteFeed(ctx, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Service) IsReady(ctx context.Context) (bool, error) {
	return true, nil
}

func NewService(hostPrefix string, opts *FeedOpts, audiobooksClient AudiobooksClient, logger loggerrific.Logger) *Service {
	return &Service{
		HostPrefix:       hostPrefix,
		Log:              logger,
		Cfg:              opts,
		AudiobooksClient: audiobooksClient,
	}
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
