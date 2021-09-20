package models

import (
	"time"
)

type User struct {
	ID        ID        `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ImageID   ID        `json:"image_id,omitempty"`
}

type Users []User

func NewUser(login, email, password string) *User {
	return &User{
		ID:        NewID(),
		Login:     login,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		ImageID:   ID{},
	}
}
