package entity

import "time"

type User struct {
	ID        ID
	Login     string
	Password  string
	CreatedAt time.Time
}

func NewUser(login, password string) *User {
	return &User{ID: NewID(), Login: login, Password: password, CreatedAt: time.Now()}
}
