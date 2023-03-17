package http

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/static"
)

type AudiobooksPodcastService interface {
	WriteAllAudiobooksFeed(ctx context.Context, w io.Writer) error
	IsReady(ctx context.Context) bool
}

type AudiobooksUpdateService interface {
	UpdateAudiobooks(ctx context.Context) error
	IsReady(ctx context.Context) bool
}

type Handler struct {
	*mux.Router
	PodcastService   AudiobooksPodcastService
	UpdateService    AudiobooksUpdateService
	Log              loggerrific.Logger
	version          string
	mediaRoot        string
	mediaServePath   string
	staticServePath  string
	podcastServePath string
	mainFeedPath     string
}

func NewHandler(podcastService AudiobooksPodcastService, updateService AudiobooksUpdateService,
	logger loggerrific.Logger, opts ...HandlerOption) *Handler {
	handler := &Handler{
		PodcastService: podcastService,
		UpdateService:  updateService,
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

	middleware := NewMiddlewares(h.Log)

	podcastSubrouter := h.PathPrefix(h.podcastServePath).Subrouter()
	podcastSubrouter.Use(middleware.LoggingMiddleware)
	podcastSubrouter.HandleFunc(h.mainFeedPath, h.getFeed)

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
