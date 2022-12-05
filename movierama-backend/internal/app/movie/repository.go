package movie

import (
	"context"
	moviesql "movierama/internal/infra/repository/sql/movie"
)

// Repository should be able to manage the movies.
type Repository interface {
	GetMoviesPublic(ctx context.Context, sortType string) ([]moviesql.Movie, error)
	GetMovies(ctx context.Context, authUsrID int, sortType string) ([]moviesql.Movie, error)
	GetUserMovies(ctx context.Context, userID, authUsrID int, sortType string) ([]moviesql.Movie, error)
	GetUserMoviesPublic(ctx context.Context, userID int, sortType string) ([]moviesql.Movie, error)
	CreateMovie(ctx context.Context, movie *moviesql.SQLMovie) error
	AddMovieAction(ctx context.Context, movieID, userID int, action string) error
	RemoveMovieAction(ctx context.Context, movieID, userID int, action string) error
}
