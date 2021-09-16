package main

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/api/handler"
	"github.com/batroff/todo-back/internal/infrastructure/repository"
	"github.com/batroff/todo-back/internal/usecase/user"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	_ = user.NewService(userRepo)
	// TODO : add service handler

	r := mux.NewRouter().
		Schemes("http", "https").
		PathPrefix("/api/")
	r.HandlerFunc(handler.UserCreateHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      r.GetHandler(),
	}

	log.Fatal(srv.ListenAndServe())
}
