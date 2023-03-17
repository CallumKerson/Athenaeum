package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
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
			name:                "feed",
			path:                "/podcast/feed.rss",
			method:              "GET",
			expectedStatus:      200,
			expectedContentType: "application/xml; charset=utf-8",
			expectedBody:        getExpectedFeed(t, host)},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			baloo.New(host).
				Method(testCase.method).
				Path(testCase.path).
				Request().
				Expect(t).
				Status(testCase.expectedStatus).
				Type(testCase.expectedBody).
				BodyEquals(testCase.expectedBody)
		})
	}
}

func startRunCommand(t *testing.T) string {
	port := getFreePort(t)
	host := fmt.Sprintf("http://localhost:%d", port)
	tempDir := t.TempDir()
	envVarCleanup := envVarSetter(t, map[string]string{
		"ATHENAEUM_HOST":       host,
		"ATHENAEUM_PORT":       fmt.Sprintf("%d", port),
		"ATHENAEUM_DB_ROOT":    tempDir,
		"ATHENAEUM_MEDIA_ROOT": dataloader.GetRootTestdata(t),
	})
	t.Cleanup(envVarCleanup)

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

func getExpectedFeed(t *testing.T, host interface{}) string {
	var b bytes.Buffer
	tpl, err := template.ParseFiles(filepath.Join("testdata", "expected.rss"))
	assert.NoError(t, err)
	err = tpl.Execute(&b, host)
	assert.NoError(t, err)
	return strings.TrimSpace(b.String())
}
