package web

import "net/http"

func (app app) notFoundHandler(wr http.ResponseWriter, r *http.Request) {
	app.logger.Errorf("Not Found!")
	wr.WriteHeader(404)
	app.render(wr, "404.html", nil)
}
