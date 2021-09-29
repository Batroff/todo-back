// Code generated by MockGen. DO NOT EDIT.
// Source: usecase_interface.go

// Package mock_task is a generated GoMock package.
package mock_task

import (
	reflect "reflect"

	models "github.com/batroff/todo-back/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockUseCase) CreateTask(arg0 *models.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockUseCaseMockRecorder) CreateTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockUseCase)(nil).CreateTask), arg0)
}

// DeleteTaskByID mocks base method.
func (m *MockUseCase) DeleteTaskByID(arg0 models.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTaskByID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTaskByID indicates an expected call of DeleteTaskByID.
func (mr *MockUseCaseMockRecorder) DeleteTaskByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskByID", reflect.TypeOf((*MockUseCase)(nil).DeleteTaskByID), arg0)
}

// DeleteTaskByTeamID mocks base method.
func (m *MockUseCase) DeleteTaskByTeamID(arg0 models.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTaskByTeamID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTaskByTeamID indicates an expected call of DeleteTaskByTeamID.
func (mr *MockUseCaseMockRecorder) DeleteTaskByTeamID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskByTeamID", reflect.TypeOf((*MockUseCase)(nil).DeleteTaskByTeamID), arg0)
}

// DeleteTaskByUserID mocks base method.
func (m *MockUseCase) DeleteTaskByUserID(arg0 models.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTaskByUserID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTaskByUserID indicates an expected call of DeleteTaskByUserID.
func (mr *MockUseCaseMockRecorder) DeleteTaskByUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskByUserID", reflect.TypeOf((*MockUseCase)(nil).DeleteTaskByUserID), arg0)
}

// GetTaskByID mocks base method.
func (m *MockUseCase) GetTaskByID(arg0 models.ID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskByID", arg0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByID indicates an expected call of GetTaskByID.
func (mr *MockUseCaseMockRecorder) GetTaskByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByID", reflect.TypeOf((*MockUseCase)(nil).GetTaskByID), arg0)
}

// GetTasksBy mocks base method.
func (m *MockUseCase) GetTasksBy(arg0 map[string]interface{}) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksBy", arg0)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksBy indicates an expected call of GetTasksBy.
func (mr *MockUseCaseMockRecorder) GetTasksBy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksBy", reflect.TypeOf((*MockUseCase)(nil).GetTasksBy), arg0)
}

// GetTasksByTeamID mocks base method.
func (m *MockUseCase) GetTasksByTeamID(arg0 models.ID) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksByTeamID", arg0)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksByTeamID indicates an expected call of GetTasksByTeamID.
func (mr *MockUseCaseMockRecorder) GetTasksByTeamID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksByTeamID", reflect.TypeOf((*MockUseCase)(nil).GetTasksByTeamID), arg0)
}

// GetTasksByUserID mocks base method.
func (m *MockUseCase) GetTasksByUserID(arg0 models.ID) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksByUserID", arg0)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksByUserID indicates an expected call of GetTasksByUserID.
func (mr *MockUseCaseMockRecorder) GetTasksByUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksByUserID", reflect.TypeOf((*MockUseCase)(nil).GetTasksByUserID), arg0)
}

// GetTasksList mocks base method.
func (m *MockUseCase) GetTasksList() ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksList")
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksList indicates an expected call of GetTasksList.
func (mr *MockUseCaseMockRecorder) GetTasksList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksList", reflect.TypeOf((*MockUseCase)(nil).GetTasksList))
}

// UpdateTask mocks base method.
func (m *MockUseCase) UpdateTask(arg0 *models.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockUseCaseMockRecorder) UpdateTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockUseCase)(nil).UpdateTask), arg0)
}
