package rest

import (
	"net/http"
	"sort"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/gorilla/mux"
)

type BooksHolder struct {
	Books []model.Book `json:"books"`
}

func addBooksAPI(router *mux.Router, ret retriever, up updater, logger logging.Logger) {
	bookController := &bookController{
		Logger:  logger,
		retreve: ret,
		upd:     up,
	}

	logger.Debugf("Adding handlers for books API")
	router.HandleFunc("/v1/books/{id}", bookController.get).Methods(http.MethodGet)
	router.HandleFunc("/v1/books/{id}", bookController.update).Methods(http.MethodPatch)
	router.HandleFunc("/v1/books/titles/{title}", bookController.getByTitle).Methods(http.MethodGet)
	router.HandleFunc("/v1/books/", bookController.getAll).Methods(http.MethodGet)
}

type bookController struct {
	logging.Logger
	retreve retriever
	upd     updater
}

func (bController *bookController) get(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookId := vars["id"]
	bController.Debugf("Getting book with id %s", bookId)
	book, err := bController.retreve.Get(bookId)
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}

	respond(wr, http.StatusOK, book)
}

func (bController *bookController) update(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookId := vars["id"]
	book, err := bController.retreve.Get(bookId)
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}
	bController.Debugf("Attempting to update book with id %s", bookId)

	bookDetails := model.Book{}
	if err := readRequest(req, &bookDetails); err != nil {
		bController.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, err)
		return
	}
	bController.Debugf("Read new book details from request")

	err = bController.upd.UpdateFrom(*book, bookDetails)
	if err != nil {
		bController.Warnf("Failed to update book with id %s: %s", book.Id, err)
		respond(wr, http.StatusBadRequest, err)
		return
	}
	bController.Debugf("Updated new book details from request")
	book, err = bController.retreve.Get(bookId)
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}

	respond(wr, http.StatusOK, book)
}

func (bController *bookController) getByTitle(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	book, err := bController.retreve.GetByTitle(vars["title"])
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}

	respond(wr, http.StatusOK, book)
}

func (bController *bookController) getAll(wr http.ResponseWriter, req *http.Request) {
	allBooks, err := bController.retreve.GetAll()
	if err != nil {
		respondErr(wr, err)
		return
	}
	sort.SliceStable(allBooks, func(i, j int) bool {
		return allBooks[i].Author[0].FamilyName < allBooks[j].Author[0].FamilyName
	})
	booksData := BooksHolder{Books: allBooks}
	respond(wr, http.StatusOK, booksData)
}

type retriever interface {
	Get(id string) (*model.Book, error)
	GetByTitle(name string) (*model.Book, error)
	GetAll() ([]model.Book, error)
}

type updater interface {
	UpdateFrom(target model.Book, source model.Book) error
}
