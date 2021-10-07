package team

import "github.com/batroff/todo-back/internal/models"

type UseCase interface {
	UseCaseReader
	UseCaseWriter
}

type UseCaseReader interface {
	SelectTeamByID(models.ID) (*models.Team, error)
	SelectTeamsList() ([]*models.Team, error)
}

type UseCaseWriter interface {
	CreateTeam(*models.Team) error
	UpdateTeam(*models.Team) error
	DeleteTeam(models.ID) error
}
