package service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alfg/mp4"
	"github.com/pelletier/go-toml/v2"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/m4b"
)

type M4BMetadataReader interface {
	Read(pathToM4BFile string) (*m4b.Metadata, error)
}

type Service struct {
	m4bMetadataReader M4BMetadataReader
	mediaRoot         string
	logger            loggerrific.Logger
}

func New(m4bMetadataReader M4BMetadataReader, mediaRoot string, opts ...Option) *Service {
	svc := &Service{
		m4bMetadataReader: m4bMetadataReader,
		mediaRoot:         mediaRoot,
		logger:            noOpLogger.New(),
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *Service) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	var books []audiobooks.Audiobook
	s.logger.Infoln("Starting scan of directory", s.mediaRoot, "for M4B audiobook files and associated TOML configuration files.")
	err := filepath.WalkDir(s.mediaRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".toml" {
			bookInfo := &audiobooks.Audiobook{}

			err := s.parseM4BInfo(getM4BPathFromTOMLPath(path), bookInfo)
			if errors.Is(err, os.ErrNotExist) {
				s.logger.Warnln("Found audiobook config file", path, "without matching .m4b file")
				return nil
			} else if err != nil {
				s.logger.WithError(err).Warnln("Problem with M4B file", getM4BPathFromTOMLPath(path))
			}
			s.logger.Infoln("Found audiobook config file", path)
			file, err := os.Open(path)
			if err != nil {
				s.logger.WithError(err).Warnln("Problem with TOML file", path)
				return nil
			}
			defer file.Close()
			err = toml.NewDecoder(file).Decode(bookInfo)
			if err != nil {
				s.logger.WithError(err).Warnln("Problem with TOML file", path)
				return nil
			}
			books = append(books, *bookInfo)
		}
		return nil
	})
	return books, err
}

func (s *Service) ScanForNewAndUpdatedAudiobooks(ctx context.Context, books []audiobooks.Audiobook) ([]audiobooks.Audiobook, bool, error) {
	booksMap := listToMap(s.mediaRoot, books)
	s.logger.Infoln("Starting scan of directory", s.mediaRoot, "for M4B audiobook files and associated TOML configuration files.")

	m4bMetadataMap := map[string]string{}
	err := filepath.WalkDir(s.mediaRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".m4b" {
			if _, tomlPathErr := os.Stat(getTOMLPathFromM4BPath(path)); tomlPathErr != nil {
				s.logger.Warnln("Found audiobook file", path, "without matching .toml metadata file")
				return nil
			}
			m4bMetadataMap[path] = getTOMLPathFromM4BPath(path)
		}
		return nil
	})

	changed := false

	for key := range booksMap {
		_, exists := m4bMetadataMap[key]
		if !exists {
			s.logger.Infoln("Audiobook", key, "not found, removing")
			delete(booksMap, key)
			changed = true
		}
	}

	for m4bPath, tomlPath := range m4bMetadataMap {
		currentVal, inMap := booksMap[m4bPath]
		if inMap {
			s.logger.Debugln("Found existing audiobook at", m4bPath)
			book, inMapErr := s.getAudiobook(m4bPath, tomlPath)
			if inMapErr != nil {
				s.logger.WithError(inMapErr).Warnln("Could not process audiobook", m4bPath)
			}
			if !audiobooks.Equal(&currentVal, &book) {
				s.logger.Infoln("Updating existing audiobook at", m4bPath)
				booksMap[m4bPath] = book
				changed = true
			}
		} else {
			s.logger.Infoln("Found new audiobook at", m4bPath)
			book, newBookErr := s.getAudiobook(m4bPath, tomlPath)
			if newBookErr != nil {
				s.logger.WithError(newBookErr).Warnln("Could not process audiobook", m4bPath)
			}
			booksMap[m4bPath] = book
			changed = true
		}
	}

	booksList := make([]audiobooks.Audiobook, 0, len(booksMap))
	for key := range booksMap {
		booksList = append(booksList, booksMap[key])
	}
	return booksList, changed, err
}

func (s *Service) getAudiobook(m4bPath, tomlPath string) (audiobooks.Audiobook, error) {
	audiobook := audiobooks.Audiobook{}

	tomlFile, err := os.Open(tomlPath)
	if err != nil {
		return audiobooks.Audiobook{}, err
	}
	defer tomlFile.Close()
	err = toml.NewDecoder(tomlFile).Decode(&audiobook)
	if err != nil {
		return audiobooks.Audiobook{}, err
	}
	err = s.parseM4BInfo(m4bPath, &audiobook)
	return audiobook, err
}

func (s *Service) parseM4BInfo(m4bPath string, audiobook *audiobooks.Audiobook) (err error) {
	defer s.trackM4BParseTime(time.Now(), m4bPath)
	fInfo, err := os.Stat(m4bPath)
	if err != nil {
		return err
	}
	file, err := os.Open(m4bPath)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := mp4.OpenFromReader(file, fInfo.Size())
	if err != nil {
		return err
	}
	audiobook.Path = strings.TrimPrefix(m4bPath, s.mediaRoot)
	audiobook.FileSize = uint64(fInfo.Size()) //nolint:gosec
	if info.Moov != nil && info.Moov.Mvhd != nil {
		audiobook.Duration = time.Duration(
			(float32(info.Moov.Mvhd.Duration) / float32(info.Moov.Mvhd.Timescale)) * float32(time.Second))
	}
	audiobook.MIMEType = "audio/mp4a-latm"
	return nil
}

func getM4BPathFromTOMLPath(tomlPath string) string {
	return fmt.Sprintf("%s.m4b", strings.TrimSuffix(tomlPath, filepath.Ext(".toml")))
}

func getTOMLPathFromM4BPath(m4bPath string) string {
	return fmt.Sprintf("%s.toml", strings.TrimSuffix(m4bPath, filepath.Ext(".m4b")))
}

func (s *Service) trackM4BParseTime(start time.Time, filename string) {
	elapsed := time.Since(start)
	s.logger.Debugln("Processing", filename, "took", elapsed.String())
}

func listToMap(root string, list []audiobooks.Audiobook) map[string]audiobooks.Audiobook {
	convertedMap := make(map[string]audiobooks.Audiobook, len(list))
	for index := range list {
		convertedMap[filepath.Join(root, list[index].Path)] = list[index]
	}
	return convertedMap
}
