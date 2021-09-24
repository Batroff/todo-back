package presenter

import (
	"github.com/batroff/todo-back/internal/models"
)

type RequestUser struct {
	Login    *string    `json:"login,omitempty"`
	Email    *string    `json:"email,omitempty"`
	Password *string    `json:"password,omitempty"`
	ImageID  *models.ID `json:"image_id,omitempty"`
}
