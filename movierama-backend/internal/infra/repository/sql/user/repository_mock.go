package user

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Mock describes a mock struct.
type Mock struct {
	mock.Mock
}

// GetUserAuthDetails mock.
func (m *Mock) GetUserAuthDetails(ctx context.Context, username string) (*AuthUserDetails, error) {
	args := m.MethodCalled("GetUserAuthDetails", ctx, username)

	return args.Get(0).(*AuthUserDetails), args.Error(1)
}

// CreateUser mock.
func (m *Mock) CreateUser(ctx context.Context, user *SQLUser) error {
	args := m.MethodCalled("CreateUser", ctx, user)

	return args.Error(0)
}
