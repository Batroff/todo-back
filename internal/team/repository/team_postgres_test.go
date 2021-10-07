package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/batroff/todo-back/internal/models"
	testUtils "github.com/batroff/todo-back/pkg/testing"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type sqlMockBehavior func(mock sqlmock.Sqlmock, query string, args ...driver.Value)

func TestTeamPostgres_SelectList(t *testing.T) {
	fixtureID := testUtils.IdFixture()
	testTable := []struct {
		name          string
		inputQuery    string
		inputArgs     []driver.Value
		sqlMock       sqlMockBehavior
		expectedTeams []*models.Team
		expectedError error
	}{
		{
			name: "OK",
			sqlMock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				rows := mock.NewRows([]string{"id_team", "name"})
				rows.AddRow(
					driver.Value(fixtureID.Expected),
					driver.Value(""),
				)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			inputQuery: `SELECT * FROM team;`,
			inputArgs:  nil,
			expectedTeams: []*models.Team{
				{ID: fixtureID.Expected},
			},
			expectedError: nil,
		},
		{
			name: "Not found any team",
			sqlMock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(sql.ErrNoRows)
			},
			inputQuery:    `SELECT * FROM team;`,
			inputArgs:     nil,
			expectedTeams: nil,
			expectedError: models.ErrNotFound,
		},
		{
			name: "Internal sql error",
			sqlMock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(errors.New("internal sql error"))
			},
			inputQuery:    `SELECT * FROM team;`,
			inputArgs:     nil,
			expectedTeams: nil,
			expectedError: errors.New("internal sql error"),
		},
		{
			name: "Rows scan error",
			sqlMock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				rows := mock.NewRows([]string{"id_team", "name"}).
					AddRow("invalid", "#1")

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			inputQuery:    `SELECT * FROM team;`,
			inputArgs:     nil,
			expectedTeams: nil,
			expectedError: errors.New("sql: Scan error on column index 0, name \"id_team\": Scan: invalid UUID length: 12"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock sql
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.sqlMock(mock, testCase.inputQuery, testCase.inputArgs)

			// Repo & repo method
			teamRepo := NewTeamPostgres(db)
			teams, err := teamRepo.SelectList()

			expectationsErr := mock.ExpectationsWereMet()

			// Assert
			assert.Equal(t, testCase.expectedTeams, teams)
			// FIXME : if err is nil -> test break
			if err != nil {
				assert.Error(t, testCase.expectedError, err)
			}
			assert.NoError(t, expectationsErr)
		})
	}
}

func TestTeamPostgres_SelectByID(t *testing.T) {
	fixtureID := testUtils.IdFixture()

	testTable := []struct {
		name          string
		testID        testUtils.TestID
		inputQuery    string
		inputArgs     []driver.Value
		mock          sqlMockBehavior
		expectedTeam  *models.Team
		expectedError error
	}{
		{
			name:       "OK",
			testID:     fixtureID,
			inputQuery: `SELECT * FROM team WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Input),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				rows := sqlmock.NewRows([]string{"id_team", "name"}).
					AddRow(
						fixtureID.Expected,
						"#1",
					)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedTeam: &models.Team{
				ID:   fixtureID.Expected,
				Name: "#1",
			},
			expectedError: nil,
		},
		{
			name:       "Error no rows",
			testID:     fixtureID,
			inputQuery: `SELECT * FROM team WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Input),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(sql.ErrNoRows)
			},
			expectedTeam:  nil,
			expectedError: models.ErrNotFound,
		},
		{
			name:       "sql internal error",
			testID:     fixtureID,
			inputQuery: `SELECT * FROM team WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Input),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(errors.New("internal error"))
			},
			expectedTeam:  nil,
			expectedError: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock sql
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// mock
			teamRepo := NewTeamPostgres(db)
			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs...)

			expTeam, err := teamRepo.SelectByID(testCase.testID.Expected)

			// Assert
			assert.Equal(t, testCase.expectedTeam, expTeam)
			assert.Equal(t, testCase.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamPostgres_Insert(t *testing.T) {
	fixtureID := testUtils.IdFixture()

	testTable := []struct {
		name          string
		inputTeam     *models.Team
		inputQuery    string
		inputArgs     []driver.Value
		mock          sqlMockBehavior
		expectedError error
	}{
		{
			name: "OK",
			inputTeam: &models.Team{
				ID:   fixtureID.Expected,
				Name: "some name",
			},
			inputQuery: `INSERT INTO team(id_team, name) VALUES($1, $2);`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Expected),
				driver.Value("some name"),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			expectedError: nil,
		},
		{
			name: "Exec error",
			inputTeam: &models.Team{
				ID:   fixtureID.Expected,
				Name: "some name",
			},
			inputQuery: `INSERT INTO team(id_team, name) VALUES($1, $2);`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Expected),
				driver.Value("some name"),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(nil).
					WillReturnError(errors.New("internal error"))
			},
			expectedError: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock sql
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			teamRepo := NewTeamPostgres(db)
			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs...)

			err = teamRepo.Insert(testCase.inputTeam)

			// Assert
			assert.Equal(t, testCase.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamPostgres_Update(t *testing.T) {
	fixtureID := testUtils.IdFixture()

	testTable := []struct {
		name          string
		inputTeam     *models.Team
		inputQuery    string
		inputArgs     []driver.Value
		mock          sqlMockBehavior
		expectedError error
	}{
		{
			name: "OK",
			inputTeam: &models.Team{
				ID:   fixtureID.Expected,
				Name: "some name",
			},
			inputQuery: `UPDATE team SET name = $1 WHERE id_team = $2;`,
			inputArgs: []driver.Value{
				driver.Value("some name"),
				driver.Value(fixtureID.Expected),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			expectedError: nil,
		},
		{
			name: "Exec error",
			inputTeam: &models.Team{
				ID:   fixtureID.Expected,
				Name: "some name",
			},
			inputQuery: `UPDATE team SET name = $1 WHERE id_team = $2;`,
			inputArgs: []driver.Value{
				driver.Value("some name"),
				driver.Value(fixtureID.Expected),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(nil).
					WillReturnError(errors.New("internal error"))
			},
			expectedError: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock sql
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			teamRepo := NewTeamPostgres(db)
			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs...)

			err = teamRepo.Update(testCase.inputTeam)

			// Assert
			assert.Equal(t, testCase.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamPostgres_Delete(t *testing.T) {
	fixtureID := testUtils.IdFixture()

	testTable := []struct {
		name          string
		testID        testUtils.TestID
		inputQuery    string
		inputArgs     []driver.Value
		mock          sqlMockBehavior
		expectedError error
	}{
		{
			name:       "OK",
			testID:     fixtureID,
			inputQuery: `DELETE FROM team WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Expected),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Exec error",
			inputQuery: `DELETE FROM team WHERE id_team = $1;`,
			testID:     fixtureID,
			inputArgs: []driver.Value{
				driver.Value(fixtureID.Expected),
			},
			mock: func(mock sqlmock.Sqlmock, query string, args ...driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(nil).
					WillReturnError(errors.New("internal error"))
			},
			expectedError: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock sql
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			teamRepo := NewTeamPostgres(db)
			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs...)

			err = teamRepo.Delete(testCase.testID.Expected)

			// Assert
			assert.Equal(t, testCase.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
