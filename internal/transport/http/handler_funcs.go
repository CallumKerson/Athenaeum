package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"

	"github.com/go-chi/chi/v5"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/templates"
)

func healthCheck(writer http.ResponseWriter, request *http.Request) {
	SendJSON(writer, http.StatusOK, Payload{
		"health": "ok",
	})
}

func (h *Handler) readiness(writer http.ResponseWriter, request *http.Request) {
	if h.PodcastService.IsReady(request.Context()) {
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
	if h.version != "" {
		SendJSON(writer, http.StatusOK, Payload{
			"version": h.version,
		})
	}
}

func (h *Handler) getFeed(writer http.ResponseWriter, request *http.Request) {
	err := h.PodcastService.WriteAllAudiobooksFeed(request.Context(), writer)
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) getGenreFeed(writer http.ResponseWriter, request *http.Request) {
	genreStr, _ := url.PathUnescape(chi.URLParam(request, "genre"))
	genre, err := audiobooks.ParseGenre(genreStr)
	if err != nil {
		SendJSONError(writer, http.StatusNotFound, err)
		return
	}
	err = h.PodcastService.WriteGenreAudiobookFeed(request.Context(), genre, writer)
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) getAuthorFeed(writer http.ResponseWriter, request *http.Request) {
	h.getFeedForStr(writer, request, "author", h.PodcastService.WriteAuthorAudiobookFeed)
}

func (h *Handler) getNarratorFeed(writer http.ResponseWriter, request *http.Request) {
	h.getFeedForStr(writer, request, "narrator", h.PodcastService.WriteNarratorAudiobookFeed)
}

func (h *Handler) getTagFeed(writer http.ResponseWriter, request *http.Request) {
	h.getFeedForStr(writer, request, "tag", h.PodcastService.WriteTagAudiobookFeed)
}

func (h *Handler) getFeedForStr(writer http.ResponseWriter, request *http.Request, pathVar string,
	writeFunc func(context.Context, string, io.Writer) (bool, error)) {
	nameStr, err := url.PathUnescape(chi.URLParam(request, pathVar))
	if err != nil {
		SendJSONError(writer, http.StatusNotFound, err)
		return
	}
	written, err := writeFunc(request.Context(), nameStr, writer)
	if err != nil {
		SendJSONError(writer, http.StatusInternalServerError, err)
		return
	}
	if !written {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

func (h *Handler) updateAudiobooks(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		if h.CacheStore != nil {
			h.CacheStore.ReleaseAll()
		}
		if err := h.PodcastService.UpdateFeeds(request.Context()); err != nil {
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
	data := map[string]interface{}{
		"Title":           "Audiobooks",
		"Description":     "Like movies in your mind!",
		"StaticServePath": StaticPath,
		"FeedLink":        fmt.Sprintf("%s/%s", PodcastPath, PodcastFeedName),
	}
	if err := tpl.Execute(writer, data); err != nil {
		return
	}
}
