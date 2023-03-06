package rest

import (
	"net/http"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/gorilla/mux"
)

type GenreHolder struct {
	Genres []string `json:"genres"`
}

type genreController struct {
	logging.Logger
	retreve retriever
}

func addGenreAPI(router *mux.Router, ret retriever, logger logging.Logger) {
	genreController := &genreController{
		Logger:  logger,
		retreve: ret,
	}

	logger.Debugf("Adding handlers for genres API")
	router.HandleFunc("/v1/genres/", genreController.getAll).Methods(http.MethodGet)
	router.HandleFunc("/v1/genres/{genre}", genreController.getBooks).Methods(http.MethodGet)

}

func (gController *genreController) getAll(wr http.ResponseWriter, req *http.Request) {
	gController.Debugf("Getting all genres")
	genres, err := gController.retreve.GetAllGenres()
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}

	respond(wr, http.StatusOK, GenreHolder{Genres: genres})
}

func (gController *genreController) getBooks(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	g := vars["genre"]
	gController.Debugf("Getting all books with genre: %s", g)
	genreBooks, err := gController.retreve.GetAllByGenre(g)
	if err != nil {
		respond(wr, http.StatusNotFound, err)
		return
	}

	respond(wr, http.StatusOK, BooksHolder{Books: genreBooks})
}
