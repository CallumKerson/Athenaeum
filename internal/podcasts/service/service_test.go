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

	"github.com/CallumKerson/loggerrific/tlogger"

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

func TestGetFeed(t *testing.T) {
	tests := []struct {
		name               string
		writeFeedTest      func(*Service, io.Writer) error
		pathToExpectedFeed string
		expectedFeed       string
	}{
		{name: "Full feed", writeFeedTest: func(svc *Service, wrt io.Writer) error {
			return svc.WriteAllAudiobooksFeed(context.Background(), wrt)
		}, pathToExpectedFeed: "full_feed.rss"},
		{name: "Sci-Fi feed", writeFeedTest: func(svc *Service, wrt io.Writer) error {
			return svc.WriteGenreAudiobookFeed(context.Background(), audiobooks.SciFi, wrt)
		}, pathToExpectedFeed: "scifi_feed.rss"},
	}

	svc := New(&testAudiobookClient{},
		tlogger.NewTLogger(t),
		WithHost("http://www.example-podcast.com/audiobooks/"),
		WithMediaPath("/media/"),
		WithPodcastFeedInfo(true, "EN", "A Person", "person@domain.test", "None", "http://www.example-podcast.com/images/itunes.jpg"),
		WithHandlePreUnixEpoch(true),
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
			err := testCase.writeFeedTest(svc, &buf)

			// then
			assert.NoError(t, err)
			assert.Equal(t, expected, buf.String())
		})
	}
}
