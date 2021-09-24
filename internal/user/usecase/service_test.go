package usecase

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/user/repository"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func newFixtureUser() *models.User {
	return models.NewUser(faker.Username(), faker.Email(), faker.Password())
}

func newTestService() *Service {
	userRepo := repository.NewMemRepo()
	return NewService(userRepo)
}

type userTestTable struct {
	login    string
	email    string
	password string
}

// TODO : add expected[errors, values] in testTable, make testTable for every test
func getTestTable() []userTestTable {
	return []userTestTable{
		{
			login:    "username",
			email:    "username@localhost",
			password: "p@ssw0rd",
		},
		{
			login:    faker.Username(),
			email:    faker.Email(),
			password: faker.Password(),
		},
	}
}

func TestService_GetUser(t *testing.T) {
	testTable := getTestTable()

	userService := newTestService()

	for _, testCase := range testTable {
		id, _ := userService.CreateUser(testCase.login, testCase.email, testCase.password)

		u, err := userService.GetUser(id)
		assert.NoError(t, err, err)

		expectedHash, err := bcrypt.GenerateFromPassword([]byte(testCase.password), bcrypt.DefaultCost)
		assert.Conditionf(t, func() bool {
			pwdErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(testCase.password))
			if u.Login != testCase.login || u.Email != testCase.email || pwdErr != nil {
				return false
			}

			return true
		}, "Expected <login, email, password>=<%q, %q, %q>\nGot <%q, %q, %q>",
			testCase.login, testCase.email, string(expectedHash), u.Login, u.Email, u.Password)
	}
}

func TestService_CreateUser(t *testing.T) {
	testTable := getTestTable()

	userService := newTestService()

	for _, testCase := range testTable {
		id, err := userService.CreateUser(testCase.login, testCase.email, testCase.password)
		assert.NoError(t, err, err)
		assert.Truef(t, models.IsIDValid(id), "Expected valid (not nil) uuid, got %q", id)

		u, _ := userService.GetUser(id)
		assert.NotEqualValues(t, time.Time{}, u.CreatedAt, "Expected not zero time")
		assert.Nilf(t, u.ImageID, "Expected nil uuid, got %s", u.ImageID)

		expectedHash, _ := bcrypt.GenerateFromPassword([]byte(testCase.password), bcrypt.DefaultCost)
		assert.Conditionf(t, func() bool {
			pwdErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(testCase.password))
			if u.Login != testCase.login || u.Email != testCase.email || pwdErr != nil {
				return false
			}

			return true
		}, "Expected <login, email, password>=<%q, %q, %q>\nGot <%q, %q, %q>",
			testCase.login, testCase.email, expectedHash, u.Login, u.Email, u.Password)
	}
}

func TestService_FindUserByEmail(t *testing.T) {
	testTable := getTestTable()

	userService := newTestService()

	for _, testCase := range testTable {
		id, _ := userService.CreateUser(testCase.login, testCase.email, testCase.password)

		u, err := userService.FindUserByEmail(testCase.email)
		assert.NoError(t, err, err)
		assert.EqualValuesf(t, id, u.ID, "Expected id=%s, got %s", id, u.ID)
		assert.EqualValuesf(t, testCase.email, u.Email,
			"Expected id=%s, got %s",
			testCase.email, u.Email)
	}
}

func TestService_GetUsersList(t *testing.T) {
	testTable := []struct {
		users []*models.User
	}{
		{
			users: nil,
		},
		{
			users: []*models.User{newFixtureUser(), newFixtureUser(), newFixtureUser()},
		},
	}

	userService := newTestService()

	for _, testCase := range testTable {
		var ids []models.ID
		for _, u := range testCase.users {
			id, _ := userService.CreateUser(u.Login, u.Email, u.Password)
			ids = append(ids, id)
		}

		users, err := userService.GetUsersList()
		assert.NoError(t, err, err)
		for _, u := range users {
			assert.Contains(t, ids, u.ID)
		}
	}
}

func TestService_DeleteUser(t *testing.T) {
	testTable := getTestTable()

	userService := newTestService()

	for _, testCase := range testTable {
		id, _ := userService.CreateUser(testCase.login, testCase.email, testCase.password)

		err := userService.DeleteUser(id)
		assert.NoError(t, err, err)

		// Try to get user
		u, err := userService.GetUser(id)
		assert.ErrorIs(t, err, models.ErrNotFound)
		assert.Nil(t, u)
	}
}

func TestService_UpdateUser(t *testing.T) {
	testTable := []struct {
		createdUser models.User
		updatedUser models.User
		expectedErr error
	}{
		{
			createdUser: models.User{},
			updatedUser: models.User{Login: "logUpdated"},
			expectedErr: nil,
		},
	}
	userService := newTestService()

	for _, testCase := range testTable {
		id, _ := userService.CreateUser(testCase.createdUser.Login, testCase.createdUser.Email, testCase.createdUser.Password)

		testCase.updatedUser.ID = id
		err := userService.UpdateUser(&testCase.updatedUser)
		assert.ErrorIs(t, err, testCase.expectedErr)

		u, _ := userService.GetUser(id)
		assert.Equal(t, u.Login, testCase.updatedUser.Login)
	}
}
