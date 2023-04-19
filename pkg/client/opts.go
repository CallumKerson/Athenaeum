package client

type Option func(c *Client)

func WithVersion(version string) Option {
	return func(c *Client) {
		c.version = version
	}
}
