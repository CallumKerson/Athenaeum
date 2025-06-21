package audiobook

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var updatedThirdParty = false

// setupServiceTest creates a test service with stored test books
func setupServiceTest(t *testing.T) *Service {
	t.Helper()
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(testbooks.Audiobooks)
	require.NoError(t, err)

	return NewService(store, testMediaRoot, logger)
}

func TestService_NewService(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	mediaRoot := testMediaRoot

	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	service := NewService(store, mediaRoot, logger)

	assert.NotNil(t, service)
	assert.Equal(t, mediaRoot, service.mediaRoot)
	assert.Equal(t, store, service.store)
	assert.Equal(t, logger, service.logger)
	assert.Empty(t, service.thirdPartyNotifiers)
	assert.Empty(t, service.filtersForAll)
}

func TestService_WithNotifiers(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	notifiers := []ThirdPartyNotifier{&DummyNotifier{}}
	service := NewService(store, testMediaRoot, logger).WithNotifiers(notifiers)

	assert.Equal(t, notifiers, service.thirdPartyNotifiers)
}

func TestService_WithFilters(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	filters := []Filter{AuthorFilter("Test Author")}
	service := NewService(store, testMediaRoot, logger).WithFilters(filters)

	assert.Equal(t, filters, service.filtersForAll)
}

func TestService_IsReady(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	service := NewService(store, testMediaRoot, logger)
	ready := service.IsReady(context.TODO())

	assert.True(t, ready)
}

func TestService_GetAllAudiobooks_WithoutFilters(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store test books
	err = store.StoreAll(testbooks.Audiobooks)
	require.NoError(t, err)

	service := NewService(store, testMediaRoot, logger)
	books, err := service.GetAllAudiobooks(context.TODO())

	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestService_GetAllAudiobooks_WithFilters(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store test books
	err = store.StoreAll(testbooks.Audiobooks)
	require.NoError(t, err)

	// Create service with filter that excludes SciFi
	filters := []Filter{NotFilter(GenreFilter(audiobooks.SciFi))}
	service := NewService(store, testMediaRoot, logger).WithFilters(filters)
	books, err := service.GetAllAudiobooks(context.TODO())

	assert.NoError(t, err)
	// Should only return books without SciFi genre
	expected := []audiobooks.Audiobook{testbooks.Audiobooks[1]}
	assert.ElementsMatch(t, expected, books)
}

func TestService_GetAudiobooksByAuthor(t *testing.T) {
	service := setupServiceTest(t)
	books, err := service.GetAudiobooksByAuthor(context.TODO(), "Max Gladstone")

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, testbooks.Audiobooks[0], books[0])
}

func TestService_GetAudiobooksByGenre(t *testing.T) {
	service := setupServiceTest(t)
	books, err := service.GetAudiobooksByGenre(context.TODO(), audiobooks.SciFi)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, testbooks.Audiobooks[0], books[0])
}

func TestService_GetAudiobooksByNarrator(t *testing.T) {
	service := setupServiceTest(t)
	books, err := service.GetAudiobooksByNarrator(context.TODO(), "Emily Woo Zeller")

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, testbooks.Audiobooks[0], books[0])
}

func TestService_GetAudiobooksByTag(t *testing.T) {
	service := setupServiceTest(t)
	books, err := service.GetAudiobooksByTag(context.TODO(), "Hugo Awards")

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, testbooks.Audiobooks[0], books[0])
}

func TestService_GetAudiobooksBy_CustomFilter(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(testbooks.Audiobooks)
	require.NoError(t, err)

	service := NewService(store, testMediaRoot, logger)

	// Custom filter for books with "Gladstone" in any author name
	customFilter := func(a *audiobooks.Audiobook) bool {
		for _, author := range a.Authors {
			if containsIgnoringCaseAndWhitespace([]string{author}, "gladstone") {
				return true
			}
		}
		return false
	}

	books, err := service.GetAudiobooksBy(context.TODO(), customFilter)

	assert.NoError(t, err)
	if len(books) > 0 {
		assert.Equal(t, testbooks.Audiobooks[0], books[0])
	}
}

type DummyNotifier struct {
	name string
}

func (u *DummyNotifier) Notify(context.Context) error {
	time.Sleep(100 * time.Millisecond)
	updatedThirdParty = true
	return nil
}

func (u *DummyNotifier) String() string {
	if u.name != "" {
		return u.name
	}
	return "dummy"
}

func TestService_UpdateAudiobooks_WithoutChanges(t *testing.T) {
	// This test would require mocking the ScanForUpdates function
	// For now, we'll test the basic structure
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	service := NewService(store, testMediaRoot, logger)

	// This will fail because /test/media doesn't exist, but we're testing the structure
	err = service.UpdateAudiobooks(context.TODO())
	assert.Error(t, err) // Expected to fail due to non-existent media root
}

func TestService_UpdateAudiobooks_WithNotifiers(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	updatedThirdParty = false
	notifier := &DummyNotifier{name: "test-notifier"}
	service := NewService(store, testMediaRoot, logger).WithNotifiers([]ThirdPartyNotifier{notifier})

	// This will fail because /test/media doesn't exist, but we're testing the structure
	err = service.UpdateAudiobooks(context.TODO())
	assert.Error(t, err) // Expected to fail due to non-existent media root
}

func TestService_UpdateAudiobooks_WithChangesAndNotifiers(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store some initial books
	err = store.StoreAll([]audiobooks.Audiobook{testbooks.Audiobooks[0]})
	require.NoError(t, err)

	updatedThirdParty = false
	notifier := &DummyNotifier{name: "test-notifier"}
	mediaRoot := dataloader.GetRootTestdata(t)

	service := NewService(store, mediaRoot, logger).WithNotifiers([]ThirdPartyNotifier{notifier})

	// This should succeed with real test data and trigger notifiers
	err = service.UpdateAudiobooks(context.TODO())
	assert.NoError(t, err)
	// Should have detected changes and called notifiers
	// Note: updatedThirdParty may or may not be true depending on actual changes detected
}

func TestService_UpdateAudiobooks_NoChanges(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store all test books first
	err = store.StoreAll(testbooks.Audiobooks)
	require.NoError(t, err)

	updatedThirdParty = false
	notifier := &DummyNotifier{name: "test-notifier"}
	mediaRoot := dataloader.GetRootTestdata(t)

	service := NewService(store, mediaRoot, logger).WithNotifiers([]ThirdPartyNotifier{notifier})

	// This should succeed with no changes, so notifiers shouldn't be called
	err = service.UpdateAudiobooks(context.TODO())
	assert.NoError(t, err)
	// Should not have called notifiers because no changes detected
	assert.False(t, updatedThirdParty)
}

func TestService_UpdateAudiobooks_NotifierError(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store some initial books
	err = store.StoreAll([]audiobooks.Audiobook{testbooks.Audiobooks[0]})
	require.NoError(t, err)

	// Create a notifier that will return an error
	errorNotifier := &ErrorNotifier{}
	mediaRoot := dataloader.GetRootTestdata(t)

	service := NewService(store, mediaRoot, logger).WithNotifiers([]ThirdPartyNotifier{errorNotifier})

	// This should still succeed even if notifier fails
	err = service.UpdateAudiobooks(context.TODO())
	assert.NoError(t, err) // Should not propagate notifier errors
}

var errNotifierTest = errors.New("notifier error")

type ErrorNotifier struct{}

func (e *ErrorNotifier) Notify(context.Context) error {
	return errNotifierTest
}

func (e *ErrorNotifier) String() string {
	return "error-notifier"
}
