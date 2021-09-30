package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/todo"
)

type Service struct {
	rep todo.Repository
}

func NewService(r todo.Repository) *Service {
	return &Service{
		rep: r,
	}
}

func (s *Service) GetTodosList() ([]*models.Todo, error) {
	return s.rep.SelectAll()
}

func (s *Service) GetTodoByID(id models.ID) (*models.Todo, error) {
	return s.rep.SelectByID(id)
}

func (s *Service) CreateTodo(t *models.Todo) error {
	t.ID = models.NewID()

	return s.rep.Insert(t)
}

func (s *Service) UpdateTodo(t *models.Todo) error {
	return s.rep.Update(t)
}

func (s *Service) DeleteTodo(id models.ID) error {
	return s.rep.Delete(id)
}
