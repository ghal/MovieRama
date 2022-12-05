package auth

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// SvcMock describes a mock struct.
type SvcMock struct {
	mock.Mock
}

// Login mock.
func (sm SvcMock) Login(c echo.Context, ul *UserLogin) (*LoginRes, error) {
	args := sm.MethodCalled("Login", c, ul)

	return args.Get(0).(*LoginRes), args.Error(1)
}

// Register mock.
func (sm SvcMock) Register(ctx context.Context, registerUser *UserRegister) error {
	args := sm.MethodCalled("Register", ctx, registerUser)

	return args.Error(0)
}
