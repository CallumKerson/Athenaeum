package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CallumKerson/loggerrific"
)

type PodcastService interface {
	GetFeed(ctx context.Context) (string, error)
	IsReady(ctx context.Context) (bool, error)
}

type Handler struct {
	*mux.Router
	Service PodcastService
	Log     loggerrific.Logger
}

func NewHandler(service PodcastService, logger loggerrific.Logger) *Handler {
	handler := &Handler{
		Service: service,
		Log:     logger,
	}
	handler.Router = mux.NewRouter()
	handler.mapRoutes()

	handler.Use(TimeoutMiddleware)
	m := NewMiddlewares(logger)
	handler.Use(m.LoggingMiddleware)

	return handler
}

func (h *Handler) mapRoutes() {
	h.HandleFunc("/health", healthCheck)
	h.HandleFunc("/ready", h.readiness)

	podcastSubrouter := h.PathPrefix("/podcast").Subrouter()
	m := NewMiddlewares(h.Log)
	podcastSubrouter.Use(m.LoggingMiddleware)
	podcastSubrouter.HandleFunc("/feed.rss", h.getFeed)
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
	feed, err := h.Service.GetFeed(request.Context())
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}
	writer.Header().Add(ContentTypeHeader, ContentTypeXML)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(feed))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
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
