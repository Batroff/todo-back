package task

import "github.com/batroff/todo-back/internal/models"

type UseCase interface {
	GetTasksList() ([]*models.Task, error)
	GetTaskByID(models.ID) (*models.Task, error)
	GetTasksByUserID(models.ID) ([]*models.Task, error)
	GetTasksByTeamID(models.ID) ([]*models.Task, error)

	CreateTask(*models.Task) error
	UpdateTask(*models.Task) error
	DeleteTaskByID(models.ID) error
	DeleteTaskByUserID(models.ID) error
	DeleteTaskByTeamID(models.ID) error
}
