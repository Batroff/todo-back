package models

import "errors"

var ErrNotFound = errors.New("entities not found")

var ErrAlreadyExists = errors.New("entity already exists in repo")

// HTTP errors

var ErrNotImplementJsonMarshaller = errors.New("object doesn't implement json.Marshaller")

var ErrBadRequest = errors.New("cannot unmarshal request object")

var ErrUnauthorized = errors.New("wrong credentials (email or password)")
