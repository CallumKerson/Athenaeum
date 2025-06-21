package audiobook

import (
	"context"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type ThirdPartyNotifier interface {
	Notify(context.Context) error
	String() string
}

type Service struct {
	store               *Store
	mediaRoot           string
	thirdPartyNotifiers []ThirdPartyNotifier
	logger              loggerrific.Logger
	filtersForAll       []Filter
}

func NewService(store *Store, mediaRoot string, logger loggerrific.Logger) *Service {
	return &Service{
		store:     store,
		mediaRoot: mediaRoot,
		logger:    logger,
	}
}

func (s *Service) WithNotifiers(notifiers []ThirdPartyNotifier) *Service {
	s.thirdPartyNotifiers = notifiers
	return s
}

func (s *Service) WithFilters(filters []Filter) *Service {
	s.filtersForAll = filters
	return s
}

func (s *Service) UpdateAudiobooks(ctx context.Context) error {
	s.logger.Infoln("Updating audiobooks")
	existingAudiobooks, err := s.store.GetAll()
	if err != nil {
		return err
	}

	audiobooksFromScan, changed, err := ScanForUpdates(ctx, s.mediaRoot, existingAudiobooks, s.logger)
	if err != nil {
		return err
	}

	if changed {
		err = s.store.StoreAll(audiobooksFromScan)
		if err != nil {
			return err
		}
		for svcIndex := range s.thirdPartyNotifiers {
			go func(ctx context.Context, notifier ThirdPartyNotifier) {
				if updateErr := notifier.Notify(ctx); updateErr != nil {
					s.logger.WithError(updateErr).Warnln("Notifying", notifier, "failed")
				}
			}(context.TODO(), s.thirdPartyNotifiers[svcIndex])
		}
		s.logger.Infoln("Update complete")
	} else {
		s.logger.Infoln("No updates detected")
	}
	return nil
}

func (s *Service) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	if len(s.filtersForAll) > 0 {
		return s.store.Get(AndFilter(s.filtersForAll...))
	} else {
		return s.store.GetAll()
	}
}

func (s *Service) IsReady(ctx context.Context) bool {
	return s.store.IsReady()
}

func (s *Service) GetAudiobooksByAuthor(ctx context.Context, name string) ([]audiobooks.Audiobook, error) {
	return s.store.Get(AuthorFilter(name))
}

func (s *Service) GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error) {
	return s.store.Get(GenreFilter(genre))
}

func (s *Service) GetAudiobooksByNarrator(ctx context.Context, name string) ([]audiobooks.Audiobook, error) {
	return s.store.Get(NarratorFilter(name))
}

func (s *Service) GetAudiobooksByTag(ctx context.Context, tag string) ([]audiobooks.Audiobook, error) {
	return s.store.Get(TagFilter(tag))
}

func (s *Service) GetAudiobooksBy(ctx context.Context, filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
	return s.store.Get(filter)
}
