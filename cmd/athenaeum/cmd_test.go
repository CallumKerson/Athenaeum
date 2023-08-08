package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/carlmjohnson/requests"
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/httpretry"
	"gopkg.in/h2non/baloo.v3"

	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
)

func TestRootCommand(t *testing.T) {
	host := startRunCommand(t)

	tests := []struct {
		name                string
		path                string
		method              string
		expectedStatus      int
		expectedContentType string
		expectedBody        string
	}{
		{
			name:                "index",
			path:                "/",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/html; charset=utf-8",
			expectedBody:        getExpected(t, "index.html", nil),
		},
		{
			name:                "version",
			path:                "/version",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        fmt.Sprintf("{\n  \"version\": %q\n}", Version),
		},
		{
			name:                "feed",
			path:                "/podcast/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/xml; charset=utf-8",
			expectedBody:        getExpected(t, "expected.rss", host),
		},
		{
			name:                "sci-fi feed",
			path:                "/podcast/genre/lgbt+/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/xml; charset=utf-8",
			expectedBody:        getExpected(t, "lgbt.rss", host),
		},
		{
			name:                "author feed",
			path:                "/podcast/authors/Ursula%20K.%20Le%20Guin/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/xml; charset=utf-8",
			expectedBody:        getExpected(t, "le_guin.rss", host),
		},
		{
			name:                "narrator feed",
			path:                "/podcast/narrators/Emily%20Woo%20Zeller/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/xml; charset=utf-8",
			expectedBody:        getExpected(t, "woo_zeller.rss", host),
		},
		{
			name:                "tag feed",
			path:                "/podcast/tags/Hugo%20Awards/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "text/xml; charset=utf-8",
			expectedBody:        getExpected(t, "hugo_awards.rss", host),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NoError(t,
				baloo.New(host).
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

func startRunCommand(t *testing.T) string {
	port := getFreePort(t)
	host := fmt.Sprintf("http://localhost:%d", port)
	tempDir := t.TempDir()
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_CONFIG_PATH": configFilePath,
	})
	t.Cleanup(envVarCleanup)

	configYAML := `---
Host: "%s"
Port: %d
DB:
    Root: "%s"
Media:
    Root: "%s"
`
	err := os.WriteFile(configFilePath, []byte(fmt.Sprintf(configYAML, host, port, tempDir, dataloader.GetRootTestdata(t))), 0644)
	assert.NoError(t, err)

	go func() {
		cmd := NewRootCommand()
		var b bytes.Buffer
		cmd.SetOut(&b)
		_ = cmd.Execute()
	}()

	assert.NoError(t,
		requests.
			URL(host).
			Client(httpretry.NewDefaultClient()).
			Path("health").
			Fetch(context.Background()),
	)
	return host
}

func getFreePort(t *testing.T) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func getExpected(t *testing.T, filename string, host interface{}) string {
	var b bytes.Buffer
	tpl, err := template.ParseFiles(filepath.Join("testdata", filename))
	assert.NoError(t, err)
	err = tpl.Execute(&b, host)
	assert.NoError(t, err)
	return strings.TrimSpace(b.String())
}
