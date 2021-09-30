package handler

import (
	"bytes"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	mockTask "github.com/batroff/todo-back/internal/task/mock"
	mockUser "github.com/batroff/todo-back/internal/user/mock"
	"github.com/batroff/todo-back/pkg/handler"
	testUtils "github.com/batroff/todo-back/pkg/testing"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTask_Get(t *testing.T) {
	type mockBehavior func(*mockTask.MockUseCase, models.ID)

	queryFixture := testUtils.IdFixture()
	testTable := []struct {
		name               string
		queryID            testUtils.TestID
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "OK",
			queryID: queryFixture,
			mockBehavior: func(useCase *mockTask.MockUseCase, id models.ID) {
				useCase.EXPECT().GetTaskByID(id).Return(&models.Task{ID: id, Title: "title", UserID: models.ID{}}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       fmt.Sprintf("{\"id\":\"%s\",\"title\":\"title\",\"id_user\":\"%s\"}\n", queryFixture.Input, models.ID{}),
		},
		{
			name: "Invalid ID length[must be 36]",
			queryID: testUtils.TestID{
				Input: "94-fd",
			},
			mockBehavior:       func(useCase *mockTask.MockUseCase, id models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("invalid UUID length: %d", 5),
		},
		{
			name: "Invalid ID",
			queryID: testUtils.TestID{
				Input: "6bc6f393-1e50-4647-9572-ce0de7b-a610",
			},
			mockBehavior:       func(useCase *mockTask.MockUseCase, id models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "invalid UUID format",
		},
		{
			name: "Not existing ID",
			queryID: testUtils.TestID{
				Input: models.ID{}.String(),
			},
			mockBehavior: func(useCase *mockTask.MockUseCase, id models.ID) {
				useCase.EXPECT().GetTaskByID(id).Return(nil, models.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       models.ErrNotFound.Error(),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			taskMock := mockTask.NewMockUseCase(c)
			testCase.mockBehavior(taskMock, testCase.queryID.Expected)

			// Test server
			getHandler := taskGetHandler(taskMock)
			r := mux.NewRouter()
			route := "/api/v1/tasks"
			r.Handle(handler.MakeURI(route, "{id}"), getHandler).Methods("GET")

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", handler.MakeURI(route, testCase.queryID.Input), bytes.NewBufferString(testCase.expectedBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestTask_Create(t *testing.T) {
	type taskMockBehavior func(*mockTask.MockUseCase, *models.Task)
	type userMockBehavior func(*mockUser.MockUseCase, models.ID)

	userFixtureID := testUtils.IdFixture()
	testTable := []struct {
		name               string
		inputBody          string
		inputTask          *models.Task
		taskMockBehavior   taskMockBehavior
		userMockBehavior   userMockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "OK",
			inputBody: fmt.Sprintf(`{"title":"task #1","id_user":"%s"}`, userFixtureID.Input),
			inputTask: &models.Task{
				Title:  "task #1",
				UserID: userFixtureID.Expected,
			},
			taskMockBehavior: func(taskCase *mockTask.MockUseCase, t *models.Task) {
				taskCase.EXPECT().CreateTask(t).Return(nil)
			},
			userMockBehavior: func(userCase *mockUser.MockUseCase, id models.ID) {
				userCase.EXPECT().GetUser(id).Return(&models.User{}, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:      "User ID doesn't exist",
			inputBody: fmt.Sprintf(`{"id_user":"%s"}`, userFixtureID.Input),
			inputTask: &models.Task{
				UserID: userFixtureID.Expected,
			},
			taskMockBehavior: func(taskCase *mockTask.MockUseCase, t *models.Task) {},
			userMockBehavior: func(userCase *mockUser.MockUseCase, id models.ID) {
				userCase.EXPECT().GetUser(id).Return(nil, models.ErrNotFound)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("user[%s] doesn't exist: %s", userFixtureID.Input, models.ErrNotFound),
		},
		{
			name:               "User ID invalid",
			inputBody:          `{"id_user":"00-fds-"}`,
			inputTask:          &models.Task{},
			taskMockBehavior:   func(taskCase *mockTask.MockUseCase, t *models.Task) {},
			userMockBehavior:   func(userCase *mockUser.MockUseCase, id models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("invalid UUID length: %d", 7),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			userMock := mockUser.NewMockUseCase(c)
			testCase.userMockBehavior(userMock, testCase.inputTask.UserID)

			taskMock := mockTask.NewMockUseCase(c)
			testCase.taskMockBehavior(taskMock, testCase.inputTask)

			// Handler
			h := taskCreateHandler(taskMock, userMock)
			r := mux.NewRouter()
			route := "/api/v1/tasks"
			r.Handle(route, h).Methods("POST").Headers("Content-Type", "application/json")

			// Http
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

func TestTask_Delete(t *testing.T) {
	type mockBehavior func(*mockTask.MockUseCase, models.ID)

	queryFixture := testUtils.IdFixture()
	testTable := []struct {
		name               string
		queryID            testUtils.TestID
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "OK",
			queryID: queryFixture,
			mockBehavior: func(taskMock *mockTask.MockUseCase, id models.ID) {
				taskMock.EXPECT().DeleteTaskByID(id).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Invalid ID length[must be 36]",
			queryID: testUtils.TestID{
				Input: "94-fd",
			},
			mockBehavior:       func(taskMock *mockTask.MockUseCase, id models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       fmt.Sprintf("invalid UUID length: %d", 5),
		},
		{
			name: "Invalid ID",
			queryID: testUtils.TestID{
				Input: "6bc6f393-1e50-4647-9572-ce0de7b-a610",
			},
			mockBehavior:       func(taskMock *mockTask.MockUseCase, id models.ID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "invalid UUID format",
		},
		{
			name: "Not existing ID",
			queryID: testUtils.TestID{
				Input: models.ID{}.String(),
			},
			mockBehavior: func(taskMock *mockTask.MockUseCase, id models.ID) {
				taskMock.EXPECT().DeleteTaskByID(id).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			taskMock := mockTask.NewMockUseCase(c)
			testCase.mockBehavior(taskMock, testCase.queryID.Expected)

			// Handler
			r := mux.NewRouter()
			h := taskDeleteHandler(taskMock)
			route := "/api/v1/tasks"
			r.Handle(handler.MakeURI(route, "{id}"), h).Methods("DELETE")

			// Request & response
			req := httptest.NewRequest("DELETE", handler.MakeURI(route, testCase.queryID.Input), nil)
			w := httptest.NewRecorder()

			// Server
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

// TODO : implement
func TestTask_Update(t *testing.T) {

}

func TestTask_List(t *testing.T) {

}
