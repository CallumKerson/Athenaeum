package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type BookData struct {
	HtmlDescription template.HTML
	PageTitle       string
}

func (app app) bookHandler(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title := vars["title"]
	author := vars["author"]
	fetchedBook, err := app.bookRetriever.GetByTitleAndAuthor(title, author)
	if err != nil {
		app.logger.Errorf("Cannot get book with title %s and author %s from request %s", title, author, req.RequestURI)
		app.notFoundHandler(wr, req)
	} else {
		data := BookData{
			HtmlDescription: template.HTML(fetchedBook.FullHTMLDescription()),
			PageTitle:       fetchedBook.Title,
		}
		app.render(wr, "book.html", data)
	}
}
