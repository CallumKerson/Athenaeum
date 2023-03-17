package http

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/CallumKerson/Athenaeum/templates"
)

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

func (h *Handler) serveHTML(writer http.ResponseWriter, request *http.Request) {
	pages := map[string]string{
		"/": "index.gohtml",
	}

	page, ok := pages[request.URL.Path]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	tpl, err := template.ParseFS(templates.Templates, page)
	if err != nil {
		h.Log.WithError(err).Errorln("page", request.RequestURI, "not found in templates cache")
		h.Log.Debugln(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set(ContentTypeHeader, ContentTypeHTML)
	writer.WriteHeader(http.StatusOK)
	data := map[string]interface{}{
		"Title":           "Audiobooks",
		"Description":     "Like movies in your mind!",
		"StaticServePath": h.staticServePath,
		"FeedLink":        fmt.Sprintf("%s%s", h.podcastServePath, h.mainFeedPath),
	}
	if err := tpl.Execute(writer, data); err != nil {
		return
	}
}
