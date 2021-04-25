package filehash

import (
	"crypto/md5"
	"io"
	"os"
)

func NewMD5FileHasher() *MD5FileHasher {
	return &MD5FileHasher{}
}

type MD5FileHasher struct {
}

func (m *MD5FileHasher) Hash(pathToFile string) ([]byte, error) {

	file, err := os.Open(pathToFile)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil

}
