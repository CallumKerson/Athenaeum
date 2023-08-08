package service

import (
	"context"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type MediaScanner interface {
	GetAllAudiobooks(context.Context) ([]audiobooks.Audiobook, error)
	ScanForNewAndUpdatedAudiobooks(context.Context, []audiobooks.Audiobook) ([]audiobooks.Audiobook, bool, error)
}

type AudiobookStore interface {
	// StoreAll should remove any audiobooks not in the current slice and store those in the slice
	StoreAll(context.Context, []audiobooks.Audiobook) error
	GetAll(context.Context) ([]audiobooks.Audiobook, error)
	Get(context.Context, func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error)
	IsReady(context.Context) bool
}

type ThirdPartyNotifier interface {
	Notify(context.Context) error
	String() string
}

type Service struct {
	mediaScanner        MediaScanner
	audiobookStore      AudiobookStore
	thirdPartyNotifiers []ThirdPartyNotifier
	logger              loggerrific.Logger
	filtersForAll       []Filter
}

func New(mediaScanner MediaScanner, audiobookStore AudiobookStore, opts ...Option) *Service {
	svc := &Service{
		mediaScanner:   mediaScanner,
		audiobookStore: audiobookStore,
		logger:         noOpLogger.New(),
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *Service) UpdateAudiobooks(ctx context.Context) error {
	s.logger.Infoln("Updating audiobooks")
	existingAudiobooks, err := s.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	audiobooksFromScan, changed, err := s.mediaScanner.ScanForNewAndUpdatedAudiobooks(ctx, existingAudiobooks)
	if err != nil {
		return err
	}
	if changed {
		err = s.audiobookStore.StoreAll(ctx, audiobooksFromScan)
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
		return s.audiobookStore.Get(ctx, AndFilter(s.filtersForAll...))
	} else {
		return s.audiobookStore.GetAll(ctx)
	}
}

func (s *Service) IsReady(ctx context.Context) bool {
	return s.audiobookStore.IsReady(ctx)
}

func (s *Service) GetAudiobooksByAuthor(ctx context.Context, name string) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, AuthorFilter(name))
}

func (s *Service) GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, GenreFilter(genre))
}

func (s *Service) GetAudiobooksByNarrator(ctx context.Context, name string) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, NarratorFilter(name))
}

func (s *Service) GetAudiobooksByTag(ctx context.Context, tag string) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, TagFilter(tag))
}

func (s *Service) GetAudiobooksBy(ctx context.Context, filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, filter)
}
