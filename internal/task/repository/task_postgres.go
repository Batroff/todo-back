package repository

import (
	"database/sql"
	"github.com/batroff/todo-back/internal/models"
	"log"
)

type TaskPostgres struct {
	db *sql.DB
}

func NewTaskPostgres(db *sql.DB) *TaskPostgres {
	return &TaskPostgres{
		db: db,
	}
}

func (taskPostgres *TaskPostgres) SelectAll() (tasks []*models.Task, err error) {
	rows, err := taskPostgres.db.Query("SELECT * FROM task;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		task := new(models.Task)

		if err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Priority,
			&task.UserID,
			&task.TeamID,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error while closing sql.Rows in TaskPostgres: %s", err)
		}
	}(rows)

	return tasks, nil
}

func (taskPostgres *TaskPostgres) SelectByID(id models.ID) (*models.Task, error) {
	t := new(models.Task)

	err := taskPostgres.db.QueryRow("SELECT * FROM task WHERE id_task = $1", id).Scan(
		&t.ID,
		&t.Title,
		&t.Priority,
		&t.UserID,
		&t.TeamID,
	)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func (taskPostgres *TaskPostgres) SelectByUserID(models.ID) ([]*models.Task, error) {
	return nil, nil
}

func (taskPostgres *TaskPostgres) SelectByTeamID(models.ID) ([]*models.Task, error) {
	return nil, nil
}

func (taskPostgres *TaskPostgres) Insert(task *models.Task) error {
	query, err := taskPostgres.db.Prepare("INSERT INTO task(id_task, title, priority, id_user, id_team) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}

	_, err = query.Exec(task.ID, task.Title, task.Priority, task.UserID, task.TeamID)
	if err != nil {
		return err
	}

	return nil
}

func (taskPostgres *TaskPostgres) Update(*models.Task) error {
	return nil
}

func (taskPostgres *TaskPostgres) DeleteByID(models.ID) error {
	return nil
}

func (taskPostgres *TaskPostgres) DeleteByUserID(models.ID) error {
	return nil
}

func (taskPostgres *TaskPostgres) DeleteByTeamID(models.ID) error {
	return nil
}
