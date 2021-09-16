package entity

import (
	"time"
)

type User struct {
	ID        ID
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
	ImageID   ID
}

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
