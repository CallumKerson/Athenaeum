package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/static"
)

type AudiobooksPodcastService interface {
	WriteAllAudiobooksFeed(context.Context, io.Writer) error
	WriteGenreAudiobookFeed(context.Context, audiobooks.Genre, io.Writer) error
	WriteAuthorAudiobookFeed(context.Context, string, io.Writer) (bool, error)
	WriteNarratorAudiobookFeed(context.Context, string, io.Writer) (bool, error)
	UpdateFeeds(context.Context) error
	IsReady(ctx context.Context) bool
}

type CacheStore interface {
	Get(key uint64) ([]byte, bool)
	Set(key uint64, content []byte, expiration time.Time)
	Release(key uint64)
	ReleaseAll()
	GetTTL() time.Duration
}

type Handler struct {
	*mux.Router
	PodcastService   AudiobooksPodcastService
	CacheStore       CacheStore
	Log              loggerrific.Logger
	version          string
	mediaRoot        string
	mediaServePath   string
	staticServePath  string
	podcastServePath string
	mainFeedPath     string
}

func NewHandler(podcastService AudiobooksPodcastService, cacheStore CacheStore,
	logger loggerrific.Logger, opts ...HandlerOption) *Handler {
	handler := &Handler{
		PodcastService: podcastService,
		CacheStore:     cacheStore,
		Log:            logger,
	}
	for _, opt := range opts {
		opt(handler)
	}
	handler.podcastServePath = "/podcast"
	handler.mainFeedPath = "/feed.rss"
	handler.Router = mux.NewRouter()
	handler.mapRoutes()
	handler.Use(TimeoutMiddleware)
	return handler
}

func (h *Handler) mapRoutes() {
	h.HandleFunc("/health", healthCheck)
	h.HandleFunc("/ready", h.readiness)
	h.HandleFunc("/version", h.printVersion)

	middleware := NewMiddlewares(h.Log, h.CacheStore)

	podcastSubrouter := h.PathPrefix(h.podcastServePath).Subrouter()
	if middleware.CacheStore != nil {
		h.Log.Infoln("Caching enabled on", h.podcastServePath, "endpoints is enabled with at TTL of", middleware.CacheStore.GetTTL().String())
	}
	podcastSubrouter.Use(middleware.LoggingMiddleware, middleware.CachingMiddleware)
	podcastSubrouter.HandleFunc(h.mainFeedPath, h.getFeed)
	podcastSubrouter.HandleFunc(fmt.Sprintf("/genre/{genre}%s", h.mainFeedPath), h.getGenreFeed)
	podcastSubrouter.HandleFunc(fmt.Sprintf("/authors/{author}%s", h.mainFeedPath), h.getAuthorFeed)
	podcastSubrouter.HandleFunc(fmt.Sprintf("/narrators/{narrator}%s", h.mainFeedPath), h.getNarratorFeed)

	updateRouter := h.PathPrefix("/update").Subrouter()
	updateRouter.Use(SevereRateLimitMiddleware, middleware.LoggingMiddleware)
	updateRouter.HandleFunc("", h.updateAudiobooks)

	mediaFS := http.StripPrefix(h.mediaServePath, http.FileServer(http.Dir(h.mediaRoot)))
	h.Log.Infoln("Serving media files from local path", h.mediaRoot, "at", h.mediaServePath)
	h.Router.PathPrefix(h.mediaServePath).Handler(middleware.LoggingMiddleware(mediaFS))

	staticFS := http.StripPrefix(h.staticServePath, http.FileServer(http.FS(static.Assets)))
	h.Log.Infoln("Serving static files at", h.staticServePath)
	h.Router.PathPrefix(h.staticServePath).Handler(middleware.LoggingMiddleware(staticFS))

	h.Handle("/", middleware.LoggingMiddleware(http.HandlerFunc(h.serveHTML)))
}
