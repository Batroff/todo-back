package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	rep user.Repository
}

func NewService(repository user.Repository) *Service {
	return &Service{rep: repository}
}

func (s *Service) GetUser(id models.ID) (u *models.User, err error) {
	return s.rep.SelectByID(id)
}

func (s *Service) FindUsersBy(key string, value interface{}) ([]*models.User, error) {
	users, err := s.rep.SelectBy(key, value)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, models.ErrNotFound
	}

	return users, nil
}

func (s *Service) FindUserByEmail(email string) (u *models.User, err error) {
	return s.rep.SelectByEmail(email)
}

func (s *Service) GetUsersList() (u []*models.User, err error) {
	return s.rep.SelectAll()
}

func (s *Service) CreateUser(login, email, password string) (id models.ID, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.ID{}, err
	}

	u := models.NewUser(login, email, string(hash))
	return s.rep.Insert(u)
}

func (s *Service) UpdateUser(u *models.User) error {
	return s.rep.Update(u)
}

func (s *Service) DeleteUser(id models.ID) error {
	return s.rep.Delete(id)
}
