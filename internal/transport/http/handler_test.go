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

func TestViewIndex(t *testing.T) {
	testHandler := NewHandler(&DummyService{}, tlogger.NewTLogger(t))

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
		name             string
		r                *http.Request
		expectedResponse string
	}{
		{name: "health check", r: newReq("GET", testServer.URL+"/health", nil), expectedResponse: "{\n  \"health\": \"ok\"\n}\n"},
		{name: "readiness check", r: newReq("GET", testServer.URL+"/ready", nil), expectedResponse: "{\n  \"readiness\": \"ok\"\n}\n"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(testCase.r)
			assert.NoError(t, err)
			defer resp.Body.Close()
			b, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, string(b))
		})
	}
}

type DummyService struct {
}

func (s *DummyService) IsReady(ctx context.Context) (bool, error) {
	return true, nil
}
