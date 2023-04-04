package client

type Option func(c *Client)

func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}
