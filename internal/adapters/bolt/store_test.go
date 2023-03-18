package bolt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

func TestAudiobookStore_CreatesDBFileOnInit(t *testing.T) {
	// given
	dbRoot := filepath.Join(t.TempDir(), "db")

	// when
	_, err := bolt.NewAudiobookStore(tlogger.NewTLogger(t), true, bolt.WithPathToDBDirectory(dbRoot), bolt.WithDBDefaults())

	// then
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	fInfo, err := os.Stat(dbFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(32768), fInfo.Size())
}

func TestAudiobookStore_CreatesNoFileWhenNotInit(t *testing.T) {
	// given
	dbRoot := t.TempDir()

	// when
	_, err := bolt.NewAudiobookStore(tlogger.NewTLogger(t), false, bolt.WithPathToDBDirectory(dbRoot), bolt.WithDBDefaults())

	// then
	assert.NoError(t, err)
	assert.NoFileExists(t, filepath.Join(dbRoot, "athenaeum.db"))
}

func TestAudiobookStore_StoreAudiobooks(t *testing.T) {
	// given
	dbRoot := t.TempDir()
	store, err := bolt.NewAudiobookStore(tlogger.NewTLogger(t), true, bolt.WithPathToDBDirectory(dbRoot), bolt.WithDBDefaults())
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	fInfo, err := os.Stat(dbFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(32768), fInfo.Size())

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
	store, err := bolt.NewAudiobookStore(tlogger.NewTLogger(t), true, bolt.WithPathToDBDirectory(dbRoot), bolt.WithDBDefaults())
	assert.NoError(t, err)
	dbFile := filepath.Join(dbRoot, "athenaeum.db")
	assert.FileExists(t, dbFile)
	fInfo, err := os.Stat(dbFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(32768), fInfo.Size())

	// when
	err = store.StoreAll(context.TODO(), testbooks.Audiobooks)

	// then
	assert.NoError(t, err)

	// when
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
