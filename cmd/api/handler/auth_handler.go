package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/cmd/api/middleware"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/user"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func loginAuthHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Decode request
		var authReq presenter.AuthRequest
		var authRes presenter.AuthResponse
		responseWriter := presenter.NewResponseWriter(rw)

		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			authRes.Msg = err.Error()
			log.Println(responseWriter.Write(http.StatusBadRequest, authRes))
			return
		}
		authRes.Request = authReq

		// Find user in repo
		u, err := useCase.FindUserByEmail(authReq.Email)
		if err != nil {
			authRes.Msg = err.Error()
			log.Println(responseWriter.Write(http.StatusNotFound, authRes))
			return
		}

		// Compare userdata
		err = bcrypt.CompareHashAndPassword([]byte(authReq.Password), []byte(u.Password))

		if authReq.Email != u.Email || err != nil {
			authRes.Msg = err.Error()
			log.Println(responseWriter.Write(http.StatusUnauthorized, authRes))
			return
		}

		// Create token
		if authRes.Token, err = middleware.CreateToken(u.ID); err != nil {
			authRes.Msg = err.Error()
			log.Println(responseWriter.Write(http.StatusInternalServerError, authRes))
			return
		}

		// Write response & cookie session_id
		cookie := http.Cookie{Name: "session_id", Path: "/", Value: authRes.Token, Secure: true, HttpOnly: true, Expires: time.Now().Add(time.Hour * 24)}
		http.SetCookie(rw, &cookie)
		if err = responseWriter.Write(http.StatusOK, authRes); err != nil {
			log.Println(err.Error())
			return
		}
	})
}

func MakeAuthHandlers(r *mux.Router, useCase user.UseCase) {
	r.Handle("/auth/login", loginAuthHandler(useCase))
}
