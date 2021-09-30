package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/task"
	"github.com/batroff/todo-back/internal/todo"
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

const todoRoute = "/todos"

func todoCreateHandler(todoCase todo.UseCase, taskCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		var t *models.Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		if _, err := taskCase.GetTaskByID(t.TaskID); err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("task[%s] doesn't exist: %s", t.TaskID, err))
			return
		}

		if err := todoCase.CreateTodo(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.SetHeaders(map[string]string{
			"Location": handler.MakeURI(r.RequestURI, t.ID.String()),
		})
		rw.WriteHeader(http.StatusCreated)
	})
}

func todoDeleteHandler(todoCase todo.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		if err := todoCase.DeleteTodo(id); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

func MakeTodoHandlers(r *mux.Router, n negroni.Negroni, todoCase todo.UseCase, taskCase task.UseCase) {
	// Create
	r.Handle(todoRoute, n.With(
		negroni.Wrap(todoCreateHandler(todoCase, taskCase)),
	)).Methods("POST").
		Headers("Content-Type", "application/json").
		Name("CreateTodoHandler")

	// Delete
	r.Handle(handler.MakeURI(todoRoute, "{id}"), n.With(
		negroni.Wrap(todoDeleteHandler(todoCase)),
	)).Methods("DELETE").
		Name("DeleteTodoHandler")
}
