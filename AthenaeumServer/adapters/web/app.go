package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

type bookWebInfoRetriever interface {
	GetAllTitles() ([]string, error)
	GetAll() ([]model.Book, error)
	GetAllByAuthor(author string) ([]model.Book, error)
	GetByTitleAndAuthor(title string, author string) (*model.Book, error)
}

type podcastBuilder interface {
	NewPodcast(title string, link string,
		description string,
		author model.PodcastAuthor,
		lastPublishedTime time.Time,
		explicit bool,
		mediaHostPrefix string,
		books []model.Book) (*model.PodcastOfBooks, error)
}

type podcastSerializer interface {
	String(podcast model.PodcastOfBooks) string
}

type app struct {
	render            func(wr http.ResponseWriter, template string, data interface{})
	bookRetriever     bookWebInfoRetriever
	podcastBuilder    podcastBuilder
	podcastSerializer podcastSerializer
	logger            logging.Logger
	host              string
	podcastAuthor     model.PodcastAuthor
}

func (a app) mediaHost() string {
	return fmt.Sprintf("%s%s", strings.TrimSuffix(a.host, "/"), mediaPath)
}
