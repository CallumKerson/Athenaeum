package bolt

import (
	"io/fs"

	"github.com/CallumKerson/loggerrific"
)

type Option func(s *AudiobookStore)

const (
	defaultDBFileName       = "athenaeum.db"
	defaultDBFilePermission = 0644
)

var (
	defaultDBBucketName = "athenaeum"
)

func WithDBFile(filename string, filePermission int) Option {
	return func(s *AudiobookStore) {
		s.dbFileName = filename
		s.dbFilePermission = fs.FileMode(filePermission)
	}
}

func WithDBBucketName(name string) Option {
	return func(s *AudiobookStore) {
		s.dbDefaultBucketName = []byte(name)
	}
}

func WithLogger(logger loggerrific.Logger) Option {
	return func(s *AudiobookStore) {
		s.log = logger
	}
}
