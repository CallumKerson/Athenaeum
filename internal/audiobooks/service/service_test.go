package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var (
	updatedThirdParty = false
)

func TestUpdate(t *testing.T) {
	testSvc := New(&DummyMediaScanner{}, &DummyAudiobookStore{}, tlogger.NewTLogger(t), &DummyUpdater{})
	err := testSvc.UpdateAudiobooks(context.TODO())

	// checks updater gets called in background
	assert.False(t, updatedThirdParty, "updated not yet called")
	time.Sleep(300 * time.Millisecond)

	assert.NoError(t, err)
	assert.True(t, updatedThirdParty, "called updater")
}

type DummyUpdater struct {
}

func (u *DummyUpdater) Update(context.Context) error {
	time.Sleep(100 * time.Millisecond)
	updatedThirdParty = true
	return nil
}

type DummyMediaScanner struct {
}

func (s *DummyMediaScanner) GetAllAudiobooks(context.Context) ([]audiobooks.Audiobook, error) {
	return testbooks.Audiobooks, nil
}

type DummyAudiobookStore struct {
}

func (s *DummyAudiobookStore) StoreAll(context.Context, []audiobooks.Audiobook) error {
	return nil
}

func (s *DummyAudiobookStore) GetAll(context.Context) ([]audiobooks.Audiobook, error) {
	return testbooks.Audiobooks, nil
}

func (s *DummyAudiobookStore) Get(ctx context.Context, filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
	var filtered []audiobooks.Audiobook
	all, _ := s.GetAll(ctx)
	for index := range all {
		if filter(&all[index]) {
			filtered = append(filtered, all[index])
		}
	}
	return filtered, nil
}

func (s *DummyAudiobookStore) IsReady(context.Context) bool {
	return true
}
