package relation_maker

import "github.com/batroff/todo-back/internal/models"

//go:generate mockgen -source=./usecase_interface.go -destination=./mock/usecase_mock.go
type UseCase interface {
	UseCaseReader
	UseCaseWriter
}

type UseCaseReader interface {
	SelectRelationsByUserID(models.ID) ([]*models.UserTeamRel, error)
	SelectRelationsByTeamID(models.ID) ([]*models.UserTeamRel, error)
	SelectRelationByIDs(teamID, userID models.ID) (*models.UserTeamRel, error)
}

type UseCaseWriter interface {
	CreateRelation(*models.UserTeamRel) error
	DeleteRelationByIDs(teamID, userID models.ID) error
	DeleteRelationsByTeamID(models.ID) error
}
