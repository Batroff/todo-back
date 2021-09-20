package presenter

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ResponseWriter struct {
	writer http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		writer: w,
	}
}

func (w *ResponseWriter) Write(status int, body interface{}) error {
	if _, err := json.Marshal(body); err != nil {
		w.writer.WriteHeader(http.StatusBadRequest)
		return errors.New("err: body doesn't implement json.Marshaller")
	}

	w.writer.WriteHeader(status)
	return json.NewEncoder(w.writer).Encode(&body)
}
