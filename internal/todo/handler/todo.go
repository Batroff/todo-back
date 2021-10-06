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
	"reflect"
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

func todoListHandler(todoCase todo.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)
		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		todos, err := todoCase.GetTodosList()
		if err == models.ErrNotFound {
			todos = make([]*models.Todo, 0)
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.Write(http.StatusOK, todos)
	})
}

func todoGetHandler(todoCase todo.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)
		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		t, err := todoCase.GetTodoByID(id)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// TODO: Error is possible? How to test
		if err := json.NewEncoder(rw).Encode(&t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})
}

func todoUpdateHandler(todoCase todo.UseCase, taskCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)
		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		// Parse request
		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		var req models.RequestTodo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		// Find existing todo_obj
		t, err := todoCase.GetTodoByID(id)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Parse request todo_obj
		refReq := reflect.ValueOf(req)
		refUpd := reflect.ValueOf(t).Elem()
		handler.ParseRequestFields(refReq, refUpd)

		// Check task ID
		if _, err := taskCase.GetTaskByID(t.TaskID); err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusNotFound, fmt.Sprintf("%s: %s", models.ErrNotFound, err))
			return
		}

		// Update
		if err := todoCase.UpdateTodo(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Is error possible? TODO : Remove error if?
		if err := json.NewEncoder(rw).Encode(&req); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
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

	// List
	r.Handle(todoRoute, n.With(
		negroni.Wrap(todoListHandler(todoCase)),
	)).Methods("GET").
		Name("TodoListHandler")

	// Get
	r.Handle(handler.MakeURI(todoRoute, "{id}"), n.With(
		negroni.Wrap(todoGetHandler(todoCase)),
	)).Methods("GET").
		Name("TodoGetHandler")

	// Update
	r.Handle(handler.MakeURI(todoRoute, "{id}"), n.With(
		negroni.Wrap(todoUpdateHandler(todoCase, taskCase)),
	)).Methods("PATCH").
		Headers("Content-Type", "application/json").
		Name("TodoUpdateHandler")
}
