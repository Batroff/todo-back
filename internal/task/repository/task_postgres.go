package repository

import (
	"database/sql"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/pkg/postgres"
	"log"
	"strings"
)

type TaskPostgres struct {
	db     *sql.DB
	helper *postgres.Helper
}

func NewTaskPostgres(db *sql.DB) *TaskPostgres {
	return &TaskPostgres{
		db:     db,
		helper: postgres.NewHelper(db),
	}
}

func scanRows(rows *sql.Rows) (tasks []*models.Task, err error) {
	for rows.Next() {
		task := new(models.Task)

		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Priority,
			&task.UserID,
			&task.TeamID,
		)
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else if err != nil {
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

func (taskPostgres *TaskPostgres) SelectAll() (tasks []*models.Task, err error) {
	rows, err := taskPostgres.db.Query("SELECT * FROM task;")
	if err != nil {
		return nil, err
	}

	return scanRows(rows)
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

// SelectBy takes map of params and returns filtered []*models.Task
// If param doesn't exist in repo returns error
func (taskPostgres *TaskPostgres) SelectBy(pairs map[string]interface{}) (tasks []*models.Task, err error) {
	if len(pairs) == 0 {
		return taskPostgres.SelectAll()
	}
	whereQueries := make([]string, len(pairs))
	whereValues := make([]interface{}, len(pairs))

	i := 0
	for k, v := range pairs {
		if err = taskPostgres.helper.IsColExists("task", k); err != nil {
			return nil, err
		}

		whereQueries[i] = fmt.Sprintf("%s = $%d", k, i+1)
		whereValues[i] = v
		i += 1
	}

	query := fmt.Sprintf("SELECT * FROM task WHERE %s", strings.Join(whereQueries, " AND "))
	rows, err := taskPostgres.db.Query(query, whereValues...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows)
}

func (taskPostgres *TaskPostgres) SelectByUserID(id models.ID) (tasks []*models.Task, err error) {
	rows, err := taskPostgres.db.Query("SELECT * FROM task WHERE id_user = $1", id)
	if err != nil {
		return nil, err
	}

	return scanRows(rows)
}

func (taskPostgres *TaskPostgres) SelectByTeamID(id models.ID) (tasks []*models.Task, err error) {
	rows, err := taskPostgres.db.Query("SELECT * FROM task WHERE id_team = $1", id)
	if err != nil {
		return nil, err
	}

	return scanRows(rows)
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

func (taskPostgres *TaskPostgres) Update(t *models.Task) error {
	if _, err := taskPostgres.db.Exec(
		"UPDATE task SET title = $1, priority = $2, id_team = $3 WHERE id_task = $4",
		t.Title, t.Priority, t.TeamID, t.ID,
	); err != nil {
		return err
	}

	return nil
}

func (taskPostgres *TaskPostgres) DeleteByID(id models.ID) error {
	if _, err := taskPostgres.db.Exec("DELETE FROM task WHERE id_task = $1", id); err != nil {
		return err
	}

	return nil
}

func (taskPostgres *TaskPostgres) DeleteByUserID(models.ID) error {
	return nil
}

func (taskPostgres *TaskPostgres) DeleteByTeamID(models.ID) error {
	return nil
}
