package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"
	bolt "go.etcd.io/bbolt"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type AudiobookStore struct {
	log                 loggerrific.Logger
	databaseRoot        string
	dbFileName          string
	dbFilePermission    fs.FileMode
	dbDefaultBucketName []byte
}

func NewAudiobookStore(pathToDatabaseDir string, opts ...Option) (*AudiobookStore, error) {
	store := &AudiobookStore{
		databaseRoot:        pathToDatabaseDir,
		log:                 noOpLogger.New(),
		dbFileName:          defaultDBFileName,
		dbFilePermission:    defaultDBFilePermission,
		dbDefaultBucketName: []byte(defaultDBBucketName),
	}
	for _, opt := range opts {
		opt(store)
	}
	err := store.Initialise()
	return store, err
}

func (s *AudiobookStore) getDBPath() string {
	return filepath.Join(s.databaseRoot, s.dbFileName)
}

func (s *AudiobookStore) Initialise() error {
	err := os.MkdirAll(s.databaseRoot, 0o755)
	if err != nil {
		return err
	}
	boltDB, err := bolt.Open(s.getDBPath(), s.dbFilePermission, nil)
	if err != nil {
		return fmt.Errorf("could not open bolt DB at %s, %w", s.getDBPath(), err)
	}
	defer boltDB.Close()
	s.log.Infoln("Setting up database at", s.getDBPath())
	return boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.dbDefaultBucketName)
		if err != nil {
			return fmt.Errorf("could not create bucket %s: %w", string(s.dbDefaultBucketName), err)
		}
		return nil
	})
}

func (s *AudiobookStore) StoreAll(ctx context.Context, allAudiobooks []audiobooks.Audiobook) error {
	boltDB, err := bolt.Open(s.getDBPath(), s.dbFilePermission, nil)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return boltDB.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket(s.dbDefaultBucketName)
		bucket, err := tx.CreateBucketIfNotExists(s.dbDefaultBucketName)
		if err != nil {
			return err
		}
		for index := range allAudiobooks {
			encoded, err := json.Marshal(allAudiobooks[index])
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(allAudiobooks[index].Path), encoded)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *AudiobookStore) GetAll(context.Context) ([]audiobooks.Audiobook, error) {
	boldDB, err := bolt.Open(s.getDBPath(), s.dbFilePermission, nil)
	if err != nil {
		return nil, err
	}
	defer boldDB.Close()
	var allAudiobooks []audiobooks.Audiobook
	return allAudiobooks, boldDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.dbDefaultBucketName)
		return b.ForEach(func(k, v []byte) error {
			audiobook := audiobooks.Audiobook{}
			err = json.Unmarshal(v, &audiobook)
			if err != nil {
				return err
			}
			allAudiobooks = append(allAudiobooks, audiobook)
			return nil
		})
	})
}

func (s *AudiobookStore) Get(
	ctx context.Context,
	filter func(*audiobooks.Audiobook) bool,
) ([]audiobooks.Audiobook, error) {
	boldDB, err := bolt.Open(s.getDBPath(), s.dbFilePermission, nil)
	if err != nil {
		return nil, err
	}
	defer boldDB.Close()
	var allAudiobooks []audiobooks.Audiobook
	return allAudiobooks, boldDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.dbDefaultBucketName)
		return b.ForEach(func(k, v []byte) error {
			audiobook := audiobooks.Audiobook{}
			err = json.Unmarshal(v, &audiobook)
			if err != nil {
				return err
			}
			if filter(&audiobook) {
				allAudiobooks = append(allAudiobooks, audiobook)
			}
			return nil
		})
	})
}

func (s *AudiobookStore) IsReady(context.Context) bool {
	return true
}
