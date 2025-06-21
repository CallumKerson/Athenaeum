package audiobook

import (
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
	dbPath              string
}

func NewStore(pathToDatabaseDir string, logger loggerrific.Logger) (*Store, error) {
	store := &Store{
		databaseRoot:        pathToDatabaseDir,
		log:                 logger,
		dbFileName:          defaultDBFileName,
		dbFilePermission:    defaultDBFilePermission,
		dbDefaultBucketName: []byte(defaultDBBucketName),
		dbPath:              filepath.Join(pathToDatabaseDir, defaultDBFileName),
	}
	err := initializeDB(store.dbPath, store.dbFilePermission, store.dbDefaultBucketName, logger)
	return store, err
}

func initializeDB(dbPath string, dbFilePermission fs.FileMode, bucketName []byte, logger loggerrific.Logger) error {
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		return err
	}
	boltDB, err := bolt.Open(dbPath, dbFilePermission, nil)
	if err != nil {
		return fmt.Errorf("could not open bolt DB at %s, %w", dbPath, err)
	}
	defer boltDB.Close()
	logger.Infoln("Setting up database at", dbPath)
	return boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("could not create bucket %s: %w", string(bucketName), err)
		}
		return nil
	})
}

func StoreAll(dbPath string, dbFilePermission fs.FileMode, bucketName []byte, allAudiobooks []audiobooks.Audiobook) error {
	boltDB, err := bolt.Open(dbPath, dbFilePermission, nil)
	if err != nil {
		return err
	}
	defer boltDB.Close()
	return boltDB.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket(bucketName)
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
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

func GetAll(dbPath string, dbFilePermission fs.FileMode, bucketName []byte) ([]audiobooks.Audiobook, error) {
	boltDB, err := bolt.Open(dbPath, dbFilePermission, nil)
	if err != nil {
		return nil, err
	}
	defer boltDB.Close()
	var allAudiobooks []audiobooks.Audiobook
	return allAudiobooks, boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
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

func Get(
	dbPath string,
	dbFilePermission fs.FileMode,
	bucketName []byte,
	filter func(*audiobooks.Audiobook) bool,
) ([]audiobooks.Audiobook, error) {
	boltDB, err := bolt.Open(dbPath, dbFilePermission, nil)
	if err != nil {
		return nil, err
	}
	defer boltDB.Close()
	var allAudiobooks []audiobooks.Audiobook
	return allAudiobooks, boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
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

func IsReady() bool {
	return true
}

// Store methods that call the package-level functions
func (s *Store) StoreAll(allAudiobooks []audiobooks.Audiobook) error {
	return StoreAll(s.dbPath, s.dbFilePermission, s.dbDefaultBucketName, allAudiobooks)
}

func (s *Store) GetAll() ([]audiobooks.Audiobook, error) {
	return GetAll(s.dbPath, s.dbFilePermission, s.dbDefaultBucketName)
}

func (s *Store) Get(filter func(*audiobooks.Audiobook) bool) ([]audiobooks.Audiobook, error) {
	return Get(s.dbPath, s.dbFilePermission, s.dbDefaultBucketName, filter)
}

func (s *Store) IsReady() bool {
	return IsReady()
}
