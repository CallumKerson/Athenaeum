package books

import "github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"

type Store interface {
	Add(book model.Book) error
	Remove(book model.Book) error
	GetAll() ([]model.Book, error)
	Get(id string) (*model.Book, error)
}
