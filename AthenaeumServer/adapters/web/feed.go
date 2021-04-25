package web

import (
	"net/http"
	"sort"
)

func (app app) feedHandler(wr http.ResponseWriter, req *http.Request) {

	books, err := app.bookRetriever.GetAll()
	if err != nil {
		app.logger.Errorf("Cannot get all books for request %s", req.RequestURI)
	}
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].ReleaseDateTime.After(books[j].ReleaseDateTime)
	})

	podcast, err := app.podcastBuilder.NewPodcast("All Books", app.host, "A Podcast Feed Containing All Books",
		app.podcastAuthor, books[0].ReleaseDateTime, true, app.mediaHost(), books)
	if err != nil {
		app.logger.Errorf("Cannot create podcast feed for request %s", req.RequestURI)
	}
	app.podcastSerializer.String(*podcast)
	_, _ = wr.Write([]byte(app.podcastSerializer.String(*podcast)))
}
