package repository

import (
	"database/sql"
	"github.com/batroff/todo-back/internal/models"
	"log"
)

type TeamPostgres struct {
	db *sql.DB
}

func NewTeamPostgres(db *sql.DB) *TeamPostgres {
	return &TeamPostgres{db: db}
}

func (teamPostgres *TeamPostgres) SelectByID(id models.ID) (*models.Team, error) {
	t := new(models.Team)

	if err := teamPostgres.db.QueryRow("SELECT * FROM team WHERE id_team = $1;", id).Scan(
		&t.ID,
		&t.Name,
	); err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func (teamPostgres *TeamPostgres) SelectList() (teams []*models.Team, err error) {
	rows, err := teamPostgres.db.Query("SELECT * FROM team;")
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := new(models.Team)

		if err := rows.Scan(
			&t.ID,
			&t.Name,
		); err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("closing rows in TeamPostgres error %s\n", err)
		}
	}(rows)

	return teams, nil
}

func (teamPostgres *TeamPostgres) Insert(t *models.Team) error {
	if _, err := teamPostgres.db.Exec("INSERT INTO team(id_team, name) VALUES($1, $2);", t.ID, t.Name); err != nil {
		return err
	}

	return nil
}

func (teamPostgres *TeamPostgres) Update(t *models.Team) error {
	if _, err := teamPostgres.db.Exec("UPDATE team SET name = $1 WHERE id_team = $2;", t.Name, t.ID); err != nil {
		return err
	}

	return nil
}

func (teamPostgres *TeamPostgres) Delete(id models.ID) error {
	if _, err := teamPostgres.db.Exec("DELETE FROM team WHERE id_team = $1;", id); err != nil {
		return err
	}

	return nil
}
