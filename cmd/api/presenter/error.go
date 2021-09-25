package presenter

import "errors"

var ErrNotImplementJsonMarshaller = errors.New("object doesn't implement json.Marshaller")

var ErrBadRequest = errors.New("cannot unmarshal request object")

var ErrUnauthorized = errors.New("wrong credentials (email or password)")
