package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type Client struct {
	host    string
	version string
}

func New(host string, opts ...Option) *Client {
	client := &Client{host: host}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Update(ctx context.Context) error {
	userAgent := "AthenaeumClient"
	if c.version != "" {
		userAgent = fmt.Sprintf("%s/%s", userAgent, c.version)
	}
	return requests.
		URL(c.host).
		Method(http.MethodPost).
		Header("User-Agent", userAgent).
		Path("update").
		Fetch(ctx)
}
