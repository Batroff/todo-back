package models

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() ID {
	return uuid.New()
}

// IsIDValid : returns true if id is not nil
func IsIDValid(id ID) bool {
	return id != [16]byte{}
}
