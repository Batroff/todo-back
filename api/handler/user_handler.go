package handler

import (
	"github.com/batroff/todo-back/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func userCreateHandler(useCase user.UseCase) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {

		id, err := useCase.CreateUser("admin", "admin@localhost", "admin")
		if err != nil {
			log.Fatalf("userCreateHandler: err %s", err.Error())
			return
		}

		log.Printf("id = %s", id.String())
		return
	})
}

func MakeUserHandlers(r *mux.Router, n negroni.Negroni, useCase user.UseCase) {
	r.Handle("/create", n.With(
		negroni.Wrap(userCreateHandler(useCase)),
	)).Methods("GET") // TODO : add post method
}
