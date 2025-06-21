package http

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"

	"github.com/CallumKerson/Athenaeum/internal/audiobook"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	testAuthor   = "Agatha Test-y"
	testNarrator = "David Crochet"
	testTag      = "Cosy Classics"
)

func TestHandler(t *testing.T) {
	audiobookService := &DummyAudiobookService{}
	feedConfig := &audiobook.FeedConfig{
		Link:      "http://example.com",
		ImageLink: "http://example.com/image.jpg",
		Host:      "http://example.com",
		MediaPath: "/media",
	}
	testHandler := NewHandler(audiobookService, feedConfig, filepath.Join("testdata", "media"), WithVersion("1.0.0-test"))

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
		{name: "feed", method: "GET", path: "/podcast/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeTextXML, expectedBody: ""},
		{name: "genre feed", method: "GET", path: "/podcast/genre/scifi/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeTextXML, expectedBody: ""},
		{name: "author feed", method: "GET", path: "/podcast/authors/Agatha%20Test-y/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeTextXML, expectedBody: ""},
		{name: "no author feed", method: "GET", path: "/podcast/authors/something/feed.rss", expectedStatus: 404, expectedContentType: "text/plain; charset=utf-8", expectedBody: "Not Found"},
		{name: "narrator feed", method: "GET", path: "/podcast/narrators/David%20Crochet/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeTextXML, expectedBody: ""},
		{name: "no narrator feed", method: "GET", path: "/podcast/narrators/something/feed.rss", expectedStatus: 404, expectedContentType: "text/plain; charset=utf-8", expectedBody: "Not Found"},
		{name: "tag feed", method: "GET", path: "/podcast/tags/Cosy%20Classics/feed.rss", expectedStatus: 200, expectedContentType: ContentTypeTextXML, expectedBody: ""},
		{name: "no tag feed", method: "GET", path: "/podcast/tags/something/feed.rss", expectedStatus: 404, expectedContentType: "text/plain; charset=utf-8", expectedBody: "Not Found"},
		{name: "media", method: "GET", path: "/media/media.txt", expectedStatus: 200, expectedContentType: "text/plain; charset=utf-8", expectedBody: "served file"},
		{name: "update", method: "POST", path: "/update", expectedStatus: 204, expectedContentType: "", expectedBody: ""},
		{name: "update on get fails", method: "GET", path: "/update", expectedStatus: 405, expectedContentType: "text/plain; charset=utf-8", expectedBody: "Method Not Allowed"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expectedBody == "" && testCase.expectedContentType == ContentTypeTextXML {
				// For XML feeds, just check status and content type
				assert.NoError(t,
					baloo.New(testServer.URL).
						Method(testCase.method).
						Path(testCase.path).
						Request().
						Expect(t).
						Status(testCase.expectedStatus).
						Type(testCase.expectedContentType).
						Done(),
				)
			} else {
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
			}
		})
	}
}

func TestHandler_Static(t *testing.T) {
	audiobookService := &DummyAudiobookService{}
	feedConfig := &audiobook.FeedConfig{
		Link:      "http://example.com",
		ImageLink: "http://example.com/image.jpg",
		Host:      "http://example.com",
		MediaPath: "/media",
	}
	testHandler := NewHandler(audiobookService, feedConfig, "testdata")

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

type DummyAudiobookService struct {
}

func (s *DummyAudiobookService) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	return []audiobooks.Audiobook{}, nil
}

func (s *DummyAudiobookService) GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error) {
	return []audiobooks.Audiobook{}, nil
}

func (s *DummyAudiobookService) GetAudiobooksByAuthor(ctx context.Context, author string) ([]audiobooks.Audiobook, error) {
	if author == testAuthor {
		return []audiobooks.Audiobook{{
			Title:       "Test Book",
			Authors:     []string{author},
			Path:        "/test.m4b",
			FileSize:    1000,
			MIMEType:    "audio/mp4a-latm",
			Duration:    time.Minute,
			ReleaseDate: &toml.LocalDate{Year: 2020, Month: 1, Day: 1},
		}}, nil
	}
	return []audiobooks.Audiobook{}, nil
}

func (s *DummyAudiobookService) GetAudiobooksByNarrator(ctx context.Context, narrator string) ([]audiobooks.Audiobook, error) {
	if narrator == testNarrator {
		return []audiobooks.Audiobook{{
			Title:       "Test Book",
			Narrators:   []string{narrator},
			Path:        "/test.m4b",
			FileSize:    1000,
			MIMEType:    "audio/mp4a-latm",
			Duration:    time.Minute,
			ReleaseDate: &toml.LocalDate{Year: 2020, Month: 1, Day: 1},
		}}, nil
	}
	return []audiobooks.Audiobook{}, nil
}

func (s *DummyAudiobookService) GetAudiobooksByTag(ctx context.Context, tag string) ([]audiobooks.Audiobook, error) {
	if tag == testTag {
		return []audiobooks.Audiobook{{
			Title:       "Test Book",
			Tags:        []string{tag},
			Path:        "/test.m4b",
			FileSize:    1000,
			MIMEType:    "audio/mp4a-latm",
			Duration:    time.Minute,
			ReleaseDate: &toml.LocalDate{Year: 2020, Month: 1, Day: 1},
		}}, nil
	}
	return []audiobooks.Audiobook{}, nil
}

func (s *DummyAudiobookService) UpdateAudiobooks(ctx context.Context) error {
	return nil
}

func (s *DummyAudiobookService) IsReady(ctx context.Context) bool {
	return true
}
