package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user"
	"github.com/batroff/todo-back/pkg/token"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

const emailRegex = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

// isPasswordValid checks for:
// 1) password length > 6;
// 2) password uses both letter cases;
// 3) password has numbers in it.
func isPasswordValid(pass string) error {
	if len(pass) <= 6 {
		return fmt.Errorf("password too short. minimum length is 7")
	}

	alphavite := "abcdefghijklmnoprstquvwxyz"
	if !strings.ContainsAny(pass, alphavite) || !strings.ContainsAny(pass, strings.ToUpper(alphavite)) {
		return fmt.Errorf("password must contain both letter cases")
	}

	numbers := "0123456789"
	if !strings.ContainsAny(pass, numbers) {
		return fmt.Errorf("password should contain atleast 1 number")
	}

	return nil
}

// isEmailValid checks for email by regex
func isEmailValid(email string) error {
	if matched, err := regexp.Match(emailRegex, []byte(email)); err != nil {
		return err
	} else if !matched {
		return fmt.Errorf("email value is incorrect")
	}

	return nil
}

func registerAuthHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		// Check for existing
		if u, err := useCase.FindUserByEmail(u.Email); u != nil || err == nil {
			responseWriter.Write(http.StatusConflict, fmt.Errorf("user with email %s already exists", u.Email))
			return
		}

		// Check email & password
		if err := isEmailValid(u.Email); err != nil {
			responseWriter.Write(http.StatusBadRequest, err)
			return
		} else if err = isPasswordValid(u.Password); err != nil {
			responseWriter.Write(http.StatusBadRequest, err)
			return
		}

		// Creating user
		if _, err := useCase.CreateUser(u.Login, u.Email, u.Password); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.SetHeaders(map[string]string{
			"Location": "/api/v1/auth/login",
		})
		rw.WriteHeader(http.StatusCreated)
	})
}

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
			Path:     "/api/v1",
			Value:    t,
			Secure:   true,
			HttpOnly: true,
		}
		responseWriter.SetCookie(&cookie)
		rw.WriteHeader(http.StatusOK)
	})
}

func logoutAuthHandler(rw http.ResponseWriter, r *http.Request) {
	responseWriter := presenter.NewResponseWriter(rw)

	_, err := r.Cookie("session_id")
	if err != nil {
		responseWriter.Write(http.StatusUnauthorized, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Path:     "/api/v1",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	}
	responseWriter.SetCookie(cookie)
	rw.WriteHeader(http.StatusOK)
}

func MakeAuthHandlers(r *mux.Router, useCase user.UseCase) {
	// Login
	r.Handle("/auth/login", loginAuthHandler(useCase)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	// Logout
	r.HandleFunc("/auth/logout", logoutAuthHandler).
		Methods("GET")

	// Register
	r.Handle("/auth/register", registerAuthHandler(useCase)).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Name("UserCreateHandler")
}
