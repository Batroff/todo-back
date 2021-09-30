package testing

import "github.com/batroff/todo-back/internal/models"

type TestID struct {
	Input    string
	Expected models.ID
}

func IdFixture() TestID {
	id := models.NewID()
	return TestID{
		Input:    id.String(),
		Expected: id,
	}
}
