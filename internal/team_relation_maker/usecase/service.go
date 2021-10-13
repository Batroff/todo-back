package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	maker "github.com/batroff/todo-back/internal/team_relation_maker"
)

type Service struct {
	rep maker.Repository
}

func NewService(r maker.Repository) *Service {
	return &Service{
		rep: r,
	}
}

func (s *Service) SelectRelationsByUserID(id models.ID) ([]*models.UserTeamRel, error) {
	return s.rep.SelectByUserID(id)
}

func (s *Service) SelectRelationsByTeamID(id models.ID) ([]*models.UserTeamRel, error) {
	return s.rep.SelectByTeamID(id)
}

func (s *Service) SelectRelationByIDs(teamID, userID models.ID) (*models.UserTeamRel, error) {
	return s.rep.SelectByIDs(teamID, userID)
}

func (s *Service) CreateRelation(rel *models.UserTeamRel) error {
	return s.rep.Insert(rel)
}

func (s *Service) DeleteRelationByIDs(teamID, userID models.ID) error {
	return s.rep.DeleteByIDs(teamID, userID)
}

func (s *Service) DeleteRelationsByTeamID(id models.ID)	error {
	return s.rep.DeleteByTeamID(id)
}