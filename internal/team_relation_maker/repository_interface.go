package relation_maker

import "github.com/batroff/todo-back/internal/models"

type Repository interface {
	RepositoryReader
	RepositoryWriter
}

type RepositoryReader interface {
	SelectByUserID(models.ID) ([]*models.UserTeamRel, error)
	SelectByTeamID(models.ID) ([]*models.UserTeamRel, error)
	SelectByIDs(teamID, userID models.ID) (*models.UserTeamRel, error)
}

type RepositoryWriter interface {
	Insert(*models.UserTeamRel) error
	DeleteByIDs(teamID, userID models.ID) error
}
