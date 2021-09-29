package task

import "github.com/batroff/todo-back/internal/models"

//go:generate mockgen -source=usecase_interface.go -destination=mock/usecase_mock.go

type UseCase interface {
	GetTasksList() ([]*models.Task, error)
	GetTaskByID(models.ID) (*models.Task, error)
	GetTasksByUserID(models.ID) ([]*models.Task, error)
	GetTasksByTeamID(models.ID) ([]*models.Task, error)
	GetTasksBy(map[string]interface{}) ([]*models.Task, error)

	CreateTask(*models.Task) error
	UpdateTask(*models.Task) error
	DeleteTaskByID(models.ID) error
	DeleteTaskByUserID(models.ID) error
	DeleteTaskByTeamID(models.ID) error
}
