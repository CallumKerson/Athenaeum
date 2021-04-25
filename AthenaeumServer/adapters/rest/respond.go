package rest

import (
	"net/http"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/errors"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/render"
)

func respond(wr http.ResponseWriter, status int, v interface{}) {
	if err := render.JSON(wr, status, v); err != nil {
		if loggable, ok := wr.(errorLogger); ok {
			loggable.Errorf("failed to write data to http ResponseWriter: %s", err)
		}
	}
}

func respondErr(wr http.ResponseWriter, err error) {
	if e, ok := err.(*errors.Error); ok {
		respond(wr, e.Code, e)
		return
	}
	respond(wr, http.StatusInternalServerError, err)
}

type errorLogger interface {
	Errorf(msg string, args ...interface{})
}
