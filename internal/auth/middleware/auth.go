package middleware

import (
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/batroff/todo-back/pkg/token"
	"net/http"
)

func Auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	responseWriter := handler.NewResponseWriter(rw)

	if ok, err := token.IsTokenValid(r); err != nil || !ok {
		responseWriter.Write(http.StatusUnauthorized, err)
		return
	}

	next(rw, r)
}
