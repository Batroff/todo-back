package repository

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	"log"
)

type UserPostgres struct {
	db *sql.DB
}

// NewUserMySQL : create UserPostgres repository
func NewUserMySQL(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// SelectByID : find in repository ONLY ONE user<entity.User> by entity.ID
func (userPostgres *UserPostgres) SelectByID(id models.ID) (*models.User, error) {
	u := new(models.User)
	err := userPostgres.db.QueryRow("select * from users where id_user = $1", id).Scan(
		&u.ID,
		&u.Login,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.ImageID,
	)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func (userPostgres *UserPostgres) SelectByEmail(email string) (*models.User, error) {
	u := new(models.User)
	err := userPostgres.db.QueryRow("select * from users where email = $1", email).Scan(
		&u.ID,
		&u.Login,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.ImageID,
	)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func scanUsersRows(rows *sql.Rows) (users []*models.User, err error) {
	for rows.Next() {
		u := new(models.User)

		err := rows.Scan(
			&u.ID,
			&u.Login,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.ImageID,
		)
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else if err != nil {
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

// SelectBy : find in repository users<entity.User> by query
func (userPostgres *UserPostgres) SelectBy(key string, value interface{}) (users []*models.User, err error) {
	rows, err := userPostgres.db.Query(fmt.Sprintf("select * from users where %s = $1", key), value)
	if err != nil {
		return nil, err
	}

	users, err = scanUsersRows(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// SelectAll : return ALL users<entity.User> of repository
func (userPostgres *UserPostgres) SelectAll() (users []*models.User, err error) {
	rows, err := userPostgres.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	users, err = scanUsersRows(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Insert : create user<entity.User> in repository
func (userPostgres *UserPostgres) Insert(u *models.User) (models.ID, error) {
	query, err := userPostgres.db.Prepare(`insert into users(id_user, login, email, password, created_at, id_image) values($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return models.ID{}, err
	}

	_, err = query.Exec(u.ID, u.Login, u.Email, u.Password, u.CreatedAt, u.ImageID)

	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

// Update : update user<entity.User> in repository
func (userPostgres *UserPostgres) Update(u *models.User) error {
	query, err := userPostgres.db.Prepare("update users set login = $1, email = $2, password = $3, id_image = $4 where id_user = $5")
	if err != nil {
		return err
	}

	res, err := query.Exec(u.Login, u.Email, u.Password, u.ImageID, u.ID)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); rows == 0 && err == nil {
		return models.ErrNotFound
	}

	return nil
}

// Delete : delete user<entity.User> in repository
func (userPostgres *UserPostgres) Delete(id models.ID) error {
	_, err := userPostgres.db.Exec(`delete from users where id_user = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
