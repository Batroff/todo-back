package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/user"
	"github.com/batroff/todo-back/pkg/token"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func loginAuthHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		// Decode request
		var authReq presenter.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		// Find user in repo
		u, err := useCase.FindUserByEmail(authReq.Email)
		if err != nil {
			responseWriter.Write(http.StatusNotFound, err)
			return
		}

		// Compare userdata
		passwordErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(authReq.Password))
		if authReq.Email != u.Email || passwordErr != nil {
			responseWriter.Write(http.StatusUnauthorized, presenter.ErrUnauthorized)
			return
		}

		// Create token
		t, err := token.CreateToken(u.ID)
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Write response & cookie session_id
		cookie := http.Cookie{
			Name:     "session_id",
			Path:     "/api/",
			Value:    t,
			Secure:   true,
			HttpOnly: true,
		}
		responseWriter.SetCookie(&cookie)
		rw.WriteHeader(http.StatusOK)
	})
}

func MakeAuthHandlers(r *mux.Router, useCase user.UseCase) {
	r.Handle("/auth/login", loginAuthHandler(useCase))
}
