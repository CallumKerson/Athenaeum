package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var (
	updatedThirdParty = false
)

func TestUpdate_Changes(t *testing.T) {
	tests := []struct {
		name                     string
		scanner                  MediaScanner
		expectedToCallThirdParty bool
	}{
		{name: "changes detected", scanner: &NewBooksMediaScanner{}, expectedToCallThirdParty: true},
		{name: "no changes detected", scanner: &NoChangesMediaScanner{}, expectedToCallThirdParty: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			updatedThirdParty = false

			testSvc := New(testCase.scanner, &DummyAudiobookStore{}, WithThirdPartyNotifier(&DummyNotifier{}))
			err := testSvc.UpdateAudiobooks(context.TODO())

			assert.False(t, updatedThirdParty, "updated not yet called")
			time.Sleep(300 * time.Millisecond)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedToCallThirdParty, updatedThirdParty)
		})
	}
}

func TestFilteringAllAudiobooks(t *testing.T) {
	tests := []struct {
		name               string
		filterOptions      []Option
		expectedAudiobooks []audiobooks.Audiobook
	}{
		{name: "no filter options", filterOptions: []Option{}, expectedAudiobooks: testbooks.Audiobooks},
		{
			name:               "filter genres from all audibooks",
			filterOptions:      []Option{WithGenresToExludeFromAllAudiobooks(audiobooks.SciFi)},
			expectedAudiobooks: []audiobooks.Audiobook{testbooks.Audiobooks[1]},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			updatedThirdParty = false

			testSvc := New(&NoChangesMediaScanner{}, &DummyAudiobookStore{}, testCase.filterOptions...)
			retrievedAudiobooks, err := testSvc.GetAllAudiobooks(context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedAudiobooks, retrievedAudiobooks)
		})
	}
}

type DummyNotifier struct {
}

func (u *DummyNotifier) Notify(context.Context) error {
	time.Sleep(100 * time.Millisecond)
	updatedThirdParty = true
	return nil
}

func (u *DummyNotifier) String() string {
	return "dummy"
}

type NewBooksMediaScanner struct {
}

func (s *NewBooksMediaScanner) GetAllAudiobooks(context.Context) ([]audiobooks.Audiobook, error) {
	return testbooks.Audiobooks, nil
}

func (s *NewBooksMediaScanner) ScanForNewAndUpdatedAudiobooks(context.Context, []audiobooks.Audiobook) ([]audiobooks.Audiobook, bool, error) {
	return testbooks.Audiobooks, true, nil
}

type NoChangesMediaScanner struct {
}

func (s *NoChangesMediaScanner) GetAllAudiobooks(context.Context) ([]audiobooks.Audiobook, error) {
	return testbooks.Audiobooks, nil
}

func (s *NoChangesMediaScanner) ScanForNewAndUpdatedAudiobooks(context.Context, []audiobooks.Audiobook) ([]audiobooks.Audiobook, bool, error) {
	return testbooks.Audiobooks, false, nil
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
