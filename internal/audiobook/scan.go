package audiobook

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

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// ScanAll scans the mediaRoot directory for all M4B audiobook files and their TOML metadata files
func ScanAll(ctx context.Context, mediaRoot string, logger loggerrific.Logger) ([]audiobooks.Audiobook, error) {
	var books []audiobooks.Audiobook
	logger.Infoln("Starting scan of directory", mediaRoot, "for M4B audiobook files and associated TOML configuration files.")

	err := filepath.WalkDir(mediaRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".toml" {
			bookInfo := &audiobooks.Audiobook{}

			err := parseM4BInfo(getM4BPathFromTOMLPath(path), mediaRoot, bookInfo, logger)
			if errors.Is(err, os.ErrNotExist) {
				logger.Warnln("Found audiobook config file", path, "without matching .m4b file")
				return nil
			} else if err != nil {
				logger.WithError(err).Warnln("Problem with M4B file", getM4BPathFromTOMLPath(path))
			}
			logger.Infoln("Found audiobook config file", path)
			file, err := os.Open(path)
			if err != nil {
				logger.WithError(err).Warnln("Problem with TOML file", path)
				return nil
			}
			defer file.Close()
			err = toml.NewDecoder(file).Decode(bookInfo)
			if err != nil {
				logger.WithError(err).Warnln("Problem with TOML file", path)
				return nil
			}
			books = append(books, *bookInfo)
		}
		return nil
	})
	return books, err
}

// ScanForUpdates scans for new and updated audiobooks compared to existing ones
func ScanForUpdates(
	ctx context.Context,
	mediaRoot string,
	existing []audiobooks.Audiobook,
	logger loggerrific.Logger,
) ([]audiobooks.Audiobook, bool, error) {
	booksMap := listToMap(mediaRoot, existing)
	logger.Infoln("Starting scan of directory", mediaRoot, "for M4B audiobook files and associated TOML configuration files.")

	m4bMetadataMap := map[string]string{}
	err := filepath.WalkDir(mediaRoot, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".m4b" {
			if _, tomlPathErr := os.Stat(getTOMLPathFromM4BPath(path)); tomlPathErr != nil {
				logger.Warnln("Found audiobook file", path, "without matching .toml metadata file")
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
			logger.Infoln("Audiobook", key, "not found, removing")
			delete(booksMap, key)
			changed = true
		}
	}

	for m4bPath, tomlPath := range m4bMetadataMap {
		currentVal, inMap := booksMap[m4bPath]
		if inMap {
			logger.Debugln("Found existing audiobook at", m4bPath)
			book, inMapErr := getAudiobook(m4bPath, tomlPath, mediaRoot, logger)
			if inMapErr != nil {
				logger.WithError(inMapErr).Warnln("Could not process audiobook", m4bPath)
			}
			if !audiobooks.Equal(&currentVal, &book) {
				logger.Infoln("Updating existing audiobook at", m4bPath)
				booksMap[m4bPath] = book
				changed = true
			}
		} else {
			logger.Infoln("Found new audiobook at", m4bPath)
			book, newBookErr := getAudiobook(m4bPath, tomlPath, mediaRoot, logger)
			if newBookErr != nil {
				logger.WithError(newBookErr).Warnln("Could not process audiobook", m4bPath)
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

func getAudiobook(m4bPath, tomlPath, mediaRoot string, logger loggerrific.Logger) (audiobooks.Audiobook, error) {
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
	err = parseM4BInfo(m4bPath, mediaRoot, &audiobook, logger)
	return audiobook, err
}

func parseM4BInfo(m4bPath, mediaRoot string, audiobook *audiobooks.Audiobook, logger loggerrific.Logger) (err error) {
	defer trackM4BParseTime(time.Now(), m4bPath, logger)
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
	audiobook.Path = strings.TrimPrefix(m4bPath, mediaRoot)
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

func trackM4BParseTime(start time.Time, filename string, logger loggerrific.Logger) {
	elapsed := time.Since(start)
	logger.Debugln("Processing", filename, "took", elapsed.String())
}

func listToMap(root string, list []audiobooks.Audiobook) map[string]audiobooks.Audiobook {
	convertedMap := make(map[string]audiobooks.Audiobook, len(list))
	for index := range list {
		convertedMap[filepath.Join(root, list[index].Path)] = list[index]
	}
	return convertedMap
}
