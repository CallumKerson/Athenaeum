package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/static"
)

const (
	StaticPath      = "/static"
	PodcastPath     = "/podcast"
	MediaPath       = "/media"
	PodcastFeedName = "feed.rss"
)

type AudiobooksPodcastService interface {
	WriteAllAudiobooksFeed(context.Context, io.Writer) error
	WriteGenreAudiobookFeed(context.Context, audiobooks.Genre, io.Writer) error
	WriteAuthorAudiobookFeed(context.Context, string, io.Writer) (bool, error)
	WriteNarratorAudiobookFeed(context.Context, string, io.Writer) (bool, error)
	WriteTagAudiobookFeed(context.Context, string, io.Writer) (bool, error)
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
	*chi.Mux
	PodcastService AudiobooksPodcastService
	mediaRoot      string
	CacheStore     CacheStore
	Log            loggerrific.Logger
	version        string
}

func NewHandler(podcastService AudiobooksPodcastService, mediaRoot string, opts ...HandlerOption) *Handler {
	handler := &Handler{
		PodcastService: podcastService,
		Log:            noOpLogger.New(),
		mediaRoot:      mediaRoot,
	}
	for _, opt := range opts {
		opt(handler)
	}
	handler.Mux = chi.NewRouter()
	handler.mapRoutes()
	return handler
}

func (h *Handler) mapRoutes() {
	h.Use(chiMiddleware.RequestID, chiMiddleware.Recoverer, GetLoggingMiddleware(h.Log), TimeoutMiddleware)

	h.HandleFunc("/health", healthCheck)
	h.HandleFunc("/ready", h.readiness)
	h.HandleFunc("/version", h.printVersion)

	h.Route(PodcastPath, func(router chi.Router) {
		if h.CacheStore != nil {
			h.Log.Infoln("Caching enabled on", PodcastPath, "endpoints is enabled with at TTL of", h.CacheStore.GetTTL().String())
			router.Use(GetCachingMiddleware(h.CacheStore))
		}
		router.HandleFunc(fmt.Sprintf("/genre/{genre}/%s", PodcastFeedName), h.getGenreFeed)
		router.HandleFunc(fmt.Sprintf("/authors/{author}/%s", PodcastFeedName), h.getAuthorFeed)
		router.HandleFunc(fmt.Sprintf("/narrators/{narrator}/%s", PodcastFeedName), h.getNarratorFeed)
		router.HandleFunc(fmt.Sprintf("/tags/{tag}/%s", PodcastFeedName), h.getTagFeed)
		router.HandleFunc(fmt.Sprintf("/%s", PodcastFeedName), h.getFeed)
	})

	h.Handle("/update", SevereRateLimitMiddleware(http.HandlerFunc(h.updateAudiobooks)))

	h.Log.Infoln("Serving media files from local path", h.mediaRoot, "at", MediaPath)
	h.routeFileServer(MediaPath, http.Dir(h.mediaRoot))

	h.Log.Infoln("Serving static files at", StaticPath)
	h.routeFileServer(StaticPath, http.FS(static.Assets))

	h.HandleFunc("/", h.serveHTML)
}

func (h *Handler) routeFileServer(path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		h.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	h.Get(path, func(w http.ResponseWriter, r *http.Request) {
		routeContext := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(routeContext.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
