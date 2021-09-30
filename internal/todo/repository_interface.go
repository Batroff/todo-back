package todo

import "github.com/batroff/todo-back/internal/models"

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	SelectAll() ([]*models.Todo, error)
	SelectByID(models.ID) (*models.Todo, error)
}

type Writer interface {
	Insert(*models.Todo) error
	Update(*models.Todo) error
	Delete(models.ID) error
}
