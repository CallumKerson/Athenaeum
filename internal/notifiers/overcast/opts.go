package overcast

import "github.com/CallumKerson/loggerrific"

type Option func(s *Notifier)

func WithLogger(logger loggerrific.Logger) Option {
	return func(s *Notifier) {
		s.logger = logger
	}
}

func WithHost(host string) Option {
	return func(s *Notifier) {
		s.host = host
	}
}
