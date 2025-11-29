package http

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	testFeed   = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\"></rss>"
	feedFormat = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\">" +
		"<channel><title>%s</title></channel></rss>"
	testAuthor   = "Agatha Test-y"
	testNarrator = "David Crochet"
	testTag      = "Cosy Classics"
)

func TestHandler(t *testing.T) {
	testHandler := NewHandler(&DummyPodcastService{}, filepath.Join("testdata", "media"), WithVersion("1.0.0-test"))

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	expectedIndex, err := os.ReadFile(filepath.Join("testdata", "expected-index.html"))
	assert.NoError(t, err)

	tests := []struct {
		name                string
		method              string
		path                string
		expectedStatus      int
		expectedContentType string
		expectedBody        string
	}{
		{
			name:                "index",
			method:              "GET",
			path:                "",
			expectedStatus:      200,
			expectedContentType: ContentTypeHTML,
			expectedBody:        strings.TrimSpace(string(expectedIndex)),
		},
		{
			name:                "health check",
			method:              "GET",
			path:                "/health",
			expectedStatus:      200,
			expectedContentType: ContentTypeJSON,
			expectedBody:        "{\n  \"health\": \"ok\"\n}",
		},
		{
			name:                "readiness check",
			method:              "GET",
			path:                "/ready",
			expectedStatus:      200,
			expectedContentType: ContentTypeJSON,
			expectedBody:        "{\n  \"readiness\": \"ok\"\n}",
		},
		{
			name:                "version",
			method:              "GET",
			path:                "/version",
			expectedStatus:      200,
			expectedContentType: ContentTypeJSON,
			expectedBody:        "{\n  \"version\": \"1.0.0-test\"\n}",
		},
		{
			name:                "feed",
			method:              "GET",
			path:                "/podcast/feed.rss",
			expectedStatus:      200,
			expectedContentType: ContentTypeTextXML,
			expectedBody:        testFeed,
		},
		{
			name:                "genre feed",
			method:              "GET",
			path:                "/podcast/genre/scifi/feed.rss",
			expectedStatus:      200,
			expectedContentType: ContentTypeTextXML,
			expectedBody:        fmt.Sprintf(feedFormat, audiobooks.SciFi.String()),
		},
		{
			name:                "author feed",
			method:              "GET",
			path:                "/podcast/authors/Agatha%20Test-y/feed.rss",
			expectedStatus:      200,
			expectedContentType: ContentTypeTextXML,
			expectedBody:        fmt.Sprintf(feedFormat, testAuthor),
		},
		{
			name:                "no author feed",
			method:              "GET",
			path:                "/podcast/authors/something/feed.rss",
			expectedStatus:      404,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Not Found",
		},
		{
			name:                "narrator feed",
			method:              "GET",
			path:                "/podcast/narrators/David%20Crochet/feed.rss",
			expectedStatus:      200,
			expectedContentType: ContentTypeTextXML,
			expectedBody:        fmt.Sprintf(feedFormat, testNarrator),
		},
		{
			name:                "no narrator feed",
			method:              "GET",
			path:                "/podcast/narrators/something/feed.rss",
			expectedStatus:      404,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Not Found",
		},
		{
			name:                "tag feed",
			method:              "GET",
			path:                "/podcast/tags/Cosy%20Classics/feed.rss",
			expectedStatus:      200,
			expectedContentType: ContentTypeTextXML,
			expectedBody:        fmt.Sprintf(feedFormat, testTag),
		},
		{
			name:                "no tag feed",
			method:              "GET",
			path:                "/podcast/tags/something/feed.rss",
			expectedStatus:      404,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Not Found",
		},
		{
			name:                "media",
			method:              "GET",
			path:                "/media/media.txt",
			expectedStatus:      200,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "served file",
		},
		{
			name:                "update",
			method:              "POST",
			path:                "/update",
			expectedStatus:      204,
			expectedContentType: "",
			expectedBody:        "",
		},
		{
			name:                "update on get fails",
			method:              "GET",
			path:                "/update",
			expectedStatus:      405,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Method Not Allowed",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NoError(t,
				baloo.New(testServer.URL).
					Method(testCase.method).
					Path(testCase.path).
					Request().
					Expect(t).
					Status(testCase.expectedStatus).
					Type(testCase.expectedContentType).
					BodyEquals(testCase.expectedBody).
					Done(),
			)
		})
	}
}

func TestHandler_Static(t *testing.T) {
	testHandler := NewHandler(&DummyPodcastService{}, "testdata")

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	tests := []struct {
		name                string
		path                string
		expectedContentType string
		expectedBodyLength  int
	}{
		{
			name:                "itunes image",
			path:                "/static/itunes_image.jpg",
			expectedContentType: "image/jpeg",
			expectedBodyLength:  261235,
		},
		{
			name:                "itunes image small",
			path:                "/static/itunes_image_small.jpg",
			expectedContentType: "image/jpeg",
			expectedBodyLength:  73150,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NoError(t,
				baloo.New(testServer.URL).
					Get(testCase.path).
					Expect(t).
					Status(200).
					Type(testCase.expectedContentType).
					BodyLength(testCase.expectedBodyLength).
					Done(),
			)
		})
	}
}

type DummyPodcastService struct{}

func (s *DummyPodcastService) WriteAllAudiobooksFeed(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(testFeed))
	return err
}

func (s *DummyPodcastService) WriteGenreAudiobookFeed(ctx context.Context, genre audiobooks.Genre, w io.Writer) error {
	_, err := fmt.Fprintf(w, feedFormat, genre.String())
	return err
}

func (s *DummyPodcastService) WriteAuthorAudiobookFeed(ctx context.Context, name string, w io.Writer) (bool, error) {
	var err error = nil
	if name == testAuthor {
		_, err = fmt.Fprintf(w, feedFormat, name)
		return true, err
	}
	return false, err
}

func (s *DummyPodcastService) WriteNarratorAudiobookFeed(ctx context.Context, name string, w io.Writer) (bool, error) {
	var err error = nil
	if name == testNarrator {
		_, err = fmt.Fprintf(w, feedFormat, name)
		return true, err
	}
	return false, err
}

func (s *DummyPodcastService) WriteTagAudiobookFeed(ctx context.Context, tag string, w io.Writer) (bool, error) {
	var err error = nil
	if tag == testTag {
		_, err = fmt.Fprintf(w, feedFormat, tag)
		return true, err
	}
	return false, err
}

func (s *DummyPodcastService) IsReady(ctx context.Context) bool {
	return true
}

func (s *DummyPodcastService) UpdateFeeds(ctx context.Context) error {
	return nil
}
