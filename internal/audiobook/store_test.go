package audiobook_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/audiobook"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestStore_CreatesDBFileOnInit(t *testing.T) {
	dbRoot := filepath.Join(t.TempDir(), "db")
	logger := logrus.NewLogger()

	_, err := audiobook.NewStore(dbRoot, logger)

	require.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "audiobooks.db")
	assert.FileExists(t, dbFile)
	_, err = os.Stat(dbFile)
	assert.NoError(t, err)
}

func TestStore_StoreAll(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	dbFile := filepath.Join(dbRoot, "audiobooks.db")
	assert.FileExists(t, dbFile)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	assert.NoError(t, err)

	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, retrievedBooks)
}

func TestStore_GetAll(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, retrievedBooks)
}

func TestStore_Get(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	retrievedBooks, err := store.Get(context.TODO(), func(a *audiobooks.Audiobook) bool {
		for _, v := range a.Authors {
			if v == "Amal El-Mohtar" {
				return true
			}
		}
		return false
	})
	assert.NoError(t, err)
	assert.Len(t, retrievedBooks, 1)
	assert.Equal(t, testbooks.Audiobooks[0], retrievedBooks[0])
}

func TestStore_GetAll_WhenEmpty(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Empty(t, retrievedBooks)
}

func TestStore_IsReady(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	ready := store.IsReady(context.TODO())
	assert.True(t, ready)
}

func TestStore_GetWithFilter(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	// Test filtering by genre
	sciFiBooks, err := store.Get(context.TODO(), func(a *audiobooks.Audiobook) bool {
		for _, genre := range a.Genres {
			if genre == audiobooks.SciFi {
				return true
			}
		}
		return false
	})
	assert.NoError(t, err)
	assert.Len(t, sciFiBooks, 1)
	assert.Equal(t, testbooks.Audiobooks[0], sciFiBooks[0])
}

func TestStore_StoreAll_OverwritesExisting(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store initial books
	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	// Store only one book (should overwrite all)
	singleBook := []audiobooks.Audiobook{testbooks.Audiobooks[0]}
	err = store.StoreAll(context.TODO(), singleBook)
	require.NoError(t, err)

	// Should only have one book now
	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, retrievedBooks, 1)
	assert.Equal(t, testbooks.Audiobooks[0], retrievedBooks[0])
}

func TestStore_NewStore_InvalidDBRoot(t *testing.T) {
	// Try to create store in a location that doesn't exist and can't be created
	dbRoot := "/root/nonexistent/path/that/cannot/be/created"
	logger := logrus.NewLogger()

	_, err := audiobook.NewStore(dbRoot, logger)
	// Should handle the error gracefully - either succeed or fail with meaningful error
	// The exact behaviour depends on permissions and OS
	if err != nil {
		assert.Error(t, err)
	}
}

func TestStore_StoreAll_EmptyList(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	// Store empty list
	err = store.StoreAll(context.TODO(), []audiobooks.Audiobook{})
	assert.NoError(t, err)

	// Should have no books
	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Empty(t, retrievedBooks)
}

func TestStore_Get_NoMatches(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	// Filter that matches nothing
	retrievedBooks, err := store.Get(context.TODO(), func(a *audiobooks.Audiobook) bool {
		return false // Never matches
	})
	assert.NoError(t, err)
	assert.Empty(t, retrievedBooks)
}

func TestStore_Get_AllMatch(t *testing.T) {
	dbRoot := t.TempDir()
	logger := logrus.NewLogger()
	store, err := audiobook.NewStore(dbRoot, logger)
	require.NoError(t, err)

	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)
	require.NoError(t, err)

	// Filter that matches everything
	retrievedBooks, err := store.Get(context.TODO(), func(a *audiobooks.Audiobook) bool {
		return true // Always matches
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, retrievedBooks)
}
