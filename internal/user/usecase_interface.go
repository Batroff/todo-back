package user

import (
	"github.com/batroff/todo-back/internal/models"
)

//go:generate mockgen -source=usecase_interface.go -destination=mock/usecase_mock.go

type UseCase interface {
	GetUser(id models.ID) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUsersBy(key string, value interface{}) ([]*models.User, error)
	GetUsersList() ([]*models.User, error)
	CreateUser(login, email, password string) (models.ID, error) // TODO : Input param should be *models.User
	UpdateUser(u *models.User) error
	DeleteUser(id models.ID) error
}
