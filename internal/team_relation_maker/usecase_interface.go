package relation_maker

import "github.com/batroff/todo-back/internal/models"

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
}
