package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CallumKerson/loggerrific"
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
	PodcastService AudiobooksPodcastService
	UpdateService  AudiobooksUpdateService
	Log            loggerrific.Logger
	version        string
	mediaRoot      string
	mediaServePath string
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

	podcastSubrouter := h.PathPrefix("/podcast").Subrouter()
	podcastSubrouter.Use(middleware.LoggingMiddleware)
	podcastSubrouter.HandleFunc("/feed.rss", h.getFeed)

	updateRouter := h.PathPrefix("/update").Subrouter()
	updateRouter.Use(SevereRateLimitMiddleware, middleware.LoggingMiddleware)
	updateRouter.HandleFunc("", h.updateAudiobooks)

	fs := http.StripPrefix(h.mediaServePath, http.FileServer(http.Dir(h.mediaRoot)))
	h.Log.Infoln("Serving files from local path", h.mediaRoot, "at", h.mediaServePath)
	h.Router.PathPrefix(h.mediaServePath).Handler(middleware.LoggingMiddleware(fs))
}

func healthCheck(writer http.ResponseWriter, request *http.Request) {
	SendJSON(writer, http.StatusOK, Payload{
		"health": "ok",
	})
}

func (h *Handler) readiness(writer http.ResponseWriter, request *http.Request) {
	if h.PodcastService.IsReady(request.Context()) && h.UpdateService.IsReady(request.Context()) {
		SendJSON(writer, http.StatusOK, Payload{
			"readiness": "ok",
		})
	} else {
		SendJSON(writer, http.StatusInternalServerError, Payload{
			"readiness": "not ready",
		})
	}
}

func (h *Handler) printVersion(writer http.ResponseWriter, request *http.Request) {
	SendJSON(writer, http.StatusOK, Payload{
		"version": h.version,
	})
}

func (h *Handler) getFeed(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add(ContentTypeHeader, ContentTypeXML)
	writer.WriteHeader(http.StatusOK)
	err := h.PodcastService.WriteAllAudiobooksFeed(request.Context(), writer)
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) updateAudiobooks(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		if err := h.UpdateService.UpdateAudiobooks(request.Context()); err != nil {
			SendJSONError(writer, http.StatusInternalServerError, err)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
		_, _ = writer.Write([]byte{})
	} else {
		writer.Header().Add("Allow", http.MethodPost)
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type Payload map[string]any

func SendJSONError(w http.ResponseWriter, status int, err error) {
	SendJSON(w, status, Payload{
		"status": status,
		"error":  err.Error(),
	})
}

func SendJSON(writer http.ResponseWriter, status int, p any) {
	writer.Header().Set(ContentTypeHeader, ContentTypeJSON)
	writer.WriteHeader(status)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(p)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
