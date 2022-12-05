package movie

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// SvcMock describes a mock struct.
type SvcMock struct {
	mock.Mock
}

// GetMoviesPublic mock.
func (m *SvcMock) GetMoviesPublic(ctx context.Context, sortType string) (*GetmoviesRes, error) {
	args := m.MethodCalled("GetMoviesPublic", ctx, sortType)

	return args.Get(0).(*GetmoviesRes), args.Error(1)
}

// GetMovies mock.
func (m *SvcMock) GetMovies(ctx context.Context, sortType string) (*GetmoviesRes, error) {
	args := m.MethodCalled("GetMovies", ctx, sortType)

	return args.Get(0).(*GetmoviesRes), args.Error(1)
}

// GetUserMovies mock.
func (m *SvcMock) GetUserMovies(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error) {
	args := m.MethodCalled("GetUserMovies", ctx, userID, sortType)

	return args.Get(0).(*GetmoviesRes), args.Error(1)
}

// GetUserMoviesPublic mock.
func (m *SvcMock) GetUserMoviesPublic(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error) {
	args := m.MethodCalled("GetUserMoviesPublic", ctx, userID, sortType)

	return args.Get(0).(*GetmoviesRes), args.Error(1)
}

// CreateMovie mock.
func (m *SvcMock) CreateMovie(ctx context.Context, movie NewMovie) error {
	args := m.MethodCalled("CreateMovie", ctx, movie)

	return args.Error(0)
}

// Action mock.
func (m *SvcMock) Action(ctx context.Context, movieID int, action string) error {
	args := m.MethodCalled("Action", ctx, movieID, action)

	return args.Error(0)
}

// RemoveAction mock.
func (m *SvcMock) RemoveAction(ctx context.Context, movieID int, action string) error {
	args := m.MethodCalled("RemoveAction", ctx, movieID, action)

	return args.Error(0)
}
