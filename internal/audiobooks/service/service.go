package service

import (
	"context"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type MediaScanner interface {
	GetAllAudiobooks(context.Context) ([]audiobooks.Audiobook, error)
}

type AudiobookStore interface {
	StoreAll(context.Context, []audiobooks.Audiobook) error
	GetAll(context.Context) ([]audiobooks.Audiobook, error)
	Get(context.Context, func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error)
	IsReady(context.Context) bool
}

type ThirdPartyUpdateService interface {
	Update(context.Context) error
}

type Service struct {
	mediaScanner             MediaScanner
	audiobookStore           AudiobookStore
	thirdPartyUpdateServices []ThirdPartyUpdateService
	logger                   loggerrific.Logger
}

func New(mediaScanner MediaScanner, audiobookStore AudiobookStore, logger loggerrific.Logger,
	thirdPartyUpdateServices ...ThirdPartyUpdateService) *Service {
	return &Service{
		mediaScanner:             mediaScanner,
		audiobookStore:           audiobookStore,
		logger:                   logger,
		thirdPartyUpdateServices: thirdPartyUpdateServices,
	}
}

func (s *Service) UpdateAudiobooks(ctx context.Context) error {
	audiobooksFromScan, err := s.mediaScanner.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	err = s.audiobookStore.StoreAll(ctx, audiobooksFromScan)
	if err != nil {
		return err
	}
	for svcIndex := range s.thirdPartyUpdateServices {
		go func(svc ThirdPartyUpdateService) {
			if updateErr := svc.Update(ctx); updateErr != nil {
				s.logger.WithError(err).Warnln("Update")
			}
		}(s.thirdPartyUpdateServices[svcIndex])
	}
	return nil
}

func (s *Service) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.GetAll(ctx)
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

func (s *Service) GetAudiobooksBy(ctx context.Context, filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.Get(ctx, filter)
}
