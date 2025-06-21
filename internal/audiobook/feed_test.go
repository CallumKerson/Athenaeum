package audiobook_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/audiobook"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func getTestFeedConfig() *audiobook.FeedConfig {
	return &audiobook.FeedConfig{
		Link:                    "http://www.example-podcast.com/audiobooks/",
		ImageLink:               "http://www.example-podcast.com/images/itunes.jpg",
		Explicit:                true,
		Language:                "EN",
		Author:                  "A Person",
		Email:                   "person@domain.test",
		Copyright:               "None",
		Host:                    "http://www.example-podcast.com/audiobooks/",
		MediaPath:               "/media/",
		HandlePreUnixEpochDates: true,
	}
}

func TestWriteFeedFromAudiobooks(t *testing.T) {
	config := getTestFeedConfig()
	config.Title = "Test Feed"
	config.Description = "Test Description"

	var buf bytes.Buffer
	err := audiobook.WriteFeedFromAudiobooks(testbooks.Audiobooks, config, &buf)

	assert.NoError(t, err)
	result := buf.String()
	assert.Contains(t, result, "Test Feed")
	assert.Contains(t, result, "Test Description")
	assert.Contains(t, result, testbooks.Audiobooks[0].Title)
}

func TestWriteFeedFromAudiobooks_NoConfig(t *testing.T) {
	var buf bytes.Buffer
	err := audiobook.WriteFeedFromAudiobooks(testbooks.Audiobooks, nil, &buf)

	assert.Error(t, err)
	assert.Equal(t, audiobook.ErrNoConfig, err)
}

func TestWriteAllAudiobooksFeed(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	err := audiobook.WriteAllAudiobooksFeed(testbooks.Audiobooks, config, &buf)

	assert.NoError(t, err)
	result := buf.String()
	assert.Contains(t, result, audiobook.AllAudiobooksFeedTitle)
	assert.Contains(t, result, audiobook.AllAudiobooksFeedDescription)
}

func TestWriteGenreAudiobookFeed(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	// Filter books to only include SciFi books
	scifiBooks := testbooks.AudiobooksFilteredBy(testbooks.GenreFilter(audiobooks.SciFi))
	err := audiobook.WriteGenreAudiobookFeed(scifiBooks, audiobooks.SciFi, config, &buf)

	assert.NoError(t, err)
	result := buf.String()
	assert.Contains(t, result, "Science Fiction")
	assert.Contains(t, result, "Science Fiction Audiobooks")
}

func TestWriteAuthorAudiobookFeed(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	// Filter books to only include books by Amal El-Mohtar
	authorBooks := testbooks.AudiobooksFilteredBy(testbooks.AuthorFilter("Amal El-Mohtar"))
	exists, err := audiobook.WriteAuthorAudiobookFeed(authorBooks, "Amal El-Mohtar", config, &buf)

	assert.NoError(t, err)
	assert.True(t, exists)
	result := buf.String()
	assert.Contains(t, result, "Amal El-Mohtar")
	assert.Contains(t, result, "Audiobooks by Amal El-Mohtar")
}

func TestWriteAuthorAudiobookFeed_NoBooks(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	exists, err := audiobook.WriteAuthorAudiobookFeed([]audiobooks.Audiobook{}, "Unknown Author", config, &buf)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Empty(t, buf.String())
}

func TestWriteNarratorAudiobookFeed(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	// Filter books to only include books narrated by Kobna Holdbrook-Smith
	narratorBooks := testbooks.AudiobooksFilteredBy(testbooks.NarratorFilter("Kobna Holdbrook-Smith"))
	exists, err := audiobook.WriteNarratorAudiobookFeed(narratorBooks, "Kobna Holdbrook-Smith", config, &buf)

	assert.NoError(t, err)
	assert.True(t, exists)
	result := buf.String()
	assert.Contains(t, result, "Kobna Holdbrook-Smith")
	assert.Contains(t, result, "Audiobooks Narrated by Kobna Holdbrook-Smith")
}

func TestWriteNarratorAudiobookFeed_NoBooks(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	exists, err := audiobook.WriteNarratorAudiobookFeed([]audiobooks.Audiobook{}, "Unknown Narrator", config, &buf)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Empty(t, buf.String())
}

func TestWriteTagAudiobookFeed(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	// Filter books to only include books with "Hugo Awards" tag
	tagBooks := testbooks.AudiobooksFilteredBy(testbooks.TagFilter("Hugo Awards"))
	exists, err := audiobook.WriteTagAudiobookFeed(tagBooks, "Hugo Awards", config, &buf)

	assert.NoError(t, err)
	assert.True(t, exists)
	result := buf.String()
	assert.Contains(t, result, "Hugo Awards")
	assert.Contains(t, result, "Hugo Awards Audiobooks")
}

func TestWriteTagAudiobookFeed_NoBooks(t *testing.T) {
	config := getTestFeedConfig()
	var buf bytes.Buffer

	exists, err := audiobook.WriteTagAudiobookFeed([]audiobooks.Audiobook{}, "Unknown Tag", config, &buf)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Empty(t, buf.String())
}

func TestGetFeed(t *testing.T) {
	config := getTestFeedConfig()
	config.Title = "String Feed"
	config.Description = "String Description"

	result, err := audiobook.GetFeed(testbooks.Audiobooks, config)

	assert.NoError(t, err)
	assert.Contains(t, result, "String Feed")
	assert.Contains(t, result, "String Description")
}

func TestWriteFeedFromAudiobooks_CompareWithExpected(t *testing.T) {
	// Test that matches the original podcast service behaviour
	config := &audiobook.FeedConfig{
		Title:                   audiobook.AllAudiobooksFeedTitle,
		Description:             audiobook.AllAudiobooksFeedDescription,
		Link:                    "http://www.example-podcast.com/audiobooks/",
		ImageLink:               "http://www.example-podcast.com/images/itunes.jpg",
		Explicit:                true,
		Language:                "EN",
		Author:                  "A Person",
		Email:                   "person@domain.test",
		Copyright:               "None",
		Host:                    "http://www.example-podcast.com/audiobooks",
		MediaPath:               "/media/",
		HandlePreUnixEpochDates: true,
	}

	var buf bytes.Buffer
	err := audiobook.WriteFeedFromAudiobooks(testbooks.Audiobooks, config, &buf)
	require.NoError(t, err)

	result := buf.String()

	// Check that it's valid RSS
	assert.Contains(t, result, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, result, `version="2.0"`)
	assert.Contains(t, result, `<channel>`)
	assert.Contains(t, result, `</channel>`)
	assert.Contains(t, result, `</rss>`)

	// Check feed metadata
	assert.Contains(t, result, audiobook.AllAudiobooksFeedTitle)
	assert.Contains(t, result, audiobook.AllAudiobooksFeedDescription)
	assert.Contains(t, result, "A Person")
	assert.Contains(t, result, "person@domain.test")

	// Check that audiobook items are included
	for _, book := range testbooks.Audiobooks {
		assert.Contains(t, result, book.Title)
		assert.Contains(t, result, book.GetAuthor())
	}
}
