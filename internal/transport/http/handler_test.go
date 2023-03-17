package http

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"

	"github.com/CallumKerson/loggerrific/tlogger"
)

const (
	testFeed = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\"></rss>"
)

func TestHandler(t *testing.T) {
	testHandler := NewHandler(&DummyPodcastService{}, &DummyUpdateService{}, tlogger.NewTLogger(t),
		WithMediaConfig("testdata", "/media/"), WithVersion("1.0.0-test"))

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	tests := []struct {
		name                string
		method              string
		path                string
		expectedStatus      int
		expectedContentType string
		expectedBody        string
	}{
		{name: "health check", method: "GET", path: "/health", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"health\": \"ok\"\n}"},
		{name: "readiness check", method: "GET", path: "/ready", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"readiness\": \"ok\"\n}"},
		{name: "version", method: "GET", path: "/version", expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"version\": \"1.0.0-test\"\n}"},
		{name: "feed", method: "GET", path: "/podcast/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeXML, expectedBody: testFeed},
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

type DummyPodcastService struct {
}

func (s *DummyPodcastService) WriteAllAudiobooksFeed(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(testFeed))
	return err
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
