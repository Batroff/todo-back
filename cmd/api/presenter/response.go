package presenter

import (
	"encoding/json"
	"log"
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

func (w *ResponseWriter) Write(status int, body interface{}) {
	if v, ok := body.(error); ok {
		w.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.writer.WriteHeader(status)
		_, _ = w.writer.Write([]byte(v.Error()))
		return
	}

	if _, err := json.Marshal(body); err != nil {
		w.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = w.writer.Write([]byte(ErrNotImplementJsonMarshaller.Error()))
		log.Printf("error while marshalling response: %s", ErrNotImplementJsonMarshaller.Error())
		return
	}

	w.writer.WriteHeader(status)
	if err := json.NewEncoder(w.writer).Encode(&body); err != nil {
		log.Printf("error while writing response: %s", err)
	}
}

func (w *ResponseWriter) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		w.writer.Header().Set(k, v)
	}
}
