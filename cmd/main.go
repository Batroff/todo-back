package main

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/configs"

	authHandler "github.com/batroff/todo-back/internal/auth/handler"
	authMiddleware "github.com/batroff/todo-back/internal/auth/middleware"

	taskHandler "github.com/batroff/todo-back/internal/task/handler"
	taskRep "github.com/batroff/todo-back/internal/task/repository"
	taskUseCase "github.com/batroff/todo-back/internal/task/usecase"

	todoHandler "github.com/batroff/todo-back/internal/todo/handler"
	todoRep "github.com/batroff/todo-back/internal/todo/repository"
	todoUseCase "github.com/batroff/todo-back/internal/todo/usecase"

	userHandler "github.com/batroff/todo-back/internal/user/handler"
	userRep "github.com/batroff/todo-back/internal/user/repository"
	userUseCase "github.com/batroff/todo-back/internal/user/usecase"

	teamHandler "github.com/batroff/todo-back/internal/team/handler"
	teamRep "github.com/batroff/todo-back/internal/team/repository"
	teamUseCase "github.com/batroff/todo-back/internal/team/usecase"

	relRep "github.com/batroff/todo-back/internal/team_relation_maker/repository"
	relUseCase "github.com/batroff/todo-back/internal/team_relation_maker/usecase"

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
	config := &configs.Config{}
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("loading config failed: %s", err.Error())
	}

	if err := os.Setenv("secret", config.Secret); err != nil {
		log.Fatalf("Error: err %s", err.Error())
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
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

	todoRepo := todoRep.NewTodoPostgres(db)
	todoService := todoUseCase.NewService(todoRepo)

	relRepo := relRep.NewTeamRelationMakerPostgres(db)
	relService := relUseCase.NewService(relRepo)

	teamRepo := teamRep.NewTeamPostgres(db)
	teamService := teamUseCase.NewService(teamRepo)

	r := mux.NewRouter()
	n := negroni.New(
		negroni.HandlerFunc(authMiddleware.Auth),
	)

	apiV1 := r.PathPrefix("/api/v1/").Subrouter()

	authHandler.MakeAuthHandlers(apiV1, userService)
	userHandler.MakeUserHandlers(apiV1, *n, userService)
	taskHandler.MakeTaskHandlers(apiV1, *n, taskService, userService)
	todoHandler.MakeTodoHandlers(apiV1, *n, todoService, taskService)
	teamHandler.MakeTeamHandlers(apiV1, *n, teamService, relService, userService)

	http.Handle("/", r)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.AppHost, config.AppPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      context.ClearHandler(http.DefaultServeMux),
	}

	log.Fatal(srv.ListenAndServe())
}
