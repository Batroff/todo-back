package team

import "github.com/batroff/todo-back/internal/models"

type Repository interface {
	RepositoryReader
	RepositoryWriter
}

type RepositoryReader interface {
	SelectByID(models.ID) (*models.Team, error)
	SelectList() ([]*models.Team, error)
}

type RepositoryWriter interface {
	Insert(*models.Team) error
	Update(*models.Team) error
	Delete(models.ID) error
}
