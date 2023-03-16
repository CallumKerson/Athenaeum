package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
)

var (
	errServerNotReady = errors.New("server not ready")
)

func TestRootCommand(t *testing.T) {
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

	err := retry.Do(
		func() error {
			req := newRequest(t, "GET", fmt.Sprintf("%s/health", host))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return errServerNotReady
			}
			return nil
		},
	)
	assert.NoError(t, err)

	req := newRequest(t, "GET", fmt.Sprintf("%s/podcast/feed.rss", host))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/xml; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, getExpectedFeed(t, host), string(body))
}

func getFreePort(t *testing.T) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func newRequest(t *testing.T, method, url string) *http.Request {
	r, err := http.NewRequestWithContext(context.TODO(), method, url, http.NoBody)
	assert.NoError(t, err)
	return r
}

func getExpectedFeed(t *testing.T, host interface{}) string {
	var b bytes.Buffer
	tpl, err := template.ParseFiles(filepath.Join("testdata", "expected.rss"))
	assert.NoError(t, err)
	err = tpl.Execute(&b, host)
	assert.NoError(t, err)
	return strings.TrimSpace(b.String())
}
