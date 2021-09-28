package presenter

import "github.com/batroff/todo-back/internal/models"

type RequestTask struct {
	Title    *string    `json:"title,omitempty"`
	Priority *uint      `json:"priority,omitempty"`
	TeamID   *models.ID `json:"id_team,omitempty"`
}
