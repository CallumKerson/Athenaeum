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
	IsReady(context.Context) bool
}

type Service struct {
	mediaScanner   MediaScanner
	audiobookStore AudiobookStore
	logger         loggerrific.Logger
}

func New(mediaScanner MediaScanner, audiobookStore AudiobookStore, logger loggerrific.Logger) *Service {
	return &Service{mediaScanner: mediaScanner, audiobookStore: audiobookStore, logger: logger}
}

func (s *Service) ScanAndUpdateAudiobooks(ctx context.Context) error {
	audiobooksFromScan, err := s.mediaScanner.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}
	return s.audiobookStore.StoreAll(ctx, audiobooksFromScan)
}

func (s *Service) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	return s.audiobookStore.GetAll(ctx)
}

func (s *Service) IsReady(ctx context.Context) bool {
	return s.audiobookStore.IsReady(ctx)
}
