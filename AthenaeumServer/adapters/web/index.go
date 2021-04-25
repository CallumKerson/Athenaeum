package web

import "net/http"

type IndexPageData struct {
	PageTitle string
}

func (app app) indexHandler(wr http.ResponseWriter, req *http.Request) {
	data := IndexPageData{
		PageTitle: "Audiobooks",
	}
	app.logger.Debugf("Rendering index")
	app.render(wr, "index.html", data)
}
