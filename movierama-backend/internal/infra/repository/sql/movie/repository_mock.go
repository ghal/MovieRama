package movie

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Mock describes a mock struct.
type Mock struct {
	mock.Mock
}

// GetMoviesPublic mock.
func (m *Mock) GetMoviesPublic(ctx context.Context, sortType string) ([]Movie, error) {
	args := m.MethodCalled("GetMoviesPublic", ctx, sortType)

	return args.Get(0).([]Movie), args.Error(1)
}

// GetMovies mock.
func (m *Mock) GetMovies(ctx context.Context, authUsrID int, sortType string) ([]Movie, error) {
	args := m.MethodCalled("GetMovies", ctx, authUsrID, sortType)

	return args.Get(0).([]Movie), args.Error(1)
}

// GetUserMovies mock.
func (m *Mock) GetUserMovies(ctx context.Context, userID, authUsrID int, sortType string) ([]Movie, error) {
	args := m.MethodCalled("GetUserMovies", ctx, userID, authUsrID, sortType)

	return args.Get(0).([]Movie), args.Error(1)
}

// GetUserMoviesPublic mock.
func (m *Mock) GetUserMoviesPublic(ctx context.Context, userID int, sortType string) ([]Movie, error) {
	args := m.MethodCalled("GetUserMoviesPublic", ctx, userID, sortType)

	return args.Get(0).([]Movie), args.Error(1)
}

// CreateMovie mock.
func (m *Mock) CreateMovie(ctx context.Context, movie *SQLMovie) error {
	args := m.MethodCalled("CreateMovie", ctx, movie)

	return args.Error(0)
}

// AddMovieAction mock.
func (m *Mock) AddMovieAction(ctx context.Context, movieID, userID int, action string) error {
	args := m.MethodCalled("AddMovieAction", ctx, movieID, userID, action)

	return args.Error(0)
}

// RemoveMovieAction mock.
func (m *Mock) RemoveMovieAction(ctx context.Context, movieID, userID int, action string) error {
	args := m.MethodCalled("RemoveMovieAction", ctx, movieID, userID, action)

	return args.Error(0)
}
