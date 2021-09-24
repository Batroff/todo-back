package task

import "github.com/batroff/todo-back/internal/models"

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	SelectAll() ([]*models.Task, error)
	SelectByID(models.ID) (*models.Task, error)
	SelectByUserID(models.ID) ([]*models.Task, error)
	SelectByTeamID(models.ID) ([]*models.Task, error)
}

type Writer interface {
	Insert(*models.Task) error
	Update(*models.Task) error
	DeleteByID(models.ID) error
	DeleteByUserID(models.ID) error
	DeleteByTeamID(models.ID) error
}
