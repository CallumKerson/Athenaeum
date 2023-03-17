package http

import (
	"fmt"
	"strings"
)

const (
	defaultStaticPath = "/static"
)

type HandlerOption func(h *Handler)

func WithMediaConfig(mediaRoot, mediaServePath string) HandlerOption {
	return func(h *Handler) {
		h.mediaRoot = mediaRoot
		if !strings.HasSuffix(mediaServePath, "/") {
			mediaServePath = fmt.Sprintf("%s/", mediaServePath)
		}
		h.mediaServePath = mediaServePath
	}
}

func WithVersion(version string) HandlerOption {
	return func(h *Handler) {
		h.version = version
	}
}

func WithStaticPath(staticServePath string) HandlerOption {
	return func(handler *Handler) {
		if staticServePath == "" {
			staticServePath = defaultStaticPath
		}
		if !strings.HasSuffix(staticServePath, "/") {
			staticServePath = fmt.Sprintf("%s/", staticServePath)
		}
		if !strings.HasPrefix(staticServePath, "/") {
			staticServePath = fmt.Sprintf("/%s", staticServePath)
		}
		handler.staticServePath = staticServePath
	}
}
