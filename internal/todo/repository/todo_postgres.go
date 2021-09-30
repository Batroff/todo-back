package repository

import (
	"database/sql"
	"github.com/batroff/todo-back/internal/models"
	"log"
)

type TodoPostgres struct {
	db *sql.DB
}

func NewTodoPostgres(db *sql.DB) *TodoPostgres {
	return &TodoPostgres{
		db: db,
	}
}

func (todoPostgres *TodoPostgres) SelectAll() (todos []*models.Todo, err error) {
	rows, err := todoPostgres.db.Query("SELECT * FROM todo;")
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := new(models.Todo)

		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Text,
			&t.Complete,
			&t.TaskID,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error while closing rows in TodoPostgres: %s", err)
		}
	}()

	return todos, nil
}

func (todoPostgres *TodoPostgres) SelectByID(id models.ID) (t *models.Todo, err error) {
	t = new(models.Todo)
	err = todoPostgres.db.QueryRow("SELECT * FROM todo WHERE id_todo = $1", id).Scan(
		&t.ID,
		&t.Title,
		&t.Text,
		&t.Complete,
		&t.TaskID,
	)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func (todoPostgres *TodoPostgres) Insert(t *models.Todo) error {
	_, err := todoPostgres.db.Exec("INSERT INTO todo(id_todo, title, text, complete, id_task) VALUES($1, $2, $3, $4, $5)", t.ID, t.Title, t.Text, t.Complete, t.TaskID)
	if err != nil {
		return err
	}

	return nil
}

func (todoPostgres *TodoPostgres) Update(t *models.Todo) error {
	_, err := todoPostgres.db.Exec("UPDATE todo SET title = $1, text = $2, complete = $3, id_task = $4", t.Title, t.Text, t.Complete, t.TaskID)
	if err != nil {
		return err
	}

	return nil
}

func (todoPostgres *TodoPostgres) Delete(id models.ID) error {
	_, err := todoPostgres.db.Exec("DELETE * FROM todo WHERE id_todo = $1", id)
	if err != nil {
		return err
	}

	return nil
}
