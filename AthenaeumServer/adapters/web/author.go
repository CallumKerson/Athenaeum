package web

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/gorilla/mux"
)

type AuthorData struct {
	Books      []model.Book
	AuthorName string
}

func (app app) authorHandler(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	author := vars["author"]
	fetchedBooks, err := app.bookRetriever.GetAllByAuthor(author)
	if err != nil || len(fetchedBooks) == 0 {
		app.logger.Errorf("Cannot get author's books for request %s", req.RequestURI)
		app.notFoundHandler(wr, req)
	} else {
		data := AuthorData{
			Books:      fetchedBooks,
			AuthorName: author,
		}
		app.render(wr, "author.html", data)
	}
}

func (app app) authorFeed(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	author := vars["author"]
	fetchedBooks, err := app.bookRetriever.GetAllByAuthor(author)
	if err != nil || len(fetchedBooks) == 0 {
		app.logger.Errorf("Cannot get author's books for request %s", req.RequestURI)
		app.notFoundHandler(wr, req)
	} else {
		sort.SliceStable(fetchedBooks, func(i, j int) bool {
			return fetchedBooks[i].ReleaseDateTime.After(fetchedBooks[j].ReleaseDateTime)
		})
		podcast, err := app.podcastBuilder.NewPodcast(
			fmt.Sprintf("Books by %s", author),
			app.host,
			fmt.Sprintf("A Podcast Feed Containing Books by %s", author),
			app.podcastAuthor,
			fetchedBooks[0].ReleaseDateTime,
			true,
			app.mediaHost(),
			fetchedBooks)
		if err != nil {
			app.logger.Errorf("Build podcast for request %s", req.RequestURI)
		}
		app.podcastSerializer.String(*podcast)
		_, _ = wr.Write([]byte(app.podcastSerializer.String(*podcast)))
	}
}
