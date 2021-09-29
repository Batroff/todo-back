package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/internal/models"
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
		if _, err := w.writer.Write([]byte(v.Error())); err != nil {
			log.Printf("unexpected error in ResponseWriter: %s", err)
		}
		return
	}

	if _, err := json.Marshal(body); err != nil {
		w.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.writer.WriteHeader(http.StatusInternalServerError)
		if _, err := w.writer.Write([]byte(models.ErrNotImplementJsonMarshaller.Error())); err != nil {
			log.Printf("unexpected error in ResponseWriter: %s", err)
		}
		log.Printf("error while marshalling response: %s", models.ErrNotImplementJsonMarshaller.Error())
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

func (w *ResponseWriter) SetCookie(cookie *http.Cookie) {
	http.SetCookie(w.writer, cookie)
}
