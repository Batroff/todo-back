package presenter

import (
	"github.com/batroff/todo-back/internal/models"
)

type User struct {
	Login    string     `json:"login"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	ImageID  *models.ID `json:"image_id,omitempty"`
}
