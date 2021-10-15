package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	mockTeam "github.com/batroff/todo-back/internal/team/mock"
	mockRel "github.com/batroff/todo-back/internal/team_relation_maker/mock"
	"github.com/batroff/todo-back/pkg/handler"
	testUtils "github.com/batroff/todo-back/pkg/testing"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_TeamCreateHandler(t *testing.T) {
	type createMockBehavior func(*mockTeam.MockUseCase, *models.Team)

	testID := testUtils.IdFixture()

	testTable := []struct {
		name               string
		inputBody          string
		inputTeam          *models.Team
		mock               createMockBehavior
		expectedLocationID string
		expectedStatus     int
		expectedBody       string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"team"}`,
			inputTeam: &models.Team{Name: "team"},
			mock: func(mock *mockTeam.MockUseCase, t *models.Team) {
				mock.EXPECT().CreateTeam(t).DoAndReturn(func(t *models.Team) error {
					t.ID = testID.Expected
					return nil
				})
			},
			expectedStatus:     http.StatusCreated,
			expectedLocationID: fmt.Sprintf("/teams/%s", testID.Expected),
			expectedBody:       "",
		},
		{
			name:               "Bad json",
			inputBody:          `{"name":,,}`,
			inputTeam:          &models.Team{},
			mock:               func(mock *mockTeam.MockUseCase, t *models.Team) {},
			expectedStatus:     http.StatusBadRequest,
			expectedLocationID: "",
			expectedBody:       fmt.Sprintf("%s: %s", models.ErrBadRequest, "invalid character ',' looking for beginning of value"),
		},
		{
			name:      "Create team internal error",
			inputBody: `{"name":"team"}`,
			inputTeam: &models.Team{Name: "team"},
			mock: func(mock *mockTeam.MockUseCase, t *models.Team) {
				mock.EXPECT().CreateTeam(t).Return(errors.New("internal"))
			},
			expectedStatus:     http.StatusInternalServerError,
			expectedLocationID: "",
			expectedBody:       "internal",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			mock := mockTeam.NewMockUseCase(c)
			testCase.mock(mock, testCase.inputTeam)

			// Request & response
			route := "/teams"
			req := httptest.NewRequest("POST", route, bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Server
			r := mux.NewRouter()
			h := teamCreateHandler(mock)
			r.Handle(route, h).Methods("POST").Headers("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
			assert.Equal(t, testCase.expectedLocationID, w.Header().Get("Location"))
		})
	}
}

func Test_TeamsListHandler(t *testing.T) {
	type listMockBehavior func(*mockTeam.MockUseCase)

	testTable := []struct {
		name           string
		inputQuery     map[string]interface{}
		mock           listMockBehavior
		expectedBody   string
		expectedStatus int
	}{
		{
			name:       "OK: Without filter",
			inputQuery: nil,
			mock: func(teamCase *mockTeam.MockUseCase) {
				teamCase.EXPECT().SelectTeamsList().Return([]*models.Team{
					{
						ID:   models.ID{},
						Name: "team #1",
					},
					{
						ID:   models.ID{},
						Name: "team #2",
					},
				}, nil)
			},
			expectedBody:   fmt.Sprintf(`[{"id":"%s","name":"team #1"},{"id":"%s","name":"team #2"}]%c`, models.ID{}, models.ID{}, '\n'),
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Empty result",
			inputQuery: nil,
			mock: func(teamCase *mockTeam.MockUseCase) {
				teamCase.EXPECT().SelectTeamsList().Return(nil, models.ErrNotFound)
			},
			expectedBody:   "[]\n",
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Internal error",
			inputQuery: nil,
			mock: func(teamCase *mockTeam.MockUseCase) {
				teamCase.EXPECT().SelectTeamsList().Return(nil, errors.New("internal"))
			},
			expectedBody:   "internal",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			mock := mockTeam.NewMockUseCase(c)
			testCase.mock(mock)

			// request & response
			route := "/teams"
			req := httptest.NewRequest("GET", handler.MakeQuery(route, testCase.inputQuery), nil)
			w := httptest.NewRecorder()

			// server
			h := teamListHandler(mock)
			r := mux.NewRouter()
			// Without query
			r.Handle(route, h).Methods("GET").Name("TeamsListHandler")
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func Test_TeamGetHandler(t *testing.T) {
	type getMockBehavior func(*mockTeam.MockUseCase, models.ID)

	queryID := testUtils.IdFixture()

	testTable := []struct {
		name           string
		queryID        testUtils.TestID
		mock           getMockBehavior
		expectedBody   string
		expectedStatus int
	}{
		{
			name:    "OK",
			queryID: queryID,
			mock: func(mock *mockTeam.MockUseCase, id models.ID) {
				mock.EXPECT().SelectTeamByID(id).Return(&models.Team{
					ID:   id,
					Name: "team #1",
				}, nil)
			},
			expectedBody:   fmt.Sprintf(`{"id":"%s","name":"team #1"}%c`, queryID.Expected, '\n'),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid id",
			queryID: testUtils.TestID{
				Input: "invalid",
			},
			mock:           func(mock *mockTeam.MockUseCase, id models.ID) {},
			expectedBody:   fmt.Sprintf("%s: %s", models.ErrBadRequest, "invalid UUID length: 7"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Team not found",
			queryID: queryID,
			mock: func(mock *mockTeam.MockUseCase, id models.ID) {
				mock.EXPECT().SelectTeamByID(id).Return(nil, models.ErrNotFound)
			},
			expectedBody:   models.ErrNotFound.Error(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "internal",
			queryID: queryID,
			mock: func(mock *mockTeam.MockUseCase, id models.ID) {
				mock.EXPECT().SelectTeamByID(id).Return(nil, errors.New("internal"))
			},
			expectedBody:   "internal",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			mock := mockTeam.NewMockUseCase(c)
			testCase.mock(mock, testCase.queryID.Expected)

			// Response & request
			route := "/teams"
			req := httptest.NewRequest("GET", handler.MakeURI(route, testCase.queryID.Input), nil)
			w := httptest.NewRecorder()

			// Server
			r := mux.NewRouter()
			h := teamGetHandler(mock)
			r.Handle(handler.MakeURI(route, "{id}"), h).Methods("GET")
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func Test_TeamDeleteHandler(t *testing.T) {
	type teamMockBehavior func(*mockTeam.MockUseCase, models.ID)
	type relMockBehavior func(*mockRel.MockUseCase, models.ID)

	queryID := testUtils.IdFixture()

	testTable := []struct {
		name           string
		queryID        testUtils.TestID
		teamMock       teamMockBehavior
		relMock        relMockBehavior
		expectedBody   string
		expectedStatus int
	}{
		{
			name:    "OK: with relation delete",
			queryID: queryID,
			teamMock: func(mock *mockTeam.MockUseCase, id models.ID) {
				mock.EXPECT().DeleteTeam(id).Return(nil)
			},
			relMock: func(mock *mockRel.MockUseCase, id models.ID) {
				mock.EXPECT().SelectRelationsByTeamID(id).Return([]*models.UserTeamRel{}, nil)
				mock.EXPECT().DeleteRelationsByTeamID(id).Return(nil)
			},
			expectedBody:   "",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:    "OK: without delete relations",
			queryID: queryID,
			teamMock: func(mock *mockTeam.MockUseCase, id models.ID) {
				mock.EXPECT().DeleteTeam(id).Return(nil)
			},
			relMock: func(mock *mockRel.MockUseCase, id models.ID) {
				mock.EXPECT().SelectRelationsByTeamID(id).Return(nil, models.ErrNotFound)
			},
			expectedBody:   "",
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "Bad request: Invalid ID",
			queryID: testUtils.TestID{
				Input: "invalid",
			},
			teamMock:       func(useCase *mockTeam.MockUseCase, id models.ID) {},
			relMock:        func(useCase *mockRel.MockUseCase, id models.ID) {},
			expectedBody:   fmt.Sprintf("%s: %s", models.ErrBadRequest, "invalid UUID length: 7"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Finding relations internal error",
			queryID:  queryID,
			teamMock: func(useCase *mockTeam.MockUseCase, id models.ID) {},
			relMock: func(useCase *mockRel.MockUseCase, id models.ID) {
				useCase.EXPECT().SelectRelationsByTeamID(id).Return(nil, errors.New("internal error"))
			},
			expectedBody:   "internal error",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Deleting relations internal error",
			queryID:  queryID,
			teamMock: func(useCase *mockTeam.MockUseCase, id models.ID) {},
			relMock: func(useCase *mockRel.MockUseCase, id models.ID) {
				useCase.EXPECT().SelectRelationsByTeamID(id).Return([]*models.UserTeamRel{}, nil)
				useCase.EXPECT().DeleteRelationsByTeamID(id).Return(errors.New("internal error"))
			},
			expectedBody:   "internal error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			relMock := mockRel.NewMockUseCase(c)
			testCase.relMock(relMock, queryID.Expected)

			teamMock := mockTeam.NewMockUseCase(c)
			testCase.teamMock(teamMock, queryID.Expected)

			// Response & request
			route := "/teams"
			req := httptest.NewRequest("GET", handler.MakeURI(route, testCase.queryID.Input), nil)
			w := httptest.NewRecorder()

			// Server
			r := mux.NewRouter()
			h := teamDeleteHandler(teamMock, relMock)
			r.Handle(handler.MakeURI(route, "{id}"), h).Methods("GET")
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}
