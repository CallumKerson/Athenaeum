package books

import (
	"os"
	"path/filepath"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

type IdProvider interface {
	Get() string
}

type DurationProvider interface {
	GetDuration(fileLocation string) (*time.Duration, error)
}

type FileHasher interface {
	Hash(pathToFile string) ([]byte, error)
}

func NewBookFactory(idProvider IdProvider, store Store,
	libraryRoot string,
	audioMetadataProvider AudioMetadataProvider,
	durationProvider DurationProvider,
	fileHasher FileHasher,
	logger logging.Logger) *Factory {
	logger.Debugf("Creating new BookFactory")
	return &Factory{idProvider, store, libraryRoot, audioMetadataProvider, durationProvider, fileHasher, logger}
}

type Factory struct {
	idProvider            IdProvider
	store                 Store
	libraryRoot           string
	audioMetadataProvider AudioMetadataProvider
	durationProvider      DurationProvider
	fileHasher            FileHasher
	logger                logging.Logger
}

func (factory *Factory) NewBook(title string, subtitle string, author []model.Person, relativePath string, releaseYear int) (*model.Book, error) {
	b := &model.Book{
		Id:              factory.idProvider.Get(),
		Title:           title,
		Subtitle:        subtitle,
		Author:          author,
		ReleaseDateTime: time.Date(releaseYear, time.January, 1, 8, 0, 0, 0, time.UTC),
	}
	hash, err := factory.fileHasher.Hash(filepath.Join(factory.libraryRoot, relativePath))
	if err != nil {
		factory.logger.Errorf("Cannot parse checksum from file %s", filepath.Join(factory.libraryRoot, relativePath))
		return nil, err
	}
	fileSize, err := getFileSize(filepath.Join(factory.libraryRoot, relativePath))
	if err != nil {
		factory.logger.Errorf("Cannot parse filesize from file %s", filepath.Join(factory.libraryRoot, relativePath))
		return nil, err
	}
	fileDuration, err := factory.durationProvider.GetDuration(filepath.Join(factory.libraryRoot, relativePath))
	if err != nil {
		factory.logger.Errorf("Cannot parse file duration from file %s", filepath.Join(factory.libraryRoot, relativePath))
		return nil, err
	}
	b.File = model.AudiobookFile{
		FileExtension: filepath.Ext(relativePath),
		FileLocation:  relativePath,
		FileHash:      hash,
		FileSize:      fileSize,
		FileDuration:  *fileDuration,
	}
	factory.logger.Debugf("Created new book with ID %s at relative location %s", b.Id, b.File.FileLocation)
	b, hasMoved, err := organiseFile(b, factory.libraryRoot)
	if err != nil {
		factory.logger.Errorf("Cannot organise book ID %s", b.Id)
		return b, err
	}
	if hasMoved {
		factory.logger.Debugf("Organised book with ID %s, moved to relative location %s", b.Id, b.File.FileLocation)
	}
	err = factory.store.Add(*b)
	if err != nil {
		factory.logger.Errorf("Cannot add new book with ID %s to store", b.Id)
		return b, err
	}
	factory.logger.Infof("Added new book '%s' with ID %s", b.Description(), b.Id)
	return b, nil
}

func (factory *Factory) Update(b model.Book) error {
	factory.logger.Infof("Updating book with ID %s", b.Id)
	organisedBook, hasMoved, err := organiseFile(&b, factory.libraryRoot)
	if err != nil {
		factory.logger.Errorf("Cannot organise book ID %s", b.Id)
		return err
	}
	if hasMoved {
		factory.logger.Debugf("Organised book with ID %s, moved to relative location %s", b.Id, b.File.FileLocation)
	}
	return factory.store.Add(*organisedBook)
}

// UpdateFrom updates a target book from a source book. The source may be incomplete. File information and ID of the target cannot be overrwritten.
func (factory *Factory) UpdateFrom(target model.Book, source model.Book) error {
	if len(source.Title) > 0 {
		target.Title = source.Title
	}
	if len(source.Summary) > 0 {
		target.Summary = source.Summary
	}
	if len(source.Genre) > 0 {
		target.Genre = source.Genre
	}
	if len(source.Author) > 0 {
		target.Author = source.Author
	}
	if !source.ReleaseDateTime.IsZero() {
		target.ReleaseDateTime = source.ReleaseDateTime
	}
	if source.Series != nil {
		target.Series = source.Series
	}
	return factory.Update(target)
}

func getFileSize(filepath string) (int64, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return 0, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
