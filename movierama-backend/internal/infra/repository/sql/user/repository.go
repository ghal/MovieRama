package user

import (
	"context"
	"database/sql"
	"errors"
)

// Repository definition.
type Repository struct {
	read  *sql.DB
	write *sql.DB
}

// NewRepository constructor.
func NewRepository(reader *sql.DB, writer *sql.DB) (*Repository, error) {
	if reader == nil {
		return nil, errors.New("db reader is nil")
	}
	if writer == nil {
		return nil, errors.New("db writer is nil")
	}
	return &Repository{read: reader, write: writer}, nil
}

// GetUserAuthDetails returns a user.
func (sr *Repository) GetUserAuthDetails(ctx context.Context, username string) (*AuthUserDetails, error) {
	sqlQuery := `SELECT id, username, password FROM users WHERE username=?`

	row := sr.read.QueryRowContext(ctx, sqlQuery, username)

	var user AuthUserDetails
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user.
func (sr *Repository) CreateUser(ctx context.Context, user *SQLUser) error {
	sqlQuery := `INSERT INTO users (
		username,
		password,
		first_name,
		last_name
	) VALUES(
		?,
		?,
		?,
		?
	);`

	var err error
	_, err = sr.write.ExecContext(ctx, sqlQuery,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return err
	}

	return nil
}
