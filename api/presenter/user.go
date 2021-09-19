package presenter

import "github.com/batroff/todo-back/internal/entity"

type User struct {
	Login    string     `json:"login"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	ImageID  *entity.ID `json:"image_id,omitempty"`
}
