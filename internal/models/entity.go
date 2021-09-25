// Package models provides lower level of application - model structs
// It also provides general types for every model like ID
// And methods like NewID
package models

import "github.com/google/uuid"

type ID = uuid.UUID

// NewID Generates random UUID
func NewID() ID {
	return uuid.New()
}

// IsIDValid returns true if id is not nil
func IsIDValid(id ID) bool {
	return id != [16]byte{}
}
