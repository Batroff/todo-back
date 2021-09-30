package todo

import "github.com/batroff/todo-back/internal/models"

//go:generate mockgen -source=usecase_interface.go -destination=mock/usecase_mock.go

type UseCase interface {
	GetTodosList() ([]*models.Todo, error)
	GetTodoByID(models.ID) (*models.Todo, error)

	CreateTodo(*models.Todo) error
	UpdateTodo(*models.Todo) error
	DeleteTodo(models.ID) error
}
