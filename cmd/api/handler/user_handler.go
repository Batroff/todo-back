package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"
)

const entityPrefix = "/users"
const userIdRegex = "{id:\\w{8}-(?:\\w{4}-){3}\\w{12}}"

func makeRegexURI(prefix, regex string) string {
	return fmt.Sprintf("%s/%s", prefix, regex)
}

func getUserID(r *http.Request) (models.ID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

func userCreateHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		// Decoding request body
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		// Creating user
		id, err := useCase.CreateUser(u.Login, u.Email, u.Password)
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
		u.ID = id

		// Encode response
		rw.Header().Set("Location", makeRegexURI(entityPrefix, id.String()))
		rw.WriteHeader(http.StatusCreated)
	})
}

func usersListHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		if err := r.ParseForm(); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
		}

		if email, ok := r.Form["email"]; ok || len(email) == 1 {
			u, err := useCase.FindUserByEmail(email[0])
			if err == models.ErrNotFound {
				responseWriter.Write(http.StatusOK, make([]string, 0))
				return
			} else if err != nil {
				responseWriter.Write(http.StatusInternalServerError, err)
				return
			}

			responseWriter.Write(http.StatusOK, u)
			return
		}

		// SelectByID users list
		users, err := useCase.GetUsersList()
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Encode response
		if users == nil {
			responseWriter.Write(http.StatusOK, make([]string, 0))
			return
		}
		responseWriter.Write(http.StatusOK, users)
	})
}

func userGetByIDHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		// Decode request
		id, err := getUserID(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, err)
			return
		}

		// SelectByID user from repo
		u, err := useCase.GetUser(id)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.Write(http.StatusOK, u)
	})
}

func userDeleteHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		// SelectByID user id
		id, err := getUserID(r)
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Deleting user
		if err := useCase.DeleteUser(id); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

func userPatchHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		responseWriter := presenter.NewResponseWriter(rw)

		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		// SelectByID user id
		id, err := getUserID(r)
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Try to find user
		u, err := useCase.GetUser(id)
		if err != nil {
			responseWriter.Write(http.StatusNotFound, err)
			return
		}

		// Decode request
		var requestUser presenter.RequestUser
		if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		// Update only requested fields
		refReq := reflect.ValueOf(requestUser)
		refUpd := reflect.ValueOf(u).Elem()

		t := refReq.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if !refReq.Field(i).IsZero() {
				v := refReq.Field(i).Elem()
				refUpd.FieldByName(f.Name).Set(v)
			}
		}

		// Update user
		if err := useCase.UpdateUser(u); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Encode response
		if err := json.NewEncoder(rw).Encode(&requestUser); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
	})
}

func MakeUserHandlers(r *mux.Router, n negroni.Negroni, useCase user.UseCase) {
	// TODO : add HEAD, OPTIONS methods for /users/:id endpoint

	// Create user
	r.Handle(entityPrefix, n.With(
		negroni.Wrap(userCreateHandler(useCase)),
	)).Methods("POST").
		Headers("Content-Type", "application/json").
		Name("UserCreateHandler")

	// Get user list (opt: with query)
	r.Handle(entityPrefix, n.With(
		negroni.Wrap(usersListHandler(useCase)),
	)).Methods("GET").
		Queries("email", "{email}").
		Name("UserQueryListHandler")

	r.Handle(entityPrefix, n.With(
		negroni.Wrap(usersListHandler(useCase)),
	)).Methods("GET").
		Name("UserListHandler")
	// End user list

	// Get user by id
	r.Handle(makeRegexURI(entityPrefix, userIdRegex), n.With(
		negroni.Wrap(userGetByIDHandler(useCase)),
	)).Methods("GET").
		Name("UserGetByIDHandler")

	// Delete user by id
	r.Handle(makeRegexURI(entityPrefix, userIdRegex), n.With(
		negroni.Wrap(userDeleteHandler(useCase)),
	)).Methods("DELETE").
		Name("UserDeleteHandler")

	// Update user by id
	r.Handle(makeRegexURI(entityPrefix, userIdRegex), n.With(
		negroni.Wrap(userPatchHandler(useCase)),
	)).Methods("PATCH").
		Headers("Content-Type", "application/json").
		Name("UserPatchHandler")
}
