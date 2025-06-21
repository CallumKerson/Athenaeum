package audiobook

import (
	"bytes"
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

const (
	AllAudiobooksFeedTitle       = "Audiobooks"
	AllAudiobooksFeedDescription = "Like movies in your mind!"
	DescriptionFormat            = "%s Audiobooks"
	AuthorDescriptionFormat      = "Audiobooks by %s"
	NarratorDescriptionFormat    = "Audiobooks Narrated by %s"
)

var (
	ErrNoConfig = errors.New("no config")
	unixEpoch   = time.Unix(0, 0).UTC()
)

// FeedConfig contains configuration for generating podcast feeds
type FeedConfig struct {
	Title                   string
	Description             string
	Link                    string
	ImageLink               string
	Explicit                bool
	Language                string
	Author                  string
	Email                   string
	Copyright               string
	Host                    string
	MediaPath               string
	HandlePreUnixEpochDates bool
}

// WriteFeedFromAudiobooks generates a podcast feed from a slice of audiobooks
func WriteFeedFromAudiobooks(books []audiobooks.Audiobook, config *FeedConfig, writer io.Writer) error {
	if config == nil {
		return ErrNoConfig
	}

	pod := &podcasts.Podcast{
		Title:       config.Title,
		Description: config.Description,
		Language:    config.Language,
		Link:        config.Link,
	}

	host := strings.Trim(config.Host, "/")
	mediaPath := strings.Trim(config.MediaPath, "/")

	for bookIndex := range books {
		hostedFile, err := url.Parse(fmt.Sprintf("%s/%s%s", host, mediaPath, books[bookIndex].Path))
		if err != nil {
			return err
		}

		pubDate := books[bookIndex].ReleaseDate.AsTime(time.UTC)
		if config.HandlePreUnixEpochDates && unixEpoch.After(pubDate) {
			pubDate = unixEpoch
		}
		pubDate = pubDate.Add(8 * time.Hour)

		pod.AddItem(&podcasts.Item{
			Title:       books[bookIndex].Title,
			Description: &podcasts.CDATAText{Value: fmt.Sprintf("%s by %s", books[bookIndex].Title, books[bookIndex].GetAuthor())},
			PubDate:     podcasts.NewPubDate(pubDate),
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

	if config.Author != "" {
		err = feed.SetOptions(podcasts.Author(config.Author), podcasts.Owner(config.Author, config.Email))
		if err != nil {
			return err
		}
	}

	if config.Explicit {
		err = feed.SetOptions(podcasts.Explicit)
		if err != nil {
			return err
		}
	}

	if config.ImageLink != "" {
		err = feed.SetOptions(podcasts.Image(config.ImageLink))
		if err != nil {
			return err
		}
	}
	return feed.Write(writer)
}

// WriteAllAudiobooksFeed generates a feed with all audiobooks
func WriteAllAudiobooksFeed(books []audiobooks.Audiobook, config *FeedConfig, writer io.Writer) error {
	feedConfig := &FeedConfig{
		Title:                   AllAudiobooksFeedTitle,
		Description:             AllAudiobooksFeedDescription,
		Link:                    config.Link,
		ImageLink:               config.ImageLink,
		Explicit:                config.Explicit,
		Language:                config.Language,
		Author:                  config.Author,
		Email:                   config.Email,
		Copyright:               config.Copyright,
		Host:                    config.Host,
		MediaPath:               config.MediaPath,
		HandlePreUnixEpochDates: config.HandlePreUnixEpochDates,
	}
	return WriteFeedFromAudiobooks(books, feedConfig, writer)
}

// WriteGenreAudiobookFeed generates a feed for a specific genre
func WriteGenreAudiobookFeed(books []audiobooks.Audiobook, genre audiobooks.Genre, config *FeedConfig, writer io.Writer) error {
	feedConfig := &FeedConfig{
		Title:                   genre.String(),
		Description:             fmt.Sprintf(DescriptionFormat, genre),
		Link:                    config.Link,
		ImageLink:               config.ImageLink,
		Explicit:                config.Explicit,
		Language:                config.Language,
		Author:                  config.Author,
		Email:                   config.Email,
		Copyright:               config.Copyright,
		Host:                    config.Host,
		MediaPath:               config.MediaPath,
		HandlePreUnixEpochDates: config.HandlePreUnixEpochDates,
	}
	return WriteFeedFromAudiobooks(books, feedConfig, writer)
}

// writeFilteredFeed is a helper function to generate feeds for specific criteria
func writeFilteredFeed(books []audiobooks.Audiobook, title, desc string, config *FeedConfig, writer io.Writer) (bool, error) {
	if len(books) < 1 {
		return false, nil
	}
	feedConfig := &FeedConfig{
		Title:                   title,
		Description:             desc,
		Link:                    config.Link,
		ImageLink:               config.ImageLink,
		Explicit:                config.Explicit,
		Language:                config.Language,
		Author:                  config.Author,
		Email:                   config.Email,
		Copyright:               config.Copyright,
		Host:                    config.Host,
		MediaPath:               config.MediaPath,
		HandlePreUnixEpochDates: config.HandlePreUnixEpochDates,
	}
	return true, WriteFeedFromAudiobooks(books, feedConfig, writer)
}

// WriteAuthorAudiobookFeed generates a feed for a specific author
func WriteAuthorAudiobookFeed(books []audiobooks.Audiobook, author string, config *FeedConfig, writer io.Writer) (bool, error) {
	return writeFilteredFeed(books, author, fmt.Sprintf(AuthorDescriptionFormat, author), config, writer)
}

// WriteNarratorAudiobookFeed generates a feed for a specific narrator
func WriteNarratorAudiobookFeed(books []audiobooks.Audiobook, narrator string, config *FeedConfig, writer io.Writer) (bool, error) {
	return writeFilteredFeed(books, narrator, fmt.Sprintf(NarratorDescriptionFormat, narrator), config, writer)
}

// WriteTagAudiobookFeed generates a feed for a specific tag
func WriteTagAudiobookFeed(books []audiobooks.Audiobook, tag string, config *FeedConfig, writer io.Writer) (bool, error) {
	return writeFilteredFeed(books, tag, fmt.Sprintf(DescriptionFormat, tag), config, writer)
}

// GetFeed generates a feed and returns it as a string
func GetFeed(books []audiobooks.Audiobook, config *FeedConfig) (string, error) {
	var buf bytes.Buffer
	if err := WriteFeedFromAudiobooks(books, config, &buf); err != nil {
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
