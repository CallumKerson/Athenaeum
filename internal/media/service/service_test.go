package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/m4b"
)

func TestService_GetAudiobooks(t *testing.T) {
	svc := New(&DummyM4BService{}, dataloader.GetRootTestdata(t))

	books, err := svc.GetAllAudiobooks(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestService_UpdateAudiobooks(t *testing.T) {
	svc := New(&DummyM4BService{}, dataloader.GetRootTestdata(t))

	books, changed, err := svc.ScanForNewAndUpdatedAudiobooks(
		context.TODO(),
		[]audiobooks.Audiobook{testbooks.Audiobooks[0]},
	)
	assert.NoError(t, err)
	assert.True(t, changed)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

type DummyM4BService struct{}

func (s *DummyM4BService) Read(pathToM4BFile string) (*m4b.Metadata, error) {
	return &m4b.Metadata{Duration: time.Nanosecond * 4671000064}, nil
}
