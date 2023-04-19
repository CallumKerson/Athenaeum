package overcast

import "github.com/CallumKerson/loggerrific"

type Option func(n *Notifier)

func WithLogger(logger loggerrific.Logger) Option {
	return func(n *Notifier) {
		n.logger = logger
	}
}

func WithHost(host string) Option {
	return func(n *Notifier) {
		n.host = host
	}
}
