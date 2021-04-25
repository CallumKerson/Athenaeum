package podcasts

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

type Factory struct {
	libraryRoot string
	logger      logging.Logger
}

func NewPodcastFactory(libraryRoot string, logger logging.Logger) *Factory {
	return &Factory{libraryRoot, logger}
}

func (factory *Factory) NewPodcast(title string, link string, description string, author model.PodcastAuthor, pubDate time.Time,
	explicit bool, mediaHostPrefix string, books []model.Book) (*model.PodcastOfBooks, error) {
	p := &model.PodcastOfBooks{
		Title:           title,
		Link:            link,
		Description:     description,
		Author:          author,
		ExplicitStatus:  explicit,
		LastBuildTime:   nil,
		PublicationDate: &pubDate,
		Category:        model.Category{MainCategory: "Arts", SubCategories: []string{"Books"}},
	}
	p.Items = factory.booksToPodcastItem(mediaHostPrefix, books)
	return p, nil
}

func (factory *Factory) booksToPodcastItem(mediaHostPrefix string, books []model.Book) []model.PodcastFeedItem {
	var items []model.PodcastFeedItem

	for _, book := range books {
		unescapedPath := strings.Split(book.File.FileLocation, "/")
		escapedPath := []string{}
		for _, pathElement := range unescapedPath {
			escapedPath = append(escapedPath, url.PathEscape(pathElement))
		}
		bookRelativeURL := filepath.Join(escapedPath...)
		item := model.PodcastFeedItem{
			GUID:              book.Id,
			Title:             book.Title,
			Subtitle:          book.AuthorString(),
			PublicationDate:   book.ReleaseDateTime,
			Description:       book.Description(),
			Summary:           book.FullHTMLDescription(),
			DurationInSeconds: int64(book.File.FileDuration.Seconds()),
			Enclosure: model.Enclosure{
				URL:    fmt.Sprintf("%s%s", mediaHostPrefix, bookRelativeURL),
				Length: book.File.FileSize,
			},
		}
		items = append(items, item)
	}

	return items
}
