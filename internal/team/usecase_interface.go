package team

import "github.com/batroff/todo-back/internal/models"

//go:generate mockgen -source=./usecase_interface.go -destination=./mock/usecase_mock.go
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
