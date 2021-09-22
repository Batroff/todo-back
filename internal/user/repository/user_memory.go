package repository

import "github.com/batroff/todo-back/internal/models"

type MemRepo struct {
	rep map[models.ID]*models.User
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		rep: make(map[models.ID]*models.User),
	}
}

func (mem *MemRepo) SelectByID(id models.ID) (*models.User, error) {
	if u, ok := mem.rep[id]; ok {
		return u, nil
	}

	return nil, models.ErrNotFound
}
func (mem *MemRepo) SelectBy(key string, value interface{}) (u []*models.User, err error) {
	return nil, nil
}

func (mem *MemRepo) SelectByEmail(email string) (*models.User, error) {
	for _, u := range mem.rep {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, models.ErrNotFound
}

func (mem *MemRepo) SelectAll() (users []*models.User, err error) {
	for _, u := range mem.rep {
		users = append(users, u)
	}

	return users, nil
}

func (mem *MemRepo) Insert(u *models.User) (id models.ID, err error) {
	if _, ok := mem.rep[u.ID]; ok {
		return u.ID, models.ErrAlreadyExists
	}
	mem.rep[u.ID] = u

	return u.ID, nil
}
func (mem *MemRepo) Update(u *models.User) error {
	if _, ok := mem.rep[u.ID]; ok {
		mem.rep[u.ID] = u
		return nil
	}

	return models.ErrNotFound
}
func (mem *MemRepo) Delete(id models.ID) error {
	if _, ok := mem.rep[id]; ok {
		delete(mem.rep, id)
		return nil
	}

	return models.ErrNotFound
}
