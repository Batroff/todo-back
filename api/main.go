package main

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/internal/infrastructure/repository"
	"github.com/batroff/todo-back/internal/usecase/user"
	"log"
)
import _ "github.com/lib/pq"

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "root", "localhost", 5432, "todo")
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
}
