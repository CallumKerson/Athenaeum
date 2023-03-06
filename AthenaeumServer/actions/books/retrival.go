package books

import (
	"bytes"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/errors"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

func NewRetriever(store Store, logger logging.Logger) *Retriever {
	return &Retriever{store, logger}
}

type Retriever struct {
	store  Store
	logger logging.Logger
}

func (retriever *Retriever) Get(id string) (*model.Book, error) {
	retriever.logger.Debugf("Getting id %s from Store", id)
	return retriever.store.Get(id)
}

func (retriever *Retriever) GetByTitle(title string) (*model.Book, error) {
	retriever.logger.Debugf("Getting title %s from Store")
	books, err := retriever.store.GetAll()
	if err != nil {
		retriever.logger.Errorf("Cannot get all books to search for title %s", title)
		return nil, err
	}
	var bookWithTitle *model.Book

	for _, bookFound := range books {
		if bookFound.Title == title {
			bookWithTitle = &bookFound
			break
		}
	}

	if bookWithTitle == nil {
		return nil, errors.ResourceNotFound("Book", title)
	}

	return bookWithTitle, nil
}

func (retriever *Retriever) GetByTitleAndAuthor(title string, author string) (*model.Book, error) {
	retriever.logger.Debugf("Getting title %s from Store")
	books, err := retriever.store.GetAll()
	if err != nil {
		retriever.logger.Errorf("Cannot get all books to search for title %s", title)
		return nil, err
	}
	var bookWithTitle *model.Book

	for _, bookFound := range books {
		if bookFound.Title == title {
			bookWithTitle = &bookFound
			break
		}
	}

	if bookWithTitle == nil {
		retriever.logger.Errorf("Could not find book with title %s", title)
		return nil, errors.ResourceNotFound("Book", title)
	}

	if bookWithTitle.AuthorString() == author {
		return bookWithTitle, nil
	} else {
		if !HasAuthor(bookWithTitle.Author, author) {
			retriever.logger.Errorf("Book has wrong author, has %s and expected %s", bookWithTitle.AuthorString(), author)
			return nil, errors.ResourceNotFound("Book", title)
		}

		return bookWithTitle, nil
	}
}

func (retriever *Retriever) GetByFileHash(hash []byte) (*model.Book, error) {
	retriever.logger.Debugf("Getting title %s from Store")
	books, err := retriever.store.GetAll()
	if err != nil {
		retriever.logger.Errorf("Cannot get all books to search for hash %s", string(hash))
		return nil, err
	}
	var bookWithHash *model.Book

	for _, bookFound := range books {
		if bytes.Equal(bookFound.File.FileHash, hash) {
			bookWithHash = &bookFound
			break
		}
	}

	if bookWithHash == nil {
		return nil, errors.ResourceNotFound("Book", string(hash))
	}

	return bookWithHash, nil
}

func (retriever *Retriever) GetAll() ([]model.Book, error) {
	retriever.logger.Debugf("Getting all books from Store")
	return retriever.store.GetAll()
}

func (retriever *Retriever) GetAllByAuthor(author string) ([]model.Book, error) {
	retriever.logger.Debugf("Getting all books by author %s from Store", author)
	allBooks, err := retriever.store.GetAll()
	var booksByAuthor []model.Book
	if err != nil {
		return booksByAuthor, err
	}
	for _, book := range allBooks {
		if HasAuthor(book.Author, author) {
			booksByAuthor = append(booksByAuthor, book)
		}
	}
	if len(booksByAuthor) < 1 {
		return booksByAuthor, errors.ResourceNotFound("model.Book", "authors")
	}
	return booksByAuthor, nil
}

func (retriever *Retriever) GetAllByGenre(genre string) ([]model.Book, error) {
	retriever.logger.Debugf("Getting all books with genre %s from Store", genre)
	allBooks, err := retriever.store.GetAll()
	var booksWithGenre []model.Book
	if err != nil {
		return booksWithGenre, err
	}
	for _, book := range allBooks {
		if stringInSlice(genre, book.Genre) {
			booksWithGenre = append(booksWithGenre, book)
		}
	}
	return booksWithGenre, nil
}

func (retriever *Retriever) GetAllGenres() ([]string, error) {
	retriever.logger.Debugf("Getting all genres Store")
	allBooks, err := retriever.store.GetAll()
	genres := []string{}
	if err != nil {
		return genres, err
	}
	for _, book := range allBooks {
		if len(book.Genre) <= 0 {
			genres = append(genres, book.Genre...)
		}
	}
	return removeDuplicateValues(genres), nil
}

func (retriever *Retriever) GetAllTitles() ([]string, error) {
	var titles []string
	retriever.logger.Debugf("Getting all titles from Store")
	fetchedBooks, err := retriever.GetAll()
	if err == nil {
		for _, b := range fetchedBooks {
			titles = append(titles, b.Title)
		}
	}
	retriever.logger.Debugf("Retrieved titles")
	return titles, err
}

func (retriever *Retriever) GetAllLocations() ([]string, error) {
	var locations []string
	retriever.logger.Debugf("Getting all locations from Store")
	fetchedBooks, err := retriever.GetAll()
	if err == nil {
		for _, b := range fetchedBooks {
			locations = append(locations, b.File.FileLocation)
		}
	}
	retriever.logger.Debugf("Retrieved files")
	return locations, err
}

func (retriever *Retriever) GetAllHashes() ([][]byte, error) {
	var hashes [][]byte
	retriever.logger.Debugf("Getting all file hashes from Store")
	fetchedBooks, err := retriever.GetAll()
	if err == nil {
		for _, b := range fetchedBooks {
			hashes = append(hashes, b.File.FileHash)
		}
	}
	retriever.logger.Debugf("Retrieved all file hashes")
	return hashes, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
