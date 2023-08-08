package service

import (
	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

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

func WithGenresToExcludeFromAllAudiobooks(genres ...audiobooks.Genre) Option {
	return func(service *Service) {
		if len(genres) > 0 {
			var genreFilters []Filter
			for _, genre := range genres {
				genreFilters = append(genreFilters, NotFilter(GenreFilter(genre)))
			}
			service.filtersForAll = append(service.filtersForAll, genreFilters...)
		}
	}
}
