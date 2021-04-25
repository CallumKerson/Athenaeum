package books

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

type Scanner struct {
	logger             logging.Logger
	dirToScan          string
	factory            Factory
	retriever          Retriever
	fileHasher         FileHasher
	booksFoundLastScan int
	lastScanned        *time.Time
}

func NewScanner(logger logging.Logger, dirToScan string, factory Factory, retriever Retriever, fileHasher FileHasher) *Scanner {
	return &Scanner{logger, dirToScan, factory, retriever, fileHasher, 0, nil}
}

func (s *Scanner) Scan() (int, error) {
	numberOfFilesScanned := 0
	s.logger.Infof("Scanning %s for audiobook files", s.dirToScan)
	currentTime := time.Now()
	s.lastScanned = &currentTime
	relativeFoundFiles, err := getRelativeM4BFiles(s.dirToScan)
	if err != nil {
		s.logger.Errorf("Cannot find audiobook files from %s", s.dirToScan)
		return numberOfFilesScanned, err
	}

	for _, relativeFoundFile := range relativeFoundFiles {

		//Check if path is duplicate of existing path
		allFilepaths, err := s.retriever.GetAllLocations()
		if err != nil {
			s.booksFoundLastScan = numberOfFilesScanned
			s.logger.Errorf("Cannot get all locations for existing files")
			return numberOfFilesScanned, err
		}
		foundDuplicateByLocation := find(allFilepaths, relativeFoundFile)
		if foundDuplicateByLocation {
			s.booksFoundLastScan = numberOfFilesScanned
			s.logger.Debugf("File already added to library %s", relativeFoundFile)
		} else {
			//Check if found file hash already exists by hashing file
			allHashes, err := s.retriever.GetAllHashes()
			if err != nil {
				s.booksFoundLastScan = numberOfFilesScanned
				s.logger.Errorf("Cannot get all hashes for existing files")
				return numberOfFilesScanned, err
			}
			fileHash, err := s.fileHasher.Hash(filepath.Join(s.dirToScan, relativeFoundFile))
			if err != nil {
				s.booksFoundLastScan = numberOfFilesScanned
				s.logger.Errorf("Cannot get hash for file %s", filepath.Join(s.dirToScan, relativeFoundFile))
				return numberOfFilesScanned, err
			}
			foundDuplicateByHash := findHash(allHashes, fileHash)
			if foundDuplicateByHash {
				existingBook, _ := s.retriever.GetByFileHash(fileHash)
				s.logger.Warnf("Will not process scanned file %s, has identical file hash to book with ID %s", filepath.Join(s.dirToScan, relativeFoundFile), existingBook.Id)

			} else {
				numberOfFilesScanned += 1
				_, err = s.factory.NewBookFromAudioMetadata(relativeFoundFile)
				if err != nil {
					s.logger.Errorf("Cannot create new audiobook from file %s", filepath.Join(s.dirToScan, relativeFoundFile))
					s.booksFoundLastScan = numberOfFilesScanned
					return numberOfFilesScanned, err
				}
			}
		}
	}
	s.booksFoundLastScan = numberOfFilesScanned
	return numberOfFilesScanned, nil
}

func (s *Scanner) LastScanned() time.Time {
	return *s.lastScanned
}

func (s *Scanner) BooksFoundLastScan() int {
	return s.booksFoundLastScan
}

func getRelativeM4BFiles(audiobookRoot string) ([]string, error) {
	var files []string
	err := filepath.Walk(audiobookRoot, func(fullPath string, info os.FileInfo, err error) error {
		if filepath.Ext(fullPath) == ".m4b" {
			relativePath := strings.TrimPrefix(fullPath, audiobookRoot)
			if !strings.HasPrefix(relativePath, "/") {
				relativePath = "/" + relativePath
			}
			files = append(files, relativePath)
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, nil
}

func findHash(slice [][]byte, value []byte) bool {
	for _, item := range slice {
		if bytes.Equal(item, value) {
			return true
		}
	}
	return false
}

func find(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
