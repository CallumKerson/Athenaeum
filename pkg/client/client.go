package client

import (
	"context"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type Client struct {
	host string
}

func New(opts ...Option) *Client {
	client := &Client{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Update(ctx context.Context) error {
	return requests.
		URL(c.host).
		Method(http.MethodPost).
		Path("update").
		Fetch(ctx)
}
