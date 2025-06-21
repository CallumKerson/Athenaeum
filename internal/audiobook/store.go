package audiobook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

const (
	defaultDBFileName       = "audiobooks.db"
	defaultDBFilePermission = 0o600
	defaultDBBucketName     = "audiobooks"
)

type Store struct {
	log                 loggerrific.Logger
	databaseRoot        string
	dbFileName          string
	dbFilePermission    fs.FileMode
	dbDefaultBucketName []byte
}

func NewStore(pathToDatabaseDir string, logger loggerrific.Logger) (*Store, error) {
	store := &Store{
		databaseRoot:        pathToDatabaseDir,
		log:                 logger,
		dbFileName:          defaultDBFileName,
		dbFilePermission:    defaultDBFilePermission,
		dbDefaultBucketName: []byte(defaultDBBucketName),
	}
	err := store.initialise()
	return store, err
}

func (s *Store) getDBPath() string {
	return filepath.Join(s.databaseRoot, s.dbFileName)
}

func (s *Store) initialise() error {
	err := os.MkdirAll(s.databaseRoot, 0755)
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

func (s *Store) StoreAll(ctx context.Context, allAudiobooks []audiobooks.Audiobook) error {
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

func (s *Store) GetAll(ctx context.Context) ([]audiobooks.Audiobook, error) {
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

func (s *Store) Get(ctx context.Context, filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
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

func (s *Store) IsReady(ctx context.Context) bool {
	return true
}
