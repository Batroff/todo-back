package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batroff/todo-back/cmd/api/presenter"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/task"
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"
	"strings"
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

func getQueryFilterIDs(r *http.Request) (ids map[string]interface{}, err error) {
	ids = make(map[string]interface{}, 0)

	for k, v := range r.URL.Query() {
		if len(v) != 1 {
			return nil, fmt.Errorf("%s: expected key=value", presenter.ErrBadRequest)
		} else if !strings.Contains(k, "id") {
			continue
		}

		id, err := uuid.Parse(v[0])
		if err != nil {
			return nil, err
		}
		ids[k] = id
	}

	return ids, nil
}

func tasksListHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)
		responseWriter.SetHeaders(map[string]string{
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Content-Type":  "application/json; charset=utf-8",
			"Pragma":        "no-cache",
		})

		err := r.ParseForm()
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		// Query list
		filterIDs, err := getQueryFilterIDs(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		tasks, err := useCase.GetTasksBy(filterIDs)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusOK, make([]string, 0))
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		} else {

		}

		responseWriter.Write(http.StatusOK, tasks)
		return
	})
}

func taskUpdateHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)
		headers := map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Cache-Control": "no-store, no-cache, must-revalidate",
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
			responseWriter.Write(http.StatusNotFound, err)
			return
		}

		var reqTask presenter.RequestTask
		if err := json.NewDecoder(r.Body).Decode(&reqTask); err != nil {
			responseWriter.Write(http.StatusBadRequest, presenter.ErrBadRequest)
			return
		}

		refReq := reflect.ValueOf(reqTask)
		refUpd := reflect.ValueOf(t).Elem()
		handler.ParseRequestFields(refReq, refUpd)

		if err := useCase.UpdateTask(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Encode response
		if err := json.NewEncoder(rw).Encode(&reqTask); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
	})
}

func taskDeleteHandler(useCase task.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := presenter.NewResponseWriter(rw)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusNotFound, err)
			return
		}

		if err := useCase.DeleteTaskByID(id); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

func MakeTaskHandlers(r *mux.Router, n negroni.Negroni, useCase task.UseCase) {
	// Get tasks list (opt: with query)
	r.Handle(tasksRoute, n.With(
		negroni.Wrap(tasksListHandler(useCase)),
	)).Methods("GET").
		Name("TaskListHandler")

	r.Handle(tasksRoute, n.With(
		negroni.Wrap(tasksListHandler(useCase)),
	)).Methods("GET").
		Queries(
			"id_user", "{id_user}",
			"id_team", "{id_team}",
		).
		Name("TaskQueryListHandler")
	// End get tasks list

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

	// Update task
	r.Handle(handler.MakeRegexURI(tasksRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(taskUpdateHandler(useCase)),
	)).Methods("PATCH").
		Headers("Content-Type", "application/json").
		Name("TaskUpdateHandler")

	// Delete task
	r.Handle(handler.MakeRegexURI(tasksRoute, handler.UUIDRegex), n.With(
		negroni.Wrap(taskDeleteHandler(useCase)),
	)).Methods("DELETE").
		Name("TaskDeleteHandler")
}
