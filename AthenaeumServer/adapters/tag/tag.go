package tag

import (
	"os"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/dhowden/tag"
)

func NewTagProvider(logger logging.Logger) *TagProvider {
	return &TagProvider{logger}
}

type TagProvider struct {
	logger logging.Logger
}

func (tagProvider *TagProvider) GetTitle(fileLocation string) (string, error) {
	metadata, err := tagProvider.getAudioMetadata(fileLocation)
	if err != nil {
		tagProvider.logger.Errorf("Cannot Get Title metadata from file %s", fileLocation)
		return "", err
	}
	return metadata.Title(), nil
}

func (tagProvider *TagProvider) GetArtist(fileLocation string) (string, error) {
	metadata, err := tagProvider.getAudioMetadata(fileLocation)
	if err != nil {
		tagProvider.logger.Errorf("Cannot Get Artist metadata from file %s", fileLocation)
		return "", err
	}
	return metadata.Artist(), nil
}

func (tagProvider *TagProvider) GetAlbumArtist(fileLocation string) (string, error) {
	metadata, err := tagProvider.getAudioMetadata(fileLocation)
	if err != nil {
		tagProvider.logger.Errorf("Cannot Get Album Artist metadata from file %s", fileLocation)
		return "", err
	}
	return metadata.AlbumArtist(), nil
}

func (tagProvider *TagProvider) GetYear(fileLocation string) (int, error) {
	metadata, err := tagProvider.getAudioMetadata(fileLocation)
	if err != nil {
		tagProvider.logger.Errorf("Cannot Get Year metadata from file %s", fileLocation)
		return 0, err
	}
	return metadata.Year(), nil
}

func (tagProvider *TagProvider) getAudioMetadata(fileLocation string) (tag.Metadata, error) {
	file, err := os.Open(fileLocation)
	if err != nil {
		tagProvider.logger.Errorf("Cannot open file %s for reading", fileLocation)
		return nil, err
	}
	defer file.Close()

	// Metadata from file
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		tagProvider.logger.Errorf("Cannot get audio metadata from file %s", fileLocation)
		return nil, err
	}
	return metadata, nil
}
