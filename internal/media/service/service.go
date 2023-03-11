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

	"github.com/abema/go-mp4"
	"github.com/pelletier/go-toml/v2"
	"github.com/sunfish-shogi/bufseekio"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type Service struct {
	mediaRoot string
	logger    loggerrific.Logger
}

func New(logger loggerrific.Logger, opts ...Option) *Service {
	svc := &Service{logger: logger}
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

			err := s.parseM4BInfo(path, bookInfo)
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

func (s *Service) parseM4BInfo(tomlPath string, audiobook *audiobooks.Audiobook) (err error) {
	expectedAudiobookPath := getM4BPathFromTOMLPath(tomlPath)
	fInfo, err := os.Stat(expectedAudiobookPath)
	if err != nil {
		return err
	}
	file, err := os.Open(expectedAudiobookPath)
	if err != nil {
		return err
	}
	info, err := mp4.Probe(bufseekio.NewReadSeeker(file, 1024, 4))
	if err != nil {
		return err
	}
	audiobook.Path = strings.TrimPrefix(expectedAudiobookPath, s.mediaRoot)
	audiobook.FileSize = uint64(fInfo.Size())
	audiobook.Duration = time.Duration((float32(info.Duration) / float32(info.Timescale)) * float32(time.Second))
	audiobook.MIMEType = "audio/mp4a-latm"
	return nil
}

func getM4BPathFromTOMLPath(tomlPath string) string {
	return fmt.Sprintf("%s.m4b", strings.TrimSuffix(tomlPath, filepath.Ext(".toml")))
}
