package auth

import (
	"context"
	"movierama/internal/infra/repository/sql/user"
)

// Repository should be able to manage the user auth.
type Repository interface {
	GetUserAuthDetails(ctx context.Context, username string) (*user.AuthUserDetails, error)
	CreateUser(ctx context.Context, user *user.SQLUser) error
}
