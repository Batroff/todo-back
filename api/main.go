package main

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/api/handler"
	"github.com/batroff/todo-back/internal/infrastructure/repository"
	"github.com/batroff/todo-back/internal/usecase/user"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/urfave/negroni"
	_ "github.com/urfave/negroni" // Add
	"log"
	"net/http"
	"time"
)

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "root", "localhost", 5432, "postgres")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error occured during db connection. err: %s", err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error occured during closing db connection. err: %s", err.Error())
		}
	}()

	userRepo := repository.NewUserMySQL(db)
	userService := user.NewService(userRepo)

	r := mux.NewRouter()
	//r.Schemes("http", "https").PathPrefix("/api")

	n := negroni.New()
	http.Handle("/", r)
	handler.MakeUserHandlers(r, *n, userService)

	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      context.ClearHandler(http.DefaultServeMux),
	}

	log.Fatal(srv.ListenAndServe())
}
