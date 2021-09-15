package repository

import (
	"database/sql"
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
func (userPostgres *UserPostgres) Get(id entity.ID) (u *entity.User, err error) {
	err = userPostgres.db.QueryRow("select * from users where id = $1", id).Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Find : find in repository users<entity.User> by query
func (userPostgres *UserPostgres) Find(query string) (users []*entity.User, err error) {
	rows, err := userPostgres.db.Query("select * from users where login = $1", query)
	if err != nil {
		return nil, err
	}

	// TODO : Create general func for appending users through rows iteration
	for rows.Next() {
		u := new(entity.User)

		err := rows.Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// TODO: logger in UserPostgres struct
			log.Fatalf("Error during closing rows in UserPostgres")
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

		err := rows.Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("Error during closing rows in UserPostgres")
		}
	}()

	return users, nil
}

// Create : create user<entity.User> in repository
func (userPostgres *UserPostgres) Create(u *entity.User) (id entity.ID, err error) {
	_, err = userPostgres.db.Exec(`insert into users(id, login, password, createdAt) values($1, $2, $3, $4)`,
		u.ID, u.Login, u.Password, u.CreatedAt)
	if err != nil {
		return entity.ID{}, err
	}

	return u.ID, nil
}

// Update : update user<entity.User> in repository
func (userPostgres *UserPostgres) Update(id entity.ID, u *entity.User) error {
	// TODO : implement Update
	return nil
}

// Delete : delete user<entity.User> in repository
func (userPostgres *UserPostgres) Delete(id entity.ID) error {
	// TODO : implement Delete
	return nil
}
