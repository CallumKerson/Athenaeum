package bolt

import "io/fs"

type Option func(s *AudiobookStore)

const (
	defaultDBFileName       = "athenaeum.db"
	defaultDBFilePermission = 0644
)

var (
	defaultDBBucketName = "athenaeum"
)

func WithPathToDBDirectory(path string) Option {
	return func(s *AudiobookStore) {
		s.databaseRoot = path
	}
}

func WithDBDefaults() Option {
	return func(s *AudiobookStore) {
		WithDBFile(defaultDBFileName, defaultDBFilePermission)(s)
		WithDBBucketName(defaultDBBucketName)(s)
	}
}

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
