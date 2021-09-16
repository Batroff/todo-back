package entity

type ErrExpectedOneEntity struct{}

func (err *ErrExpectedOneEntity) Error() string {
	return "Expected only one entity, got many"
}

type ErrNotFound struct{}

func (err *ErrNotFound) Error() string {
	return "Entities not found"
}

type ErrEntityAlreadyExists struct{}

func (err *ErrEntityAlreadyExists) Error() string {
	return "Entity already exists in repo"
}
