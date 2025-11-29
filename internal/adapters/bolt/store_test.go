package bolt_test

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestAudiobookStore_CreatesDBFileOnInit(t *testing.T) {
	// given
	dbRoot := filepath.Join(t.TempDir(), "db")

	// when
	_, err := bolt.NewAudiobookStore(dbRoot)

	// then
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	_, err = os.Stat(dbFile)
	assert.NoError(t, err)
}

func TestAudiobookStore_StoreAudiobooks(t *testing.T) {
	// given
	dbRoot := t.TempDir()
	store, err := bolt.NewAudiobookStore(dbRoot)
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	_, err = os.Stat(dbFile)
	assert.NoError(t, err)

	// when
	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)

	// then
	assert.NoError(t, err)

	// when
	retrievedBooks, err := store.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, retrievedBooks)
}

func TestAudiobookStore_GetAudiobooks(t *testing.T) {
	// given
	dbRoot := t.TempDir()
	store, err := bolt.NewAudiobookStore(dbRoot)
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	_, err = os.Stat(dbFile)
	assert.NoError(t, err)

	// when
	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)

	// then
	assert.NoError(t, err)

	// when
	retrievedBooks, err := store.Get(context.TODO(), func(a *audiobooks.Audiobook) bool {
		return slices.Contains(a.Authors, "Amal El-Mohtar")
	})
	assert.NoError(t, err)
	assert.Len(t, retrievedBooks, 1)
	assert.Equal(t, testbooks.Audiobooks[0], retrievedBooks[0])
}

func TestAudiobookStore_GetAllAudiobooks_WhenEmpty(t *testing.T) {
	// given
	dbRoot := t.TempDir()
	store, err := bolt.NewAudiobookStore(dbRoot)
	assert.NoError(t, err)

	// when
	retrievedBooks, err := store.GetAll(context.TODO())

	// then
	assert.NoError(t, err)
	assert.Empty(t, retrievedBooks)
}
