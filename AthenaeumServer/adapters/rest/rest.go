package rest

import (
	"encoding/json"
	"net/http"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/errors"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/render"
	"github.com/gorilla/mux"
)

// New initializes the server with routes exposing the given usecases.
func New(logger logging.Logger, bookRetriever retriever, bookUpdater updater, scan scanner) http.Handler {

	// setup router with default handlers
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	// setup api endpoints
	addBooksAPI(router, bookRetriever, bookUpdater, logger)
	addScannerAPI(router, scan, logger)
	addGenreAPI(router, bookRetriever, logger)
	// addPostsAPI(router, postPub, postsRet, logger)

	return router
}

func notFoundHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusNotFound, errors.ResourceNotFound("path", req.URL.Path))
}

func methodNotAllowedHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusMethodNotAllowed, errors.ResourceNotFound("path", req.URL.Path))
}

func readRequest(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return errors.Validation("Failed to read request body")
	}

	return nil
}
