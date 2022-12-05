package movie_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"movierama/internal/app/movie"
	sqlMovieMock "movierama/internal/infra/repository/sql/movie"
	"testing"
	"time"
)

func Test_GetMoviesPublic(t *testing.T) {
	time1HourAgo := time.Now().Add(time.Duration(-1) * time.Hour)
	tests := map[string]struct {
		sqlRepo  *sqlMovieMock.Mock
		sortType string
		expRes   *movie.GetmoviesRes
		expErr   error
	}{
		"Should get public movies": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetMoviesPublic", context.TODO(), "date").Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("GetMoviesPublic", context.TODO(), "date").Return([]sqlMovieMock.Movie{}, errors.New("random error"))

				return &repo
			}(),
			sortType: "date",
			expErr:   errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			res, err := app.GetMoviesPublic(context.TODO(), tt.sortType)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes, res)
		})
	}
}

func Test_GetMovies(t *testing.T) {
	time1HourAgo := time.Now().Add(time.Duration(-1) * time.Hour)
	tests := map[string]struct {
		sqlRepo  *sqlMovieMock.Mock
		sortType string
		ctx      context.Context
		expRes   *movie.GetmoviesRes
		expErr   error
	}{
		"Should get movies": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						UserLiked:   true,
						UserHated:   false,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 3, "date").
					Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						UserLiked:   true,
						UserHated:   false,
						IsSameUser:  false,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"movie should contain IsSameUser on same user posted movie": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      3,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 3, "date").
					Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      3,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						IsSameUser:  true,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("GetMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 3, "date").
					Return([]sqlMovieMock.Movie{}, errors.New("random error"))

				return &repo
			}(),
			sortType: "date",
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expErr:   errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			res, err := app.GetMovies(tt.ctx, tt.sortType)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes, res)
		})
	}
}

func Test_GetUserMoviesPublic(t *testing.T) {
	time1HourAgo := time.Now().Add(time.Duration(-1) * time.Hour)
	tests := map[string]struct {
		sqlRepo  *sqlMovieMock.Mock
		sortType string
		userID   int
		expRes   *movie.GetmoviesRes
		expErr   error
	}{
		"Should get public movies": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetUserMoviesPublic", context.TODO(), 2, "date").Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			userID:   2,
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("GetUserMoviesPublic", context.TODO(), 2, "date").Return([]sqlMovieMock.Movie{}, errors.New("random error"))

				return &repo
			}(),
			sortType: "date",
			userID:   2,
			expErr:   errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			res, err := app.GetUserMoviesPublic(context.TODO(), tt.userID, tt.sortType)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes, res)
		})
	}
}

func Test_GetUserMovies(t *testing.T) {
	time1HourAgo := time.Now().Add(time.Duration(-1) * time.Hour)
	tests := map[string]struct {
		sqlRepo  *sqlMovieMock.Mock
		sortType string
		ctx      context.Context
		userID   int
		expRes   *movie.GetmoviesRes
		expErr   error
	}{
		"Should get movies": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						UserLiked:   true,
						UserHated:   false,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetUserMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 2, 3, "date").
					Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			userID:   2,
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      2,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						UserLiked:   true,
						UserHated:   false,
						IsSameUser:  false,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"movie should contain IsSameUser on same user posted movie": {
			sqlRepo: func() *sqlMovieMock.Mock {
				movies := []sqlMovieMock.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      3,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						CreatedAt:   time1HourAgo,
					},
				}

				repo := sqlMovieMock.Mock{}
				repo.On("GetUserMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 3, 3, "date").
					Return(movies, nil)

				return &repo
			}(),
			sortType: "date",
			userID:   3,
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expRes: &movie.GetmoviesRes{
				Movies: []movie.Movie{
					{
						ID:          1,
						Title:       "movie title",
						Description: "movie description",
						UserID:      3,
						PostedBy:    "Full Name",
						Likes:       3,
						Hates:       4,
						IsSameUser:  true,
						TimeAgo:     "about an hour ago",
					},
				},
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("GetUserMovies", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 2, 3, "date").
					Return([]sqlMovieMock.Movie{}, errors.New("random error"))

				return &repo
			}(),
			sortType: "date",
			userID:   2,
			ctx:      context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expErr:   errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			res, err := app.GetUserMovies(tt.ctx, tt.userID, tt.sortType)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes, res)
		})
	}
}

func Test_CreateMovie(t *testing.T) {
	tests := map[string]struct {
		sqlRepo *sqlMovieMock.Mock
		ctx     context.Context
		movie   movie.NewMovie
		expErr  error
	}{
		"Should get movies": {
			sqlRepo: func() *sqlMovieMock.Mock {
				m := &sqlMovieMock.SQLMovie{
					Title:       "movie title",
					Description: "movie description",
					UserID:      3,
				}

				repo := sqlMovieMock.Mock{}
				repo.On("CreateMovie", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), m).
					Return(nil)

				return &repo
			}(),
			ctx: context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			movie: movie.NewMovie{
				Title:       "movie title",
				Description: "movie description",
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				m := &sqlMovieMock.SQLMovie{
					Title:       "movie title",
					Description: "movie description",
					UserID:      3,
				}
				repo := sqlMovieMock.Mock{}
				repo.On("CreateMovie", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), m).
					Return(errors.New("random error"))

				return &repo
			}(),
			movie: movie.NewMovie{
				Title:       "movie title",
				Description: "movie description",
			}, ctx: context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expErr: errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			err := app.CreateMovie(tt.ctx, tt.movie)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_Action(t *testing.T) {
	tests := map[string]struct {
		sqlRepo *sqlMovieMock.Mock
		ctx     context.Context
		movieID int
		action  string
		expErr  error
	}{
		"Should add movie action": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("AddMovieAction", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 1, 3, "like").
					Return(nil)

				return &repo
			}(),
			ctx:     context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			movieID: 1,
			action:  "like",
			expErr:  nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("AddMovieAction", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 1, 3, "like").
					Return(errors.New("random error"))

				return &repo
			}(),
			movieID: 1,
			action:  "like",
			ctx:     context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expErr:  errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			err := app.Action(tt.ctx, tt.movieID, tt.action)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_RemoveAction(t *testing.T) {
	tests := map[string]struct {
		sqlRepo *sqlMovieMock.Mock
		ctx     context.Context
		movieID int
		action  string
		expErr  error
	}{
		"Should add movie action": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("RemoveMovieAction", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 1, 3, "like").
					Return(nil)

				return &repo
			}(),
			ctx:     context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			movieID: 1,
			action:  "like",
			expErr:  nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqlMovieMock.Mock {
				repo := sqlMovieMock.Mock{}
				repo.On("RemoveMovieAction", context.WithValue(
					context.TODO(), movie.AuthUserIDContextKey, 3), 1, 3, "like").
					Return(errors.New("random error"))

				return &repo
			}(),
			movieID: 1,
			action:  "like",
			ctx:     context.WithValue(context.TODO(), movie.AuthUserIDContextKey, 3),
			expErr:  errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := movie.NewService(tt.sqlRepo)

			err := app.RemoveAction(tt.ctx, tt.movieID, tt.action)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
