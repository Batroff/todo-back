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

func TestTodo_List(t *testing.T) {
	type todoMockBehavior func(useCase *mockTodo.MockUseCase)

	testTable := []struct {
		name             string
		todoMockBehavior todoMockBehavior
		expectedBody     string
		expectedStatus   int
	}{
		{
			name: "OK",
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase) {
				todoCase.EXPECT().GetTodosList().Return([]*models.Todo{
					{Text: "Some text"},
				}, nil)
			},
			expectedBody:   fmt.Sprintf(`[{"id":"%s","text":"Some text","id_task":"%s"}]%s`, models.ID{}, models.ID{}, "\n"),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty result - should return empty slice",
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase) {
				todoCase.EXPECT().GetTodosList().Return(nil, models.ErrNotFound)
			},
			expectedBody:   "[]\n",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Internal error",
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase) {
				todoCase.EXPECT().GetTodosList().Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			// Mock & router
			todoMock := mockTodo.NewMockUseCase(c)
			testCase.todoMockBehavior(todoMock)

			r := mux.NewRouter()
			route := "/api/v1/todos"
			h := todoListHandler(todoMock)
			r.Handle(route, h).Methods("GET")

			// Request & response
			req := httptest.NewRequest("GET", route, nil)
			w := httptest.NewRecorder()

			// Server
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestTodo_Get(t *testing.T) {
	type todoMockBehavior func(*mockTodo.MockUseCase, models.ID)

	fixtureID := testUtils.IdFixture()
	testTable := []struct {
		name             string
		queryID          testUtils.TestID
		todoMockBehavior todoMockBehavior
		expectedBody     string
		expectedStatus   int
	}{
		{
			name:    "OK",
			queryID: fixtureID,
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(&models.Todo{ID: id}, nil)
			},
			expectedBody:   fmt.Sprintf(`{"id":"%s","id_task":"%s"}%s`, fixtureID.Input, models.ID{}, "\n"),
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Not found todo",
			queryID: fixtureID,
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(nil, models.ErrNotFound)
			},
			expectedBody:   models.ErrNotFound.Error(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Internal error",
			queryID: fixtureID,
			todoMockBehavior: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			// Mock & router
			todoMock := mockTodo.NewMockUseCase(c)
			testCase.todoMockBehavior(todoMock, testCase.queryID.Expected)

			r := mux.NewRouter()
			h := todoGetHandler(todoMock)
			route := "/api/v1/todos"
			r.Handle(handler.MakeURI(route, "{id}"), h)

			// Request & response
			req := httptest.NewRequest("GET", handler.MakeURI(route, testCase.queryID.Input), nil)
			w := httptest.NewRecorder()

			// Server
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestTodo_Update(t *testing.T) {
	type todoUpdateMockBehavior func(*mockTodo.MockUseCase, *models.Todo)
	type todoGetMockBehavior func(*mockTodo.MockUseCase, models.ID)
	type taskGetMockBehavior func(*mockTask.MockUseCase, models.ID)

	fixtureID := testUtils.IdFixture()

	testTable := []struct {
		name               string
		inputBody          string
		inputTodo          *models.Todo
		queryID            testUtils.TestID
		todoUpdateMock     todoUpdateMockBehavior
		todoGetMock        todoGetMockBehavior
		taskGetMock        taskGetMockBehavior
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"title":"Title #1","text":"New Body"}`,
			inputTodo: &models.Todo{
				ID:    fixtureID.Expected,
				Title: func(s string) *string { return &s }("Title #1"),
				Text:  "New Body",
			},
			queryID: fixtureID,
			taskGetMock: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(&models.Task{}, nil)
			},
			todoUpdateMock: func(todoCase *mockTodo.MockUseCase, t *models.Todo) {
				todoCase.EXPECT().UpdateTodo(t).Return(nil)
			},
			todoGetMock: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(&models.Todo{
					ID:       fixtureID.Expected,
					Title:    nil,
					Text:     "Old string",
					Complete: false,
					TaskID:   models.ID{},
				}, nil)
			},
			expectedBody:       fmt.Sprintln(`{"title":"Title #1","text":"New Body"}`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Bad request[invalid ID]",
			queryID: testUtils.TestID{
				Input: "invalid-id",
			},
			inputBody:          `{}`,
			inputTodo:          &models.Todo{},
			todoGetMock:        func(*mockTodo.MockUseCase, models.ID) {},
			todoUpdateMock:     func(*mockTodo.MockUseCase, *models.Todo) {},
			taskGetMock:        func(*mockTask.MockUseCase, models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("%s: invalid UUID length: %d", models.ErrBadRequest, 10),
		},
		{
			name:               "Bad request[invalid JSON]",
			queryID:            fixtureID,
			inputBody:          `{"Title": true}`,
			inputTodo:          &models.Todo{},
			todoGetMock:        func(*mockTodo.MockUseCase, models.ID) {},
			todoUpdateMock:     func(*mockTodo.MockUseCase, *models.Todo) {},
			taskGetMock:        func(*mockTask.MockUseCase, models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("%s: json: cannot unmarshal bool into Go struct field RequestTodo.title of type string", models.ErrBadRequest),
		},
		{
			name:      "Todo doesn't exist",
			queryID:   fixtureID,
			inputBody: `{}`,
			inputTodo: &models.Todo{},
			todoGetMock: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(nil, models.ErrNotFound)
			},
			todoUpdateMock:     func(*mockTodo.MockUseCase, *models.Todo) {},
			taskGetMock:        func(*mockTask.MockUseCase, models.ID) {},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       models.ErrNotFound.Error(),
		},
		{
			name:      "Todo get internal error",
			queryID:   fixtureID,
			inputBody: `{}`,
			inputTodo: &models.Todo{},
			todoGetMock: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(nil, errors.New("internal error"))
			},
			todoUpdateMock:     func(*mockTodo.MockUseCase, *models.Todo) {},
			taskGetMock:        func(*mockTask.MockUseCase, models.ID) {},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "internal error",
		},
		{
			name:      "Task ID in request doesn't exist",
			queryID:   fixtureID,
			inputBody: fmt.Sprintf(`{"id_task":"%s"}`, models.ID{}),
			inputTodo: &models.Todo{},
			todoGetMock: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(&models.Todo{}, nil)
			},
			taskGetMock: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(nil, models.ErrNotFound)
			},
			todoUpdateMock:     func(*mockTodo.MockUseCase, *models.Todo) {},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       models.ErrNotFound.Error(),
		},
		{
			name:      "Update todo internal error",
			queryID:   fixtureID,
			inputBody: fmt.Sprintf(`{"id_task":"%s"}`, models.ID{}),
			inputTodo: &models.Todo{ID: fixtureID.Expected},
			todoGetMock: func(todoCase *mockTodo.MockUseCase, id models.ID) {
				todoCase.EXPECT().GetTodoByID(id).Return(&models.Todo{ID: id}, nil)
			},
			taskGetMock: func(taskCase *mockTask.MockUseCase, id models.ID) {
				taskCase.EXPECT().GetTaskByID(id).Return(&models.Task{}, nil)
			},
			todoUpdateMock: func(todoCase *mockTodo.MockUseCase, t *models.Todo) {
				todoCase.EXPECT().UpdateTodo(t).Return(errors.New("internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "internal error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			// Router & mock
			todoMock := mockTodo.NewMockUseCase(c)
			taskMock := mockTask.NewMockUseCase(c)

			testCase.todoGetMock(todoMock, testCase.queryID.Expected)

			testCase.taskGetMock(taskMock, testCase.inputTodo.TaskID)

			testCase.todoUpdateMock(todoMock, testCase.inputTodo)

			r := mux.NewRouter()
			h := todoUpdateHandler(todoMock, taskMock)
			route := "/api/v1/todos"
			r.Handle(handler.MakeURI(route, "{id}"), h).Methods("PATCH").Headers("Content-Type", "application/json")

			// Request & response
			req := httptest.NewRequest("PATCH", handler.MakeURI(route, testCase.queryID.Input), bytes.NewBufferString(testCase.inputBody))
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
