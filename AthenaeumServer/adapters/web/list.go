package web

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

type booksPageData struct {
	BooksData []singleBookData
}

type singleBookData struct {
	Title       template.HTML
	Author      template.HTML
	ReleaseDate string
}

func (app app) bookListHandler(wr http.ResponseWriter, req *http.Request) {
	fetchedBooks, err := app.bookRetriever.GetAll()
	if err != nil {
		app.logger.Errorf("Cannot list all books for request %s", req.RequestURI)
	}
	sort.SliceStable(fetchedBooks, func(i, j int) bool {
		return fetchedBooks[i].Author[0].FamilyName < fetchedBooks[j].Author[0].FamilyName
	})
	var sbd []singleBookData

	for _, fb := range fetchedBooks {
		sbd = append(sbd,
			singleBookData{
				Title:       titleToHTML(fb.Title, fb.AuthorString()),
				Author:      authorToHTML(fb.Author),
				ReleaseDate: fb.ReleaseDate(),
			})
	}

	data := booksPageData{
		BooksData: sbd,
	}
	app.render(wr, "books.html", data)
}

func titleToHTML(title string, authorString string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"/book/%s/%s\">%s</a>", authorString, title, title))
}

func authorToHTML(authors []model.Person) template.HTML {

	authorLinkFormat := "<a href=\"/author/%s\">%s</a>"
	if len(authors) == 1 {
		return template.HTML(fmt.Sprintf(authorLinkFormat, authors[0].String(), authors[0].String()))
	} else {
		andAuthor := authors[len(authors)-1]
		commaSeparatedAuthorStrings := []string{}
		for _, author := range authors[:len(authors)-1] {
			commaSeparatedAuthorStrings = append(commaSeparatedAuthorStrings, fmt.Sprintf(authorLinkFormat, author.String(), author.String()))
		}
		return template.HTML(fmt.Sprintf("%s and %s", strings.Join(commaSeparatedAuthorStrings[:], ", "), fmt.Sprintf(authorLinkFormat, andAuthor.String(), andAuthor.String())))
	}
}
