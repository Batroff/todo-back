package handler

import (
	"encoding/json"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/task"
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

const tasksRoute = "/tasks"

func taskGetHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		headers := map[string]string{
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Content-Type":  "application/json; charset=utf-8",
			"Pragma":        "no-cache",
		}
		responseWriter.SetHeaders(headers)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusNotFound, err)
			return
		}

		t, err := useCase.GetTaskByID(id)
		if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.Write(http.StatusOK, t)
	})
}

func taskCreateHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		var t *models.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			responseWriter.Write(http.StatusBadRequest, err)
			return
		}

		t = models.NewTask(t.Title, t.Priority, t.UserID, t.TeamID)
		if err := useCase.CreateTask(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.SetHeaders(map[string]string{
			"Location": handler.MakeRegexURI(tasksRoute, t.ID.String()),
		})
		rw.WriteHeader(http.StatusCreated)
	})
}

func tasksListHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)
		responseWriter.SetHeaders(map[string]string{
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Content-Type":  "application/json; charset=utf-8",
			"Pragma":        "no-cache",
		})

		tasks, err := useCase.GetTasksList()
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusOK, make([]string, 0))
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.Write(http.StatusOK, tasks)
	})
}

func MakeTaskHandlers(r *mux.Router, n negroni.Negroni, useCase task.UseCase) {
	// Get tasks list
	r.Handle(tasksRoute, n.With(
		negroni.Wrap(tasksListHandler(useCase)),
	)).Methods("GET").
		Name("TaskListHandler")

	// Get task by ID
	r.Handle(handler.MakeRegexURI(tasksRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(taskGetHandler(useCase)),
	)).Methods("GET").
		Name("TaskGetHandler")

	// Create task
	r.Handle(tasksRoute, n.With(
		negroni.Wrap(taskCreateHandler(useCase)),
	)).Methods("POST").
		Headers("Content-Type", "application/json").
		Name("TaskCreateHandler")
}
