package http

import (
	"encoding/json"
	"net/http"
)

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
