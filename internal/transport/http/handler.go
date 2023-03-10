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
	IsReady(ctx context.Context) (bool, error)
}

type Handler struct {
	*mux.Router
	Service        AudiobooksPodcastService
	Log            loggerrific.Logger
	mediaRoot      string
	mediaServePath string
}

func NewHandler(service AudiobooksPodcastService, logger loggerrific.Logger, opts ...HandlerOption) *Handler {
	handler := &Handler{
		Service: service,
		Log:     logger,
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

	middleware := NewMiddlewares(h.Log)

	podcastSubrouter := h.PathPrefix("/podcast").Subrouter()
	podcastSubrouter.Use(middleware.LoggingMiddleware)
	podcastSubrouter.HandleFunc("/feed.rss", h.getFeed)

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
	ready, err := h.Service.IsReady(request.Context())
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}

	if ready {
		SendJSON(writer, http.StatusOK, Payload{
			"readiness": "ok",
		})
	} else {
		SendJSON(writer, http.StatusInternalServerError, Payload{
			"readiness": "not ready",
		})
	}
}

func (h *Handler) getFeed(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add(ContentTypeHeader, ContentTypeXML)
	writer.WriteHeader(http.StatusOK)
	err := h.Service.WriteAllAudiobooksFeed(request.Context(), writer)
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
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
