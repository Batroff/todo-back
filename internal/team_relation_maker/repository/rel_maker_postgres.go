package repository

import (
	"database/sql"
	"github.com/batroff/todo-back/internal/models"
	"log"
)

type TeamRelationMakerPostgres struct {
	db *sql.DB
}

func NewTeamRelationMakerPostgres(db *sql.DB) *TeamRelationMakerPostgres {
	return &TeamRelationMakerPostgres{
		db: db,
	}
}

func selectByID(relMaker *TeamRelationMakerPostgres, query string, id models.ID) (relations []*models.UserTeamRel, err error) {
	rows, err := relMaker.db.Query(query, id)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := new(models.UserTeamRel)

		if err := rows.Scan(
			&r.UserID,
			&r.TeamID,
		); err != nil {
			return nil, err
		}

		relations = append(relations, r)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("closing rows in TeamRelationMakerPostgres error: %s\n", err)
		}
	}(rows)

	return relations, nil
}

func (relMaker *TeamRelationMakerPostgres) SelectByUserID(id models.ID) (relations []*models.UserTeamRel, err error) {
	return selectByID(relMaker, "SELECT * FROM users_team_xref WHERE id_user = $1;", id)
}

func (relMaker *TeamRelationMakerPostgres) SelectByTeamID(id models.ID) (relations []*models.UserTeamRel, err error) {
	return selectByID(relMaker, "SELECT * FROM users_team_xref WHERE id_team = $1;", id)
}

func (relMaker *TeamRelationMakerPostgres) SelectByIDs(teamID, userID models.ID) (*models.UserTeamRel, error) {
	rel := new(models.UserTeamRel)

	if err := relMaker.db.QueryRow("SELECT * FROM users_team_xref WHERE id_user = $1 AND id_team = $2",
		userID, teamID).Scan(
		&rel.UserID,
		&rel.TeamID,
	); err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return rel, nil
}

func (relMaker *TeamRelationMakerPostgres) Insert(rel *models.UserTeamRel) error {
	if _, err := relMaker.db.Exec("INSERT INTO users_team_xref(id_user, id_team) VALUES($1, $2);", rel.UserID, rel.TeamID); err != nil {
		return err
	}

	return nil
}

func (relMaker *TeamRelationMakerPostgres) DeleteByIDs(teamID, userID models.ID) error {
	if _, err := relMaker.db.Exec("DELETE FROM users_team_xref WHERE id_team = $1 AND id_user = $2;", teamID, userID); err != nil {
		return err
	}

	return nil
}

// DeleteByTeamID TODO : add tests
func (relMaker *TeamRelationMakerPostgres) DeleteByTeamID(id models.ID) error {
	if _, err := relMaker.db.Exec("DELETE FROM users_team_xref WHERE id_team = $1;", id); err != nil {
		return err
	}

	return nil
}
