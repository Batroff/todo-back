package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/team"
)

type Service struct {
	rep team.Repository
}

func NewService(r team.Repository) *Service {
	return &Service{rep: r}
}

func (s *Service) SelectTeamByID(id models.ID) (*models.Team, error) {
	return s.rep.SelectByID(id)
}

func (s *Service) SelectTeamsList() ([]*models.Team, error) {
	return s.rep.SelectList()
}

func (s *Service) CreateTeam(t *models.Team) error {
	return s.rep.Insert(t)
}

func (s *Service) UpdateTeam(t *models.Team) error {
	return s.rep.Update(t)
}

func (s *Service) DeleteTeam(id models.ID) error {
	return s.rep.Delete(id)
}
