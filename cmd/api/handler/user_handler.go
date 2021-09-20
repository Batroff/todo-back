package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// TODO : Deprecated - remove
func writeResponseErr(rw http.ResponseWriter, statusCode int, err error) error {
	var res struct {
		Msg string `json:"msg"`
	}
	res.Msg = err.Error()

	rw.WriteHeader(statusCode)
	e := json.NewEncoder(rw).Encode(res)
	return e
}

func getUserID(r *http.Request) (models.ID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

func userCreateHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// Decoding request body
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("Error while decoding user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusBadRequest, err)
			return
		}

		// Creating user
		id, err := useCase.CreateUser(u.Login, u.Email, u.Password)
		if err != nil {
			log.Printf("Error while creating user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
		u.ID = id
		log.Printf("Created user: id = %s\n", id.String())

		// Encode response
		// FIXME : CreatedAt, ImageID shouldn't encode? Add GetUser(id)?
		rw.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(rw).Encode(&u)

		if err != nil {
			log.Printf("Error while encoding response. err %s", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func userListHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// SelectByID users list
		users, err := useCase.GetUsersList()
		if err != nil {
			log.Printf("Error while getting users. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Encode response
		err = json.NewEncoder(rw).Encode(&users)
		if err != nil {
			log.Printf("Error while encoding user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func userFindHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// Decode request
		var req struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error while decoding request. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusBadRequest, err)
			return
		}

		// Finding user in repo
		users, err := useCase.FindUsersBy(req.Key, req.Value)
		if err == models.ErrNotFound {
			log.Printf("Error user not found")
			_ = writeResponseErr(rw, http.StatusNotFound, err)
			return
		} else if err != nil {
			log.Printf("Error while finding user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Encode response
		err = json.NewEncoder(rw).Encode(&users)
		if err != nil {
			log.Printf("Error while encoding response. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func userGetHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// Decode request
		id, err := getUserID(r)
		if err != nil {
			log.Printf("Error while decoding request. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// SelectByID user from repo
		u, err := useCase.GetUser(id)
		if err == models.ErrNotFound {
			log.Printf("Error user not found")
			_ = writeResponseErr(rw, http.StatusNotFound, err)
			return
		} else if err != nil {
			log.Printf("Error while getting user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Encode & write response
		err = json.NewEncoder(rw).Encode(&u)
		if err != nil {
			log.Printf("Error while encoding response. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func userDeleteHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// SelectByID user id
		id, err := getUserID(r)
		if err != nil {
			log.Printf("Error while decoding request. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Deleting user
		if err := useCase.DeleteUser(id); err == models.ErrNotFound {
			log.Printf("Error while deleting user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusNotFound, err)
			return
		} else if err != nil {
			log.Printf("Error while deleting user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Encoding response
		var resp struct {
			ID string `json:"id"`
		}
		resp.ID = id.String()

		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			log.Printf("Error while encoding response. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func userUpdateHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ResponseWriter headers
		rw.Header().Set("Content-Type", "application/json")

		// SelectByID user id
		id, err := getUserID(r)
		if err != nil {
			log.Printf("Error while getting id. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Decode request
		var reqUser presenter.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			log.Printf("Error while decoding request. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusBadRequest, err)
			return
		}

		// Update user
		u := &models.User{
			ID:       id,
			Login:    reqUser.Login,
			Email:    reqUser.Email,
			Password: reqUser.Password,
			ImageID:  *reqUser.ImageID,
		}
		if err := useCase.UpdateUser(u); err != nil {
			log.Printf("Error while updating user. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}

		// Encode response
		if err := json.NewEncoder(rw).Encode(reqUser); err != nil {
			log.Printf("Error while encoding response. err %s\n", err.Error())
			_ = writeResponseErr(rw, http.StatusInternalServerError, err)
			return
		}
	})
}

func MakeUserHandlers(r *mux.Router, n negroni.Negroni, useCase user.UseCase) {
	r.Handle("/users", n.With(
		negroni.Wrap(userCreateHandler(useCase)),
	)).Methods("POST").Headers("Content-Type", "application/json")

	r.Handle("/users", n.With(
		negroni.Wrap(userListHandler(useCase)),
	)).Methods("GET")

	r.Handle("/users/find", n.With(
		negroni.Wrap(userFindHandler(useCase)),
	)).Methods("POST").Headers("Content-Type", "application/json")

	r.Handle("/users/{id}", n.With(
		negroni.Wrap(userGetHandler(useCase)),
	)).Methods("GET")

	r.Handle("/users/{id}", n.With(
		negroni.Wrap(userDeleteHandler(useCase)),
	)).Methods("DELETE")

	r.Handle("/users/{id}", n.With(
		negroni.Wrap(userUpdateHandler(useCase)),
	)).Methods("PUT").Headers("Content-Type", "application/json")
}
