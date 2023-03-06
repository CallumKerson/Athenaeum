package rest

import (
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

type BooksHolder struct {
	Books []model.Book `json:"books"`
}

type retriever interface {
	Get(id string) (*model.Book, error)
	GetByTitle(name string) (*model.Book, error)
	GetAll() ([]model.Book, error)
	GetAllGenres() ([]string, error)
	GetAllByGenre(genre string) ([]model.Book, error)
}

type updater interface {
	UpdateFrom(target model.Book, source model.Book) error
}
