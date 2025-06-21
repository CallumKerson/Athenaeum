package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"

	"github.com/CallumKerson/Athenaeum/internal/audiobook"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/static"
)

const (
	StaticPath      = "/static"
	PodcastPath     = "/podcast"
	MediaPath       = "/media"
	PodcastFeedName = "feed.rss"
)

type AudiobookService interface {
	GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error)
	GetAudiobooksByGenre(ctx context.Context, genre audiobooks.Genre) ([]audiobooks.Audiobook, error)
	GetAudiobooksByAuthor(ctx context.Context, author string) ([]audiobooks.Audiobook, error)
	GetAudiobooksByNarrator(ctx context.Context, narrator string) ([]audiobooks.Audiobook, error)
	GetAudiobooksByTag(ctx context.Context, tag string) ([]audiobooks.Audiobook, error)
	UpdateAudiobooks(ctx context.Context) error
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
	audiobookService AudiobookService
	feedConfig       *audiobook.FeedConfig
	mediaRoot        string
	CacheStore       CacheStore
	Log              loggerrific.Logger
	version          string
}

func NewHandler(audiobookService AudiobookService, feedConfig *audiobook.FeedConfig, mediaRoot string, opts ...HandlerOption) *Handler {
	handler := &Handler{
		audiobookService: audiobookService,
		feedConfig:       feedConfig,
		Log:              noOpLogger.New(),
		mediaRoot:        mediaRoot,
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
