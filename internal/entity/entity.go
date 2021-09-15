package entity

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() ID {
	return uuid.New()
}

func IDFromString(idStr string) (id ID, err error) {
	return uuid.FromBytes([]byte(idStr))
}
