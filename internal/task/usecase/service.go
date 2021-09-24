package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/task"
)

type Service struct {
	rep task.Repository
}

func NewService(repository task.Repository) *Service {
	return &Service{
		rep: repository,
	}
}

func (s *Service) GetTasksList() ([]*models.Task, error) {
	return s.rep.SelectAll()
}

func (s *Service) GetTaskByID(id models.ID) (*models.Task, error) {
	return s.rep.SelectByID(id)
}

func (s *Service) GetTasksByUserID(id models.ID) ([]*models.Task, error) {
	return s.rep.SelectByUserID(id)
}

func (s *Service) GetTasksByTeamID(id models.ID) ([]*models.Task, error) {
	return s.rep.SelectByTeamID(id)
}

func (s *Service) CreateTask(t *models.Task) error {
	return s.rep.Insert(t)
}

func (s *Service) UpdateTask(t *models.Task) error {
	return s.rep.Update(t)
}

func (s *Service) DeleteTaskByID(id models.ID) error {
	return s.rep.DeleteByID(id)
}

func (s *Service) DeleteTaskByUserID(id models.ID) error {
	return s.rep.DeleteByUserID(id)
}

func (s *Service) DeleteTaskByTeamID(id models.ID) error {
	return s.rep.DeleteByTeamID(id)
}
