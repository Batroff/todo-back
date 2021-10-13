package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/batroff/todo-back/internal/models"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type teamRelationMakerMockBehavior func(mock sqlmock.Sqlmock, query string, args []driver.Value)

func TestTeamRelationMakerPostgres_SelectByUserID(t *testing.T) {
	testUserID := models.NewID()
	testTeamID := models.NewID()

	testTable := []struct {
		name              string
		inputQuery        string
		inputArgs         []driver.Value
		inputUserID       models.ID
		mock              teamRelationMakerMockBehavior
		expectedRelations []*models.UserTeamRel
		expectedError     error
	}{
		{
			name:       "OK",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1;`,
			inputArgs: []driver.Value{
				testUserID,
			},
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				rows := sqlmock.NewRows([]string{"id_user", "id_team"}).
					AddRow(testUserID, testTeamID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedRelations: []*models.UserTeamRel{
				{
					UserID: testUserID,
					TeamID: testTeamID,
				},
			},
			expectedError: nil,
		},
		{
			name:       "Not found relations",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1;`,
			inputArgs: []driver.Value{
				testUserID,
			},
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRelations: nil,
			expectedError:     models.ErrNotFound,
		},
		{
			name:       "Query unexpected error",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1;`,
			inputArgs: []driver.Value{
				testUserID,
			},
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(errors.New("unexpected"))
			},
			expectedRelations: nil,
			expectedError:     errors.New("unexpected"),
		},
		{
			name:       "Scan error",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1;`,
			inputArgs: []driver.Value{
				testUserID,
			},
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				rows := sqlmock.NewRows([]string{"id_user", "id_team"}).
					AddRow(true, testTeamID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedRelations: nil,
			expectedError:     errors.New("sql: Scan error on column index 0, name \"id_user\": Scan: unable to scan type bool into UUID"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init db mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs)

			// Init repo
			makerRepo := NewTeamRelationMakerPostgres(db)
			rel, err := makerRepo.SelectByUserID(testCase.inputUserID)

			// Assert
			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedRelations, rel)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamRelationMakerPostgres_SelectByTeamID(t *testing.T) {
	testTeamID := models.NewID()
	testUserID := models.NewID()

	testTable := []struct {
		name              string
		inputQuery        string
		inputArgs         []driver.Value
		inputTeamID       models.ID
		mock              teamRelationMakerMockBehavior
		expectedRelations []*models.UserTeamRel
		expectedError     error
	}{
		{
			name:       "OK",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				testTeamID,
			},
			inputTeamID: testTeamID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				rows := sqlmock.NewRows([]string{"id_user", "id_team"}).
					AddRow(testTeamID, testUserID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedRelations: []*models.UserTeamRel{
				{
					UserID: testTeamID,
					TeamID: testUserID,
				},
			},
			expectedError: nil,
		},
		{
			name:       "Not found relations",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				testTeamID,
			},
			inputTeamID: testTeamID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRelations: nil,
			expectedError:     models.ErrNotFound,
		},
		{
			name:       "Query unexpected error",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				testTeamID,
			},
			inputTeamID: testTeamID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(errors.New("unexpected"))
			},
			expectedRelations: nil,
			expectedError:     errors.New("unexpected"),
		},
		{
			name:       "Scan error",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_team = $1;`,
			inputArgs: []driver.Value{
				testTeamID,
			},
			inputTeamID: testTeamID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				rows := sqlmock.NewRows([]string{"id_user", "id_team"}).
					AddRow(true, testUserID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedRelations: nil,
			expectedError:     errors.New("sql: Scan error on column index 0, name \"id_user\": Scan: unable to scan type bool into UUID"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init db mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs)

			// Init repo
			makerRepo := NewTeamRelationMakerPostgres(db)
			rel, err := makerRepo.SelectByTeamID(testCase.inputTeamID)

			// Assert
			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedRelations, rel)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamRelationMakerPostgres_SelectByIDs(t *testing.T) {
	testUserID := models.NewID()
	testTeamID := models.NewID()

	testTable := []struct {
		name             string
		inputQuery       string
		inputArgs        []driver.Value
		inputTeamID      models.ID
		inputUserID      models.ID
		mock             teamRelationMakerMockBehavior
		expectedRelation *models.UserTeamRel
		expectedError    error
	}{
		{
			name:       "OK",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1 AND id_team = $2`,
			inputArgs: []driver.Value{
				testUserID,
				testTeamID,
			},
			inputTeamID: testTeamID,
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				rows := sqlmock.NewRows([]string{"id_user", "id_team"}).
					AddRow(testUserID, testTeamID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			expectedRelation: &models.UserTeamRel{
				UserID: testUserID,
				TeamID: testTeamID,
			},
			expectedError: nil,
		},
		{
			name:       "Not found",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1 AND id_team = $2`,
			inputArgs: []driver.Value{
				testUserID,
				testTeamID,
			},
			inputTeamID: testTeamID,
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRelation: nil,
			expectedError:    models.ErrNotFound,
		},
		{
			name:       "Scan error",
			inputQuery: `SELECT * FROM users_team_xref WHERE id_user = $1 AND id_team = $2`,
			inputArgs: []driver.Value{
				testUserID,
				testTeamID,
			},
			inputTeamID: testTeamID,
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(&sqlmock.Rows{}).
					WillReturnError(errors.New("scan error"))
			},
			expectedRelation: nil,
			expectedError:    errors.New("scan error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init db mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs)

			// Init repo
			makerRepo := NewTeamRelationMakerPostgres(db)
			rel, err := makerRepo.SelectByIDs(testCase.inputTeamID, testCase.inputUserID)

			// Assert
			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedRelation, rel)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamRelationMakerPostgres_Insert(t *testing.T) {
	testUserID := models.NewID()
	testTeamID := models.NewID()

	testTable := []struct {
		name          string
		inputQuery    string
		inputArgs     []driver.Value
		inputRel      *models.UserTeamRel
		mock          teamRelationMakerMockBehavior
		expectedError error
	}{
		{
			name:       "OK",
			inputQuery: `INSERT INTO users_team_xref(id_user, id_team) VALUES($1, $2);`,
			inputArgs: []driver.Value{
				testUserID,
				testTeamID,
			},
			inputRel: &models.UserTeamRel{
				UserID: testUserID,
				TeamID: testTeamID,
			},
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Exec error",
			inputQuery: `INSERT INTO users_team_xref(id_user, id_team) VALUES($1, $2);`,
			inputArgs: []driver.Value{
				testUserID,
				testTeamID,
			},
			inputRel: &models.UserTeamRel{
				UserID: testUserID,
				TeamID: testTeamID,
			},
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(nil).
					WillReturnError(errors.New("exec error"))
			},
			expectedError: errors.New("exec error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs)

			// Init repo
			makerRepo := NewTeamRelationMakerPostgres(db)
			err = makerRepo.Insert(testCase.inputRel)

			// Assert
			assert.Equal(t, testCase.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamRelationMakerPostgres_DeleteByIDs(t *testing.T) {
	testUserID := models.NewID()
	testTeamID := models.NewID()

	testTable := []struct {
		name        string
		inputQuery  string
		inputArgs   []driver.Value
		inputTeamID models.ID
		inputUserID models.ID
		mock        teamRelationMakerMockBehavior
		expectedErr error
	}{
		{
			name:       "OK",
			inputQuery: `DELETE * FROM users_team_xref WHERE id_team = $1 AND id_user = $2;`,
			inputArgs: []driver.Value{
				testTeamID,
				testUserID,
			},
			inputTeamID: testTeamID,
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Exec error",
			inputQuery: `DELETE * FROM users_team_xref WHERE id_team = $1 AND id_user = $2;`,
			inputArgs: []driver.Value{
				testTeamID,
				testUserID,
			},
			inputTeamID: testTeamID,
			inputUserID: testUserID,
			mock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(nil).
					WillReturnError(errors.New("exec error"))
			},
			expectedErr: errors.New("exec error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			testCase.mock(mock, testCase.inputQuery, testCase.inputArgs)

			// Repo
			makerRepo := NewTeamRelationMakerPostgres(db)
			err = makerRepo.DeleteByIDs(testCase.inputTeamID, testCase.inputUserID)

			// Assert
			assert.Equal(t, testCase.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
