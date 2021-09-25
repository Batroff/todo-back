package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user"
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"
)

const usersRoute = "/users"

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
		rw.Header().Set("Location", handler.MakeRegexURI(usersRoute, id.String()))
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
		id, err := handler.GetIDFromURI(r)
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
		id, err := handler.GetIDFromURI(r)
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
		id, err := handler.GetIDFromURI(r)
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

			if refReq.Field(i).IsZero() && refUpd.FieldByName(f.Name).Kind() != reflect.Ptr {
				continue
			}

			var v reflect.Value

			switch refUpd.FieldByName(f.Name).Type().Kind() {
			case reflect.Ptr:
				v = refReq.Field(i)
			default:
				v = refReq.Field(i).Elem()
			}

			refUpd.FieldByName(f.Name).Set(v)
		}

		// Update user
		if err := useCase.UpdateUser(u); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// FIXME : do not return if image id is nil [bug]
		// Encode response
		if err := json.NewEncoder(rw).Encode(&requestUser); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
	})
}

func userOptionsHandler(rw http.ResponseWriter, _ *http.Request) {
	responseWriter := presenter.NewResponseWriter(rw)
	responseWriter.SetHeaders(map[string]string{
		"Allow": "GET, DELETE, PATCH, OPTIONS",
	})

	rw.WriteHeader(http.StatusOK)
}

// MakeUserHandlers sets up all user http handlers
func MakeUserHandlers(r *mux.Router, n negroni.Negroni, useCase user.UseCase) {
	// Create user
	r.Handle(usersRoute, n.With(
		negroni.Wrap(userCreateHandler(useCase)),
	)).Methods("POST").
		Headers("Content-Type", "application/json").
		Name("UserCreateHandler")

	// Get user list (opt: with query)
	r.Handle(usersRoute, n.With(
		negroni.Wrap(usersListHandler(useCase)),
	)).Methods("GET").
		Queries("email", "{email}").
		Name("UserQueryListHandler")

	r.Handle(usersRoute, n.With(
		negroni.Wrap(usersListHandler(useCase)),
	)).Methods("GET").
		Name("UserListHandler")
	// End user list

	// Get user by id
	r.Handle(handler.MakeRegexURI(usersRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(userGetByIDHandler(useCase)),
	)).Methods("GET").
		Name("UserGetByIDHandler")

	// Delete user by id
	r.Handle(handler.MakeRegexURI(usersRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(userDeleteHandler(useCase)),
	)).Methods("DELETE").
		Name("UserDeleteHandler")

	// Update user by id
	r.Handle(handler.MakeRegexURI(usersRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(userPatchHandler(useCase)),
	)).Methods("PATCH").
		Headers("Content-Type", "application/json").
		Name("UserPatchHandler")

	// users/:id OPTIONS handler
	r.Handle(handler.MakeRegexURI(usersRoute, handler.UUIDRegex), n.With(
		negroni.WrapFunc(userOptionsHandler),
	)).Methods("OPTIONS").
		Name("UserOptionsHandler")
}
