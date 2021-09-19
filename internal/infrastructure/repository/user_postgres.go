package repository

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/internal/entity"
	"log"
)

type UserPostgres struct {
	db *sql.DB
}

// NewUserMySQL : create UserPostgres repository
func NewUserMySQL(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// Get : find in repository ONLY ONE user<entity.User> by entity.ID
func (userPostgres *UserPostgres) Get(id entity.ID) (*entity.User, error) {
	u := new(entity.User)
	err := userPostgres.db.QueryRow("select * from users where id_user = $1", id).Scan(
		&u.ID,
		&u.Login,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.ImageID,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

// Find : find in repository users<entity.User> by query
func (userPostgres *UserPostgres) Find(key string, value interface{}) (users []*entity.User, err error) {
	rows, err := userPostgres.db.Query(fmt.Sprintf("select * from users where %s = $1", key), value)
	if err != nil {
		return nil, err
	}
	// TODO : ? Create general func for appending users through rows iteration
	for rows.Next() {
		u := new(entity.User)

		err := rows.Scan(
			&u.ID,
			&u.Login,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.ImageID,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// TODO: logger in UserPostgres struct
			log.Printf("Error during closing rows in UserPostgres")
		}
	}()

	return users, nil
}

// List : return ALL users<entity.User> of repository
func (userPostgres *UserPostgres) List() (users []*entity.User, err error) {
	rows, err := userPostgres.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := new(entity.User)

		err := rows.Scan(
			&u.ID,
			&u.Login,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.ImageID,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error during closing rows in UserPostgres")
		}
	}()

	return users, nil
}

// Create : create user<entity.User> in repository
func (userPostgres *UserPostgres) Create(u *entity.User) (entity.ID, error) {
	query, err := userPostgres.db.Prepare(`insert into users(id_user, login, email, password, created_at, id_image) values($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return entity.ID{}, err
	}

	if &u.ImageID == new(entity.ID) {
		_, err = query.Exec(u.ID, u.Login, u.Email, u.Password, u.CreatedAt, nil)
	} else {
		_, err = query.Exec(u.ID, u.Login, u.Email, u.Password, u.CreatedAt, u.ImageID)
	}

	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

// Update : update user<entity.User> in repository
func (userPostgres *UserPostgres) Update(u *entity.User) error {
	query, err := userPostgres.db.Prepare("update users set login = $1, email = $2, password = $3, id_image = $4 where id_user = $5")
	if err != nil {
		return err
	}

	_, err = query.Exec(u.Login, u.Email, u.Password, u.ImageID, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete : delete user<entity.User> in repository
func (userPostgres *UserPostgres) Delete(id entity.ID) error {
	res, err := userPostgres.db.Exec(`delete from users where id_user = $1`, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); rows == 0 && err == nil {
		return entity.ErrNotFound
	}

	return nil
}
