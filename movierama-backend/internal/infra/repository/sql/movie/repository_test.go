package movie_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"movierama/internal/infra/repository/sql/movie"
	"testing"
	"time"
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
			got, err := movie.NewRepository(tt.reader, tt.writer)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_GetMoviesPublic(t *testing.T) {
	nowTime := time.Now()
	cases := map[string]struct {
		dbMock   dbMock
		sortType string
		expRes   []movie.Movie
		expErr   error
	}{
		"should return public movies": {
			sortType: "date",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "created_at", "likes", "hates", "posted_by"}).
					AddRow(1, "movie title", "movie description", 4, nowTime, 2, 3, "user 1")

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT CONCAT_WS(' ', first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date"))).
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expRes: []movie.Movie{
				{
					ID:          1,
					UserID:      4,
					Title:       "movie title",
					Description: "movie description",
					PostedBy:    "user 1",
					Likes:       2,
					Hates:       3,
					CreatedAt:   nowTime,
				},
			},
		},
		"should return error on sql error": {
			sortType: "date",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT CONCAT_WS(' ', first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date"))).
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			resp, err := repo.GetMoviesPublic(context.TODO(), tt.sortType)
			assert.Equal(t, tt.expRes, resp)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr.Error(), err.Error())
			}
		})
	}
}

func Test_GetMovies(t *testing.T) {
	nowTime := time.Now()
	cases := map[string]struct {
		dbMock     dbMock
		sortType   string
		authUserID int
		expRes     []movie.Movie
		expErr     error
	}{
		"should return movies": {
			sortType:   "date",
			authUserID: 6,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "created_at", "likes", "hates", "usr_liked", "usr_hated", "posted_by"}).
					AddRow(1, "movie title", "movie description", 4, nowTime, 2, 3, 1, 0, "user 1")

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='like' AND a.user_id=?) AS usr_liked,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='hate' AND a.user_id=?) AS usr_hated,
    (SELECT CONCAT_WS(' ', first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date")),
				).
					WithArgs(6, 6).
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expRes: []movie.Movie{
				{
					ID:          1,
					UserID:      4,
					Title:       "movie title",
					Description: "movie description",
					PostedBy:    "user 1",
					Likes:       2,
					Hates:       3,
					UserLiked:   true,
					UserHated:   false,
					CreatedAt:   nowTime,
				},
			},
		},
		"should return error on sql error": {
			authUserID: 6,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='like' AND a.user_id=?) AS usr_liked,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='hate' AND a.user_id=?) AS usr_hated,
    (SELECT CONCAT_WS(' ', first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date")),
				).
					WithArgs(6, 6).
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			resp, err := repo.GetMovies(context.TODO(), tt.authUserID, tt.sortType)
			assert.Equal(t, tt.expRes, resp)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr.Error(), err.Error())
			}
		})
	}
}

func Test_GetUserMoviesPublic(t *testing.T) {
	nowTime := time.Now()
	cases := map[string]struct {
		dbMock   dbMock
		sortType string
		userID   int
		expRes   []movie.Movie
		expErr   error
	}{
		"should return public user movies": {
			sortType: "date",
			userID:   8,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "created_at", "likes", "hates", "posted_by"}).
					AddRow(1, "movie title", "movie description", 4, nowTime, 2, 3, "user 1")

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT CONCAT_WS(" ", first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    	WHERE user_id=?
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date"))).
					WithArgs(8).
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expRes: []movie.Movie{
				{
					ID:          1,
					UserID:      4,
					Title:       "movie title",
					Description: "movie description",
					PostedBy:    "user 1",
					Likes:       2,
					Hates:       3,
					CreatedAt:   nowTime,
				},
			},
		},
		"should return error on sql error": {
			sortType: "likes",
			userID:   8,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT CONCAT_WS(" ", first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    	WHERE user_id=?
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("likes"))).
					WithArgs(8).
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			resp, err := repo.GetUserMoviesPublic(context.TODO(), tt.userID, tt.sortType)
			assert.Equal(t, tt.expRes, resp)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr.Error(), err.Error())
			}
		})
	}
}

func Test_GetUserMovies(t *testing.T) {
	nowTime := time.Now()
	cases := map[string]struct {
		dbMock     dbMock
		sortType   string
		userID     int
		authUserID int
		expRes     []movie.Movie
		expErr     error
	}{
		"should return movies": {
			sortType:   "hates",
			userID:     9,
			authUserID: 6,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "created_at", "likes", "hates", "usr_liked", "usr_hated", "posted_by"}).
					AddRow(1, "movie title", "movie description", 4, nowTime, 2, 3, 1, 0, "user 1")

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='like' AND a.user_id=?) AS usr_liked,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='hate' AND a.user_id=?) AS usr_hated,
    (SELECT CONCAT_WS(" ", first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    	WHERE user_id=?
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("hates")),
				).
					WithArgs(6, 6, 9).
					WillReturnRows(rows)

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expRes: []movie.Movie{
				{
					ID:          1,
					UserID:      4,
					Title:       "movie title",
					Description: "movie description",
					PostedBy:    "user 1",
					Likes:       2,
					Hates:       3,
					UserLiked:   true,
					UserHated:   false,
					CreatedAt:   nowTime,
				},
			},
		},
		"should return error on sql error": {
			userID:     9,
			authUserID: 6,
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectQuery(fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='like' AND a.user_id=?) AS usr_liked,
    (SELECT COUNT(*) FROM movies_users_actions a WHERE movie_id=movie.id AND action='hate' AND a.user_id=?) AS usr_hated,
    (SELECT CONCAT_WS(" ", first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    	WHERE user_id=?
    		ORDER BY %s DESC;`, movie.ConvertSortTypeToOrderByColumn("date")),
				).
					WithArgs(6, 6, 9).
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			resp, err := repo.GetUserMovies(context.TODO(), tt.userID, tt.authUserID, tt.sortType)
			assert.Equal(t, tt.expRes, resp)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr.Error(), err.Error())
			}
		})
	}
}

func Test_CreateMovie(t *testing.T) {
	cases := map[string]struct {
		dbMock dbMock
		movie  *movie.SQLMovie
		expErr error
	}{
		"should create movie": {
			movie: &movie.SQLMovie{
				UserID:      6,
				Title:       "movie title",
				Description: "movie description",
			},
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(`INSERT INTO movies (
		title,
		user_id,
		description
	) VALUES(
		?,
		?,
		?
	);`,
				).
					WithArgs("movie title", 6, "movie description").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
		},
		"should return error on sql error": {
			movie: &movie.SQLMovie{
				UserID:      6,
				Title:       "movie title",
				Description: "movie description",
			},
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectExec(`INSERT INTO movies (
		title,
		user_id,
		description
	) VALUES(
		?,
		?,
		?
	);`,
				).
					WithArgs("movie title", 6, "movie description").
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			err := repo.CreateMovie(context.TODO(), tt.movie)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_AddMovieAction(t *testing.T) {
	cases := map[string]struct {
		dbMock          dbMock
		movieID, userID int
		action          string
		expErr          error
	}{
		"should add movie action": {
			movieID: 1,
			userID:  2,
			action:  "like",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(`INSERT INTO movies_users_actions (
		movie_id,
		user_id,
		action
	) VALUES(
		?,
		?,
		?
	);`,
				).
					WithArgs(1, 2, "like").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
		},
		"should return error on sql error": {
			movieID: 1,
			userID:  2,
			action:  "like",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectExec(`INSERT INTO movies_users_actions (
		movie_id,
		user_id,
		action
	) VALUES(
		?,
		?,
		?
	);`,
				).
					WithArgs(1, 2, "like").
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			err := repo.AddMovieAction(context.TODO(), tt.movieID, tt.userID, tt.action)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_RemoveMovieAction(t *testing.T) {
	cases := map[string]struct {
		dbMock          dbMock
		movieID, userID int
		action          string
		expErr          error
	}{
		"should remove movie action": {
			movieID: 1,
			userID:  2,
			action:  "like",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(`DELETE FROM movies_users_actions WHERE movie_id=? AND user_id=? AND action=?;`).
					WithArgs(1, 2, "like").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
		},
		"should return error on sql error": {
			movieID: 1,
			userID:  2,
			action:  "like",
			dbMock: func() dbMock {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				mock.ExpectExec(`DELETE FROM movies_users_actions WHERE movie_id=? AND user_id=? AND action=?;`).
					WithArgs(1, 2, "like").
					WillReturnError(errors.New("sql error"))

				return dbMock{
					db:   db,
					mock: mock,
				}
			}(),
			expErr: errors.New("sql error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			repo, _ := movie.NewRepository(tt.dbMock.db, tt.dbMock.db)
			err := repo.RemoveMovieAction(context.TODO(), tt.movieID, tt.userID, tt.action)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
