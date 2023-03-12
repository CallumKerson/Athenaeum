package http

import (
	"fmt"
	"strings"
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
