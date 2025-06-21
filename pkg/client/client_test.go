package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBaseURL = "http://example.com"
	testVersion = "v1.2.3"
)

func TestNew(t *testing.T) {
	client := New(testBaseURL)

	assert.NotNil(t, client)
}

func TestNew_WithOptions(t *testing.T) {
	client := New(testBaseURL, WithVersion(testVersion))

	assert.NotNil(t, client)
}

func TestClient_Update_Success(t *testing.T) {
	// Create a test server that responds with 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/update", r.URL.Path)

		// Check User-Agent header if version is set
		userAgent := r.Header.Get("User-Agent")
		if userAgent != "" {
			assert.Contains(t, userAgent, "AthenaeumClient")
		}

		responseWriter.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)

	err := client.Update(context.Background())
	assert.NoError(t, err)
}

func TestClient_Update_WithVersion(t *testing.T) {
	expectedUserAgent := "AthenaeumClient/" + testVersion

	// Create a test server that checks the User-Agent header
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/update", r.URL.Path)
		assert.Equal(t, expectedUserAgent, r.Header.Get("User-Agent"))

		responseWriter.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, WithVersion(testVersion))

	err := client.Update(context.Background())
	assert.NoError(t, err)
}

func TestClient_Update_ServerError(t *testing.T) {
	// Create a test server that responds with 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := New(server.URL)

	err := client.Update(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestClient_Update_ClientError(t *testing.T) {
	// Create a test server that responds with 404 Not Found
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := New(server.URL)

	err := client.Update(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestClient_Update_InvalidURL(t *testing.T) {
	// Use an invalid URL that will cause a connection error
	client := New("http://invalid-host-that-does-not-exist.local:12345")

	err := client.Update(context.Background())
	assert.Error(t, err)
}

func TestClient_Update_ContextCancellation(t *testing.T) {
	// Create a test server that delays the response
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		// This handler will never respond, simulating a slow server
		select {}
	}))
	defer server.Close()

	client := New(server.URL)

	// Create a context that is immediately cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Update(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestClient_Update_EmptyBaseURL(t *testing.T) {
	client := New("")

	err := client.Update(context.Background())
	assert.Error(t, err)
}

func TestClient_Update_MalformedBaseURL(t *testing.T) {
	client := New("not-a-valid-url")

	err := client.Update(context.Background())
	assert.Error(t, err)
}

func TestWithVersion(t *testing.T) {
	version := "test-version"

	client := New(testBaseURL, WithVersion(version))

	assert.NotNil(t, client)
}

func TestClient_UpdatePath(t *testing.T) {
	// Create a test server that captures the request path
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		responseWriter.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)

	err := client.Update(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "/update", capturedPath)
}
