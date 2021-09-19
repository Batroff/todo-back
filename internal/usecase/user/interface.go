package user

import "github.com/batroff/todo-back/internal/entity"

type Reader interface {
	Get(id entity.ID) (u *entity.User, err error)
	Find(key string, value interface{}) (u []*entity.User, err error)
	List() (u []*entity.User, err error)
}

type Writer interface {
	Create(u *entity.User) (id entity.ID, err error)
	Update(u *entity.User) error
	Delete(id entity.ID) error
}

type UseCase interface {
	GetUser(id entity.ID) (u *entity.User, err error)
	FindUsersBy(key string, value interface{}) (u []*entity.User, err error)
	GetUsersList() (u []*entity.User, err error)
	CreateUser(login, email, password string) (id entity.ID, err error)
	UpdateUser(u *entity.User) error
	DeleteUser(id entity.ID) error
}

type Repository interface {
	Reader
	Writer
}
