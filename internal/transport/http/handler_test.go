package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"
)

const (
	testFeed = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\"></rss>"
)

func TestHandler(t *testing.T) {
	testHandler := NewHandler(&DummyService{}, tlogger.NewTLogger(t), WithMediaConfig("testdata", "/media/"))

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	newReq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequestWithContext(context.TODO(), method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name                string
		r                   *http.Request
		expectedStatus      int
		expectedContentType string
		expectedBody        string
	}{
		{name: "health check", r: newReq("GET", testServer.URL+"/health", nil), expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"health\": \"ok\"\n}\n"},
		{name: "readiness check", r: newReq("GET", testServer.URL+"/ready", nil), expectedStatus: 200, expectedContentType: ContentTypeJSON, expectedBody: "{\n  \"readiness\": \"ok\"\n}\n"},
		{name: "feed", r: newReq("GET", testServer.URL+"/podcast/feed.rss", nil), expectedStatus: 200, expectedContentType: ContentTypeXML, expectedBody: testFeed},
		{name: "media", r: newReq("GET", testServer.URL+"/media/media.txt", nil), expectedStatus: 200, expectedContentType: "text/plain; charset=utf-8", expectedBody: "served file\n"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(testCase.r)
			assert.NoError(t, err)
			defer resp.Body.Close()
			b, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)
			assert.Equal(t, testCase.expectedContentType, resp.Header.Get(ContentTypeHeader))
			assert.Equal(t, testCase.expectedBody, string(b))
		})
	}
}

type DummyService struct {
}

func (s *DummyService) WriteAllAudiobooksFeed(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(testFeed))
	return err
}

func (s *DummyService) IsReady(ctx context.Context) (bool, error) {
	return true, nil
}
