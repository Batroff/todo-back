package models

import (
	"time"
)

type User struct {
	ID        ID        `json:"id,omitempty"`
	Login     string    `json:"login,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ImageID   *ID       `json:"image_id"`
}

type Users []User

func NewUser(login, email, password string) *User {
	return &User{
		ID:        NewID(),
		Login:     login,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		ImageID:   nil,
	}
}
