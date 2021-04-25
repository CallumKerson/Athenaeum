package web

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/gorilla/mux"
)

var mediaPath = "/media/"

// New initializes a new webapp server.
func New(cfg config, bookRetriever bookWebInfoRetriever, podcastBuilder podcastBuilder, podcastSerializer podcastSerializer, logger logging.Logger) (http.Handler, error) {
	template, err := initTemplate(logger, "", cfg.TemplateDir())
	if err != nil {
		return nil, err
	}

	app := &app{
		render: func(wr http.ResponseWriter, tplName string, data interface{}) {
			err = template.ExecuteTemplate(wr, tplName, data)
			if err != nil {
				logger.Errorf("Cannot render template %s, error was %s", tplName, err)
			}
		},
		bookRetriever:     bookRetriever,
		podcastBuilder:    podcastBuilder,
		podcastSerializer: podcastSerializer,
		logger:            logger,
		host:              cfg.Host(),
		podcastAuthor:     model.PodcastAuthor{AuthorName: cfg.PodcastAuthorName(), AuthorEmail: cfg.PodcastAuthorEmail()},
	}

	staticFilesServer := newSafeFileSystemServer(cfg.StaticDir())
	audioFilesServer := newSafeFileSystemServer(cfg.LibraryHome())

	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticFilesServer))
	router.PathPrefix(mediaPath).Handler(http.StripPrefix(mediaPath, audioFilesServer))

	// web app routes
	router.HandleFunc("/", app.indexHandler)
	router.HandleFunc("/books", app.bookListHandler)
	router.HandleFunc("/feed.rss", app.feedHandler)

	router.PathPrefix("/favicon.ico").Handler(faviconHandler(cfg.StaticDir()))

	router.HandleFunc("/book/{author}/{title}", app.bookHandler)
	router.HandleFunc("/author/{author}", app.authorHandler)
	router.HandleFunc("/author/{author}/feed.rss", app.authorFeed)

	router.NotFoundHandler = http.HandlerFunc(app.notFoundHandler)
	return router, nil
}

// Config represents server configuration.
type config interface {
	TemplateDir() string
	StaticDir() string
	LibraryHome() string
	Host() string
	PodcastAuthorName() string
	PodcastAuthorEmail() string
}

func initTemplate(logger logging.Logger, name, path string) (*template.Template, error) {
	apath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(apath)
	if err != nil {
		return nil, err
	}

	logger.Infof("loading templates from '%s'", path)
	tpl := template.New(name)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fp := filepath.Join(apath, f.Name())
		logger.Debugf("parsing template file '%s'", f.Name())
		tpl.New(f.Name()).ParseFiles(fp)
	}

	return tpl, nil
}
