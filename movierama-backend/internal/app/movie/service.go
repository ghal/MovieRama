package movie

import (
	"context"
	sqlmovie "movierama/internal/infra/repository/sql/movie"

	"github.com/xeonx/timeago"
)

// AuthUserIDContextKey contains the context key for the user id.
const AuthUserIDContextKey = UserIDContextKey("auth_user_id")

// UserIDContextKey contains the type of user id context key.
type UserIDContextKey string

// Service handles movie information.
type Service interface {
	GetMoviesPublic(ctx context.Context, sortType string) (*GetmoviesRes, error)
	GetMovies(ctx context.Context, sortType string) (*GetmoviesRes, error)
	GetUserMovies(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error)
	GetUserMoviesPublic(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error)
	CreateMovie(ctx context.Context, movie NewMovie) error
	Action(ctx context.Context, movieID int, action string) error
	RemoveAction(ctx context.Context, movieID int, action string) error
}

type movieService struct {
	mr Repository
}

// NewService constructor.
func NewService(movieRepo Repository) Service {
	return &movieService{
		mr: movieRepo,
	}
}

// GetMoviesPublic function.
func (a movieService) GetMoviesPublic(ctx context.Context, sortType string) (*GetmoviesRes, error) {

	movies, err := a.mr.GetMoviesPublic(ctx, sortType)
	if err != nil {
		return nil, err
	}

	res := &GetmoviesRes{}
	for _, movie := range movies {
		m := Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			UserID:      movie.UserID,
			PostedBy:    movie.PostedBy,
			Likes:       movie.Likes,
			Hates:       movie.Hates,
			UserLiked:   movie.UserLiked,
			UserHated:   movie.UserHated,
			TimeAgo:     timeago.English.Format(movie.CreatedAt),
		}
		res.Movies = append(res.Movies, m)
	}

	return res, err
}

// GetMovies function.
func (a movieService) GetMovies(ctx context.Context, sortType string) (*GetmoviesRes, error) {
	authUserID := ctx.Value(AuthUserIDContextKey).(int)

	movies, err := a.mr.GetMovies(ctx, authUserID, sortType)
	if err != nil {
		return nil, err
	}

	res := &GetmoviesRes{}
	for _, movie := range movies {
		m := Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			UserID:      movie.UserID,
			PostedBy:    movie.PostedBy,
			Likes:       movie.Likes,
			Hates:       movie.Hates,
			UserLiked:   movie.UserLiked,
			UserHated:   movie.UserHated,
			IsSameUser:  movie.UserID == authUserID,
			TimeAgo:     timeago.English.Format(movie.CreatedAt),
		}
		res.Movies = append(res.Movies, m)
	}

	return res, err
}

// GetUserMovies function.
func (a movieService) GetUserMovies(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error) {
	authUserID := ctx.Value(AuthUserIDContextKey).(int)
	movies, err := a.mr.GetUserMovies(ctx, userID, authUserID, sortType)
	if err != nil {
		return nil, err
	}

	res := &GetmoviesRes{}
	for _, movie := range movies {
		m := Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			UserID:      movie.UserID,
			PostedBy:    movie.PostedBy,
			Likes:       movie.Likes,
			Hates:       movie.Hates,
			UserLiked:   movie.UserLiked,
			UserHated:   movie.UserHated,
			IsSameUser:  movie.UserID == authUserID,
			TimeAgo:     timeago.English.Format(movie.CreatedAt),
		}
		res.Movies = append(res.Movies, m)
	}

	return res, err
}

// GetUserMoviesPublic function.
func (a movieService) GetUserMoviesPublic(ctx context.Context, userID int, sortType string) (*GetmoviesRes, error) {
	movies, err := a.mr.GetUserMoviesPublic(ctx, userID, sortType)
	if err != nil {
		return nil, err
	}

	res := &GetmoviesRes{}
	for _, movie := range movies {
		m := Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			UserID:      movie.UserID,
			PostedBy:    movie.PostedBy,
			Likes:       movie.Likes,
			Hates:       movie.Hates,
			UserLiked:   movie.UserLiked,
			UserHated:   movie.UserHated,
			TimeAgo:     timeago.English.Format(movie.CreatedAt),
		}
		res.Movies = append(res.Movies, m)
	}

	return res, err
}

// CreateMovie function.
func (a movieService) CreateMovie(ctx context.Context, movie NewMovie) error {
	authUserID := ctx.Value(AuthUserIDContextKey).(int)
	println(ctx)
	m := &sqlmovie.SQLMovie{
		UserID:      authUserID,
		Title:       movie.Title,
		Description: movie.Description,
	}
	err := a.mr.CreateMovie(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

// Action function.
func (a movieService) Action(ctx context.Context, movieID int, action string) error {
	authUserID := ctx.Value(AuthUserIDContextKey).(int)

	err := a.mr.AddMovieAction(ctx, movieID, authUserID, action)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAction function.
func (a movieService) RemoveAction(ctx context.Context, movieID int, action string) error {
	authUserID := ctx.Value(AuthUserIDContextKey).(int)
	err := a.mr.RemoveMovieAction(ctx, movieID, authUserID, action)
	if err != nil {
		return err
	}

	return nil
}
