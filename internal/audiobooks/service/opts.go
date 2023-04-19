package service

import "github.com/CallumKerson/loggerrific"

type Option func(s *Service)

func WithLogger(logger loggerrific.Logger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithThirdPartyNotifier(notifier ThirdPartyNotifier) Option {
	return func(s *Service) {
		s.thirdPartyNotifiers = append(s.thirdPartyNotifiers, notifier)
	}
}
