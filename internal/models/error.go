package models

import "errors"

var ErrNotFound = errors.New("entities not found")

var ErrAlreadyExists = errors.New("entity already exists in repo")
