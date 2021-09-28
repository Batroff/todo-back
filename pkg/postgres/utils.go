// Package postgres provides useful checks for postgres db
package postgres

import (
	"database/sql"
	"errors"
)

type Helper struct {
	db *sql.DB
}

// NewHelper returns *Helper, which provides methods to sql.DB
func NewHelper(db *sql.DB) *Helper {
	return &Helper{db: db}
}

// IsColExists checks existence of tableName.colName
// Returns nil if column exists
func (h *Helper) IsColExists(tableName, colName string) error {
	var exists bool

	err := h.db.QueryRow(
		"SELECT 1 as exists FROM information_schema.columns WHERE table_name = $1 AND column_name = $2;",
		tableName, colName,
	).Scan(&exists)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("column doesn't exist")
	}

	return nil
}
