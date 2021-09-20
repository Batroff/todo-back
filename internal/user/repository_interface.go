package user

import "github.com/batroff/todo-back/internal/models"

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	SelectBy(key string, value interface{}) ([]*models.User, error)
	SelectByID(id models.ID) (*models.User, error)
	SelectByEmail(email string) (*models.User, error)
	SelectAll() ([]*models.User, error)
}

type Writer interface {
	Insert(u *models.User) (models.ID, error)
	Update(u *models.User) error
	Delete(id models.ID) error
}
