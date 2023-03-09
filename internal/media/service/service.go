package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type Service struct {
	mediaRoot string
	logger    loggerrific.Logger
}

func New(pathToMediaRoot string, logger loggerrific.Logger) *Service {
	return &Service{mediaRoot: pathToMediaRoot, logger: logger}
}

func (s *Service) GetAudiobooks() ([]audiobooks.Audiobook, error) {
	var books []audiobooks.Audiobook
	err := filepath.WalkDir(s.mediaRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".toml" {
			hasM4BFile, m4bPathWithoutRoot, m4bSize := s.parseM4BInfo(path)
			if !hasM4BFile {
				s.logger.Warnln("Found audiobook config file", path, "without matching .m4b file")
				return nil
			}
			s.logger.Infoln("Found audiobook config file", path)
			bookInfo := &audiobooks.Audiobook{}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			err = toml.NewDecoder(file).Decode(bookInfo)
			if err != nil {
				return err
			}
			bookInfo.Path = m4bPathWithoutRoot
			bookInfo.FileSize = uint64(m4bSize)
			bookInfo.MIMEType = "audio/mp4a-latm"
			books = append(books, *bookInfo)
		}
		return nil
	})
	return books, err
}

// func hasMatchingM4BFile(path string) bool {
// 	expectedAudiobookPath := fmt.Sprintf("%s.m4b", strings.TrimSuffix(path, filepath.Ext(".toml")))
// 	if _, err := os.Stat(expectedAudiobookPath); errors.Is(err, os.ErrNotExist) {
// 		return false
// 	}
// 	return true
// }

func (s *Service) getM4BPath(path string) string {
	return strings.TrimPrefix(path, s.mediaRoot)
}

func (s *Service) parseM4BInfo(tomlPath string) (exists bool, pathWithoutRoot string, size int64) {
	expectedAudiobookPath := fmt.Sprintf("%s.m4b", strings.TrimSuffix(tomlPath, filepath.Ext(".toml")))
	fInfo, err := os.Stat(expectedAudiobookPath)
	if err != nil {
		return false, "", 0
	}
	return true, s.getM4BPath(expectedAudiobookPath), fInfo.Size()
}
