package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	mockTask "github.com/batroff/todo-back/internal/task/mock"
	mockTodo "github.com/batroff/todo-back/internal/todo/mock"
	"github.com/batroff/todo-back/pkg/handler"
	testUtils "github.com/batroff/todo-back/pkg/testing"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTodo_Create(t *testing.T) {
	type taskMockBehavior func(*mockTask.MockUseCase, models.ID)
	type todoMockBehavior func(*mockTodo.MockUseCase, *models.Todo)
	idFixture := testUtils.IdFixture()

	testTable := []struct {
		name               string
		inputBody          string
		inputTodo          *models.Todo
		taskMockBehavior   taskMockBehavior
		todoMockBehavior   todoMockBehavior
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: fmt.Sprintf(`{"title":"wash the dishes","text":"now","complete":false,"id_task":"%s"}`, idFixture.Input),
			inputTodo: &models.Todo{Title: func(s string) *string { return &s }("wash the dishes"), Text: "now", Complete: false, TaskID: idFixture.Expected},
			taskMockBehavior: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(&models.Task{}, nil)
			},
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, t *models.Todo) {
				todoCase.EXPECT().CreateTodo(t).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Bad request",
			inputBody:          `{"title":,}`,
			inputTodo:          &models.Todo{},
			expectedBody:       fmt.Sprintf("%s: invalid character ',' looking for beginning of value", models.ErrBadRequest),
			taskMockBehavior:   func(*mockTask.MockUseCase, models.ID) {},
			todoMockBehavior:   func(*mockTodo.MockUseCase, *models.Todo) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Task ID invalid",
			inputBody:          `{"id_task":"invalid-0000"}`,
			inputTodo:          &models.Todo{},
			expectedBody:       fmt.Sprintf("%s: invalid UUID length: 12", models.ErrBadRequest),
			taskMockBehavior:   func(*mockTask.MockUseCase, models.ID) {},
			todoMockBehavior:   func(*mockTodo.MockUseCase, *models.Todo) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:         "Task ID doesn't exist",
			inputBody:    fmt.Sprintf(`{"id_task":"%s"}`, idFixture.Input),
			inputTodo:    &models.Todo{TaskID: idFixture.Expected},
			expectedBody: fmt.Sprintf("task[%s] doesn't exist: %s", idFixture.Input, models.ErrNotFound),
			taskMockBehavior: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(nil, models.ErrNotFound)
			},
			todoMockBehavior:   func(*mockTodo.MockUseCase, *models.Todo) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Internal error",
			inputBody: fmt.Sprintf(`{"id_task":"%s"}`, idFixture.Input),
			inputTodo: &models.Todo{TaskID: idFixture.Expected},
			taskMockBehavior: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(&models.Task{}, nil)
			},
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, t *models.Todo) {
				todoCase.EXPECT().CreateTodo(t).Return(errors.New("some internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "some internal error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			// UseCase mocks
			taskCase := mockTask.NewMockUseCase(c)
			testCase.taskMockBehavior(taskCase, testCase.inputTodo.TaskID)

			todoCase := mockTodo.NewMockUseCase(c)
			testCase.todoMockBehavior(todoCase, testCase.inputTodo)

			// Router
			r := mux.NewRouter()
			route := handler.MakeURI("/api/v1", todoRoute)
			fmt.Println(route)
			h := todoCreateHandler(todoCase, taskCase)
			r.Handle(route, h).
				Methods("POST").
				Headers("Content-Type", "application/json")

			// Request & ResponseWriter
			req := httptest.NewRequest("POST", route, bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Server
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestTodo_Delete(t *testing.T) {
	testID := testUtils.IdFixture()
	type todoMockBehavior func(*mockTodo.MockUseCase, models.ID)

	testTable := []struct {
		name               string
		queryID            testUtils.TestID
		todoMockBehavior   todoMockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "OK",
			queryID: testID,
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().DeleteTodo(id).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Bad request[invalid ID]",
			queryID: testUtils.TestID{
				Input: "invalid-id",
			},
			todoMockBehavior:   func(*mockTodo.MockUseCase, models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("%s: invalid UUID length: %d", models.ErrBadRequest, 10),
		}, {
			name:    "Internal error: connection lost to db", // TODO : 500 error testing? OK?
			queryID: testID,
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().DeleteTodo(id).Return(errors.New("some internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "some internal error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			// Mock & router
			todoCase := mockTodo.NewMockUseCase(c)
			testCase.todoMockBehavior(todoCase, testCase.queryID.Expected)

			route := handler.MakeURI("/api/v1", todoRoute)
			r := mux.NewRouter()
			h := todoDeleteHandler(todoCase)
			r.Handle(handler.MakeURI(route, "{id}"), h).Methods("DELETE")

			// Request & response
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", handler.MakeURI(route, testCase.queryID.Input), nil)

			// Start server
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}
