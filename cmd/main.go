package main

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/cmd/api/handler"
	"github.com/batroff/todo-back/cmd/api/middleware"

	taskRep "github.com/batroff/todo-back/internal/task/repository"
	taskUseCase "github.com/batroff/todo-back/internal/task/usecase"

	userRep "github.com/batroff/todo-back/internal/user/repository"
	userUseCase "github.com/batroff/todo-back/internal/user/usecase"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// TODO: put secret to .env.local
	if err := os.Setenv("secret", "secret"); err != nil {
		log.Fatalf("Error: err %s", err.Error())
	}

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

	userRepo := userRep.NewUserMySQL(db)
	userService := userUseCase.NewService(userRepo)
	taskRepo := taskRep.NewTaskPostgres(db)
	taskService := taskUseCase.NewService(taskRepo)

	r := mux.NewRouter()
	n := negroni.New(
		negroni.HandlerFunc(middleware.Auth),
	)

	apiV1 := r.PathPrefix("/api/v1/").Subrouter()

	handler.MakeAuthHandlers(apiV1, userService)
	handler.MakeUserHandlers(apiV1, *n, userService)
	handler.MakeTaskHandlers(apiV1, *n, taskService)

	http.Handle("/", r)

	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      context.ClearHandler(http.DefaultServeMux),
	}

	log.Fatal(srv.ListenAndServe())
}
