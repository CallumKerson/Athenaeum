package books

import (
	"path/filepath"
	"strings"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

// Creates a new Book from an audiofile at the relative path to the library root
func (factory *Factory) NewBookFromAudioMetadata(relativeFilePath string) (*model.Book, error) {

	title, subtitle, err := factory.parseTitle(filepath.Join(factory.libraryRoot, relativeFilePath))
	if err != nil {
		factory.logger.Errorf("Cannot get title from file %s", relativeFilePath)
		return nil, err
	}
	author, err := factory.combineArtists(filepath.Join(factory.libraryRoot, relativeFilePath))
	if err != nil {
		factory.logger.Errorf("Cannot get author from file %s", relativeFilePath)
		return nil, err
	}
	year, err := factory.audioMetadataProvider.GetYear(filepath.Join(factory.libraryRoot, relativeFilePath))
	if err != nil {
		factory.logger.Errorf("Cannot get year from file %s", relativeFilePath)
		return nil, err
	}
	return factory.NewBook(
		title,
		subtitle,
		ParseAuthors(author),
		relativeFilePath,
		year)
}

// This checks album artist first, and if there is no album artist it uses artist
func (factory *Factory) combineArtists(file string) (string, error) {
	albumArtist, err := factory.audioMetadataProvider.GetAlbumArtist(file)
	if err != nil || albumArtist == "" {
		factory.logger.Warnf("No album artist found in file %s", file)
		return factory.audioMetadataProvider.GetArtist(file)
	} else {
		return albumArtist, err
	}
}

// parseTitle gets the title from the metadata provider returns a title and, if a colon is present, a subtitle
func (factory *Factory) parseTitle(file string) (string, string, error) {
	metadataTitle, err := factory.audioMetadataProvider.GetTitle(file)
	if err != nil {
		factory.logger.Errorf("Cannot get title from metadata of file %s", file)
		return "", "", err
	}
	if strings.Contains(metadataTitle, ":") {
		ss := strings.SplitN(metadataTitle, ":", 2)
		if len(ss) < 2 {
			return metadataTitle, "", nil
		} else {
			for i := range ss {
				ss[i] = strings.TrimSpace(ss[i])
			}
			return ss[0], ss[1], nil
		}
	} else {
		return metadataTitle, "", nil
	}
}

// Provides an interface to get relevant metadata from audio metadata
type AudioMetadataProvider interface {
	GetTitle(fileLocation string) (string, error)
	GetArtist(fileLocation string) (string, error)
	GetAlbumArtist(fileLocation string) (string, error)
	GetYear(fileLocation string) (int, error)
}
