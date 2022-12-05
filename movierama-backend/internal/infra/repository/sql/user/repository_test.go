package user_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"movierama/internal/infra/repository/sql/user"
	"testing"
)

type dbMock struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func Test_NewRepository(t *testing.T) {
	tests := map[string]struct {
		reader *sql.DB
		writer *sql.DB
		expErr error
	}{
		"success": {
			reader: &sql.DB{},
			writer: &sql.DB{},
		},
		"missing reader": {
			writer: &sql.DB{},
			expErr: errors.New("db reader is nil"),
		},
		"missing writer": {
			reader: &sql.DB{},
			expErr: errors.New("db writer is nil"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := user.NewRepository(tt.reader, tt.writer)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_GetUserAuthDetails(t *testing.T) {
	cases := map[string]struct {
		dbMock   dbMock
		username string
		expRes   *user.AuthUserDetails
		expErr   error
	}{
		"should return auth user details": {
			username: "test_username",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, "test_username", "pass")

				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE username=?`).
					WithArgs("test_username").
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expRes: &user.AuthUserDetails{
				ID:       1,
				Username: "test_username",
				Password: "pass",
			},
		},
		"should return error on sql error": {
			username: "test_username",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE username=?`).
					WithArgs("test_username").
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
		"should return error when scan.Row() returns error": {
			username: "test_username",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(nil, "test_username", "pass")

				mock.ExpectQuery(`SELECT id, username, password FROM users WHERE username=?`).
					WithArgs("test_username").
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql: Scan error on column index 0, name \"id\": converting NULL to int is unsupported"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := user.NewRepository(tt.dbMock.db, tt.dbMock.db)
			resp, err := repo.GetUserAuthDetails(context.TODO(), tt.username)
			assert.Equal(t, tt.expRes, resp)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr.Error(), err.Error())
			}
		})
	}
}

func Test_CreateUser(t *testing.T) {
	cases := map[string]struct {
		dbMock dbMock
		user   *user.SQLUser
		expErr error
	}{
		"should create user": {
			user: &user.SQLUser{
				Username:  "username",
				Password:  "pass",
				FirstName: "first name",
				LastName:  "last name",
			},
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(`INSERT INTO users (
		username,
		password,
		first_name,
		last_name
	) VALUES(
		?,
		?,
		?,
		?
	);`).
					WithArgs([]driver.Value{"username", "pass", "first name", "last name"}...).
					WillReturnResult(sqlmock.NewResult(1, 1))
				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
		},
		"should return error on exec error": {
			user: &user.SQLUser{
				Username:  "username",
				Password:  "pass",
				FirstName: "first name",
				LastName:  "last name",
			}, dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(`INSERT INTO users (
		username,
		password,
		first_name,
		last_name
	) VALUES(
		?,
		?,
		?,
		?
	);`).
					WithArgs([]driver.Value{"username", "pass", "first name", "last name"}...).
					WillReturnError(errors.New("exec error"))
				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("exec error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := user.NewRepository(tt.dbMock.db, tt.dbMock.db)
			err := repo.CreateUser(context.TODO(), tt.user)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
