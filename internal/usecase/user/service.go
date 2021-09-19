package user

import (
	"github.com/batroff/todo-back/internal/entity"
)

type Service struct {
	rep Repository
}

func NewService(repository Repository) *Service {
	return &Service{rep: repository}
}

func (s *Service) GetUser(id entity.ID) (u *entity.User, err error) {
	return s.rep.Get(id)
}

func (s *Service) FindUsersBy(key string, value interface{}) ([]*entity.User, error) {
	users, err := s.rep.Find(key, value)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, entity.ErrNotFound
	}

	return users, nil
}

func (s *Service) GetUsersList() (u []*entity.User, err error) {
	return s.rep.List()
}

func (s *Service) CreateUser(login, email, password string) (id entity.ID, err error) {
	u := entity.NewUser(login, email, password)
	return s.rep.Create(u)
}

func (s *Service) UpdateUser(u *entity.User) error {
	return s.rep.Update(u)
}

func (s *Service) DeleteUser(id entity.ID) error {
	return s.rep.Delete(id)
}
