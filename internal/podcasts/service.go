package podcasts

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"

	"github.com/CallumKerson/loggerrific"
	"github.com/CallumKerson/podcasts"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
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

func (s *Service) GetFeed(ctx context.Context) (string, error) {
	if s.Cfg == nil {
		return "", errNoConfig
	}
	pod := &podcasts.Podcast{
		Title:       s.Cfg.Title,
		Description: s.Cfg.Description,
		Language:    s.Cfg.Language,
		Link:        s.Cfg.Link,
	}

	books, err := s.GetAllAudiobooks(ctx)
	if err != nil {
		return "", nil
	}

	for bookIndex := range books {
		hostedFile := fmt.Sprintf("%s%s", s.HostPrefix, books[bookIndex].Path)
		pod.AddItem(&podcasts.Item{
			Title:       books[bookIndex].Title,
			Description: &podcasts.CDATAText{Value: fmt.Sprintf("%s by %s", books[bookIndex].Title, books[bookIndex].GetAuthor())},
			PubDate:     podcasts.NewPubDate(books[bookIndex].ReleaseDate.Time),
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
		return "", nil
	}

	if s.Cfg.Author != "" {
		err = feed.SetOptions(podcasts.Author(s.Cfg.Author), podcasts.Owner(s.Cfg.Author, s.Cfg.Email))
		if err != nil {
			return "", nil
		}
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
	if book.Description != nil {
		switch book.Description.Format {
		case audiobooks.HTML:
			_, _ = builder.WriteString(book.Description.Text)
		case audiobooks.Markdown:
			md := []byte(book.Description.Text)
			_, _ = builder.WriteString(string(markdown.ToHTML(md, nil, nil)))
		case audiobooks.Plain, audiobooks.Undefined:
			lines := strings.Split(book.Description.Text, "\n")
			for _, line := range lines {
				_, _ = builder.WriteString(fmt.Sprintf("<p>%s</p>", line))
			}
		}
	}
	return builder.String()
}
