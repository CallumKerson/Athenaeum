package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type testAudiobookClient struct{}

func (c *testAudiobookClient) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	return testbooks.Audiobooks, nil
}

func (c *testAudiobookClient) GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error) {
	return testbooks.AudiobooksFilteredBy(testbooks.GenreFilter(genre)), nil
}

func (c *testAudiobookClient) GetAudiobooksByAuthor(ctx context.Context, author string) ([]audiobooks.Audiobook, error) {
	return testbooks.AudiobooksFilteredBy(testbooks.AuthorFilter(author)), nil
}

func (c *testAudiobookClient) GetAudiobooksByNarrator(ctx context.Context, narrator string) ([]audiobooks.Audiobook, error) {
	return testbooks.AudiobooksFilteredBy(testbooks.NarratorFilter(narrator)), nil
}

func (c *testAudiobookClient) GetAudiobooksByTag(ctx context.Context, tag string) ([]audiobooks.Audiobook, error) {
	return testbooks.AudiobooksFilteredBy(testbooks.TagFilter(tag)), nil
}

func (c *testAudiobookClient) UpdateAudiobooks(ctx context.Context) error {
	return nil
}

func TestGetFeed(t *testing.T) {
	tests := []struct {
		name               string
		writeFeedTest      func(*Service, io.Writer) (bool, error)
		pathToExpectedFeed string
		expectedFeed       string
		expectedFeedExists bool
	}{
		{name: "Feed", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return true, svc.WriteAllAudiobooksFeed(context.Background(), wrt)
		}, pathToExpectedFeed: "feed.rss", expectedFeedExists: true},
		{name: "Sci-Fi feed", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return true, svc.WriteGenreAudiobookFeed(context.Background(), audiobooks.SciFi, wrt)
		}, pathToExpectedFeed: "scifi_feed.rss", expectedFeedExists: true},
		{name: "Amal El-Mohtar Author feed", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteAuthorAudiobookFeed(context.Background(), "Amal El-Mohtar", wrt)
		}, pathToExpectedFeed: "el_mohtar_feed.rss", expectedFeedExists: true},
		{name: "Feed for author that does not exist in library", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteAuthorAudiobookFeed(context.Background(), "Octavia Butler", wrt)
		}, expectedFeed: "", expectedFeedExists: false},
		{name: "Kobna Holdbrook-Smith Narrator feed", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteNarratorAudiobookFeed(context.Background(), "Kobna Holdbrook-Smith", wrt)
		}, pathToExpectedFeed: "holdbrook_smith_feed.rss", expectedFeedExists: true},
		{name: "Feed for narrator that does not exist in library", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteNarratorAudiobookFeed(context.Background(), "Simon Vance", wrt)
		}, expectedFeed: "", expectedFeedExists: false},
		{name: "Tag feed", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteTagAudiobookFeed(context.Background(), "Hugo Awards", wrt)
		}, pathToExpectedFeed: "hugo_awards_feed.rss", expectedFeedExists: true},
		{name: "Feed for tag that does not exist in library", writeFeedTest: func(svc *Service, wrt io.Writer) (bool, error) {
			return svc.WriteNarratorAudiobookFeed(context.Background(), "Nebula Awards", wrt)
		}, expectedFeed: "", expectedFeedExists: false},
	}

	svc := New(&testAudiobookClient{},
		"http://www.example-podcast.com/audiobooks/",
		"/media/",
	)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			expected := testCase.expectedFeed
			if testCase.pathToExpectedFeed != "" {
				expectedBytes, err := os.ReadFile(filepath.Join("testdata", testCase.pathToExpectedFeed))
				assert.NoError(t, err)
				expected = strings.TrimSpace(string(expectedBytes))
			}
			var buf bytes.Buffer

			// when
			feedExists, err := testCase.writeFeedTest(svc, &buf)

			// then
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedFeedExists, feedExists)
			assert.Equal(t, expected, buf.String())
		})
	}
}

func TestGetFeed_WithOptions(t *testing.T) {
	// given
	svc := New(&testAudiobookClient{},
		"http://www.example-podcast.com/audiobooks",
		"/media/",
		WithPodcastFeedInfo(true, "EN", "A Person", "person@domain.test", "None", "http://www.example-podcast.com/images/itunes.jpg"),
		WithHandlePreUnixEpoch(true),
	)

	expectedBytes, err := os.ReadFile(filepath.Join("testdata", "full_feed.rss"))
	assert.NoError(t, err)
	expected := strings.TrimSpace(string(expectedBytes))
	var buf bytes.Buffer

	// when
	err = svc.WriteAllAudiobooksFeed(context.Background(), &buf)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
