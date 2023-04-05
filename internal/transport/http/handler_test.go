package http

import (
	"context"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	testFeed  = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\"></rss>"
	sciFiFeed = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\">" +
		"<channel><title>Science Fiction</title></channel></rss>"
)

func TestHandler(t *testing.T) {
	testHandler := NewHandler(&DummyPodcastService{}, &DummyUpdateService{}, tlogger.NewTLogger(t),
		WithMediaConfig(filepath.Join("testdata", "media"), "/media/"), WithVersion("1.0.0-test"), WithStaticPath("/static"))

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
		{name: "index", method: "GET", path: "", expectedStatus: 200, expectedContentType: ContentTypeHTML, expectedBody: strings.TrimSpace(string(expectedIndex))},
		{name: "health check", method: "GET", path: "/health", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"health\": \"ok\"\n}"},
		{name: "readiness check", method: "GET", path: "/ready", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"readiness\": \"ok\"\n}"},
		{name: "version", method: "GET", path: "/version", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"version\": \"1.0.0-test\"\n}"},
		{name: "feed", method: "GET", path: "/podcast/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeXML, expectedBody: testFeed},
		{name: "feed", method: "GET", path: "/podcast/genre/scifi/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeXML, expectedBody: sciFiFeed},
		{name: "media", method: "GET", path: "/media/media.txt", expectedStatus: 200, expectedContentType: "text/plain; charset=utf-8", expectedBody: "served file"},
		{name: "update", method: "POST", path: "/update", expectedStatus: 204, expectedContentType: "", expectedBody: ""},
		{name: "update on get fails", method: "GET", path: "/update", expectedStatus: 405, expectedContentType: "text/plain; charset=utf-8", expectedBody: "Method Not Allowed"},
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
	testHandler := NewHandler(&DummyPodcastService{}, &DummyUpdateService{}, tlogger.NewTLogger(t),
		WithMediaConfig("testdata", "/media/"), WithVersion("1.0.0-test"), WithStaticPath("/static"))

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	tests := []struct {
		name                string
		path                string
		expectedContentType string
		expectedBodyLength  int
	}{
		{name: "itunes image", path: "/static/itunes_image.jpg", expectedContentType: "image/jpeg", expectedBodyLength: 261235},
		{name: "itunes image small", path: "/static/itunes_image_small.jpg", expectedContentType: "image/jpeg", expectedBodyLength: 73150},
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

type DummyPodcastService struct {
}

func (s *DummyPodcastService) WriteAllAudiobooksFeed(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(testFeed))
	return err
}

func (s *DummyPodcastService) WriteGenreAudiobookFeed(ctx context.Context, genre audiobooks.Genre, w io.Writer) error {
	if genre == audiobooks.SciFi {
		_, err := w.Write([]byte(sciFiFeed))
		return err
	} else {
		_, err := w.Write([]byte(testFeed))
		return err
	}
}

func (s *DummyPodcastService) IsReady(ctx context.Context) bool {
	return true
}

type DummyUpdateService struct {
}

func (s *DummyUpdateService) UpdateAudiobooks(ctx context.Context) error {
	return nil
}

func (s *DummyUpdateService) IsReady(ctx context.Context) bool {
	return true
}
