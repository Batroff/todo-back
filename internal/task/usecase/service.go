// Package usecase represents 2nd layer of application.
// It provides Service with methods to operate repository
package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/task"
)

type Service struct {
	rep task.Repository
}

// NewService returns *Service to operate task.Repository
func NewService(repository task.Repository) *Service {
	return &Service{
		rep: repository,
	}
}

// GetTasksList returns slice of *models.Task
func (s *Service) GetTasksList() ([]*models.Task, error) {
	return s.rep.SelectAll()
}

// GetTaskByID takes task id models.ID and returns *models.Task
func (s *Service) GetTaskByID(id models.ID) (*models.Task, error) {
	return s.rep.SelectByID(id)
}

// GetTasksByUserID takes user's id models.ID, who created task and returns slice of *models.Task
func (s *Service) GetTasksByUserID(id models.ID) ([]*models.Task, error) {
	return s.rep.SelectByUserID(id)
}

// GetTasksByTeamID takes team's id models.ID, which task relates to and returns slice of *models.Task
func (s *Service) GetTasksByTeamID(id models.ID) ([]*models.Task, error) {
	return s.rep.SelectByTeamID(id)
}

// GetTasksBy takes map with key = param name and value = param value and returns slice of *models.Task
func (s *Service) GetTasksBy(pairs map[string]interface{}) ([]*models.Task, error) {
	return s.rep.SelectBy(pairs)
}

// CreateTask takes *models.Task and inserts it into repository
func (s *Service) CreateTask(t *models.Task) error {
	return s.rep.Insert(t)
}

// UpdateTask takes *models.Task and update it in repository
func (s *Service) UpdateTask(t *models.Task) error {
	return s.rep.Update(t)
}

// DeleteTaskByID takes task models.ID and deletes it from repository
func (s *Service) DeleteTaskByID(id models.ID) error {
	return s.rep.DeleteByID(id)
}

// DeleteTaskByUserID takes user's models.ID, who created task and deletes it from repository
func (s *Service) DeleteTaskByUserID(id models.ID) error {
	return s.rep.DeleteByUserID(id)
}

// DeleteTaskByTeamID takes team's models.ID, which relates to and deletes it from repository
func (s *Service) DeleteTaskByTeamID(id models.ID) error {
	return s.rep.DeleteByTeamID(id)
}
