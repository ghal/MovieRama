package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Ordering types.
const (
	SortCaseDate           = "date"
	SortCaseLikes          = "likes"
	SortCaseHates          = "hates"
	TypeOrderByCreatedDate = "created_at"
	TypeOrderByLikes       = "likes"
	TypeOrderByHates       = "hates"
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

// GetMoviesPublic returns a list of movies without auth.
func (sr *Repository) GetMoviesPublic(ctx context.Context, sortType string) ([]Movie, error) {
	sqlQuery := fmt.Sprintf(`SELECT 
    movie.id,
    movie.title,
    movie.description,
    movie.user_id,
    movie.created_at,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='like') AS likes,
    (SELECT COUNT(*) FROM movies_users_actions WHERE movie_id=movie.id AND action='hate') AS hates,
    (SELECT CONCAT_WS(' ', first_name, last_name) FROM users WHERE id=movie.user_id) AS posted_by
    	FROM movies AS movie 
    		ORDER BY %s DESC;`, ConvertSortTypeToOrderByColumn(sortType))

	rows, err := sr.read.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []Movie
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.UserID,
			&movie.CreatedAt,
			&movie.Likes,
			&movie.Hates,
			&movie.PostedBy,
		); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return movies, err
	}

	return movies, nil
}

// GetMovies returns a list of movies.
func (sr *Repository) GetMovies(ctx context.Context, authUsrID int, sortType string) ([]Movie, error) {
	sqlQuery := fmt.Sprintf(`SELECT 
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
    		ORDER BY %s DESC;`, ConvertSortTypeToOrderByColumn(sortType))

	rows, err := sr.read.QueryContext(ctx, sqlQuery, authUsrID, authUsrID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []Movie
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.UserID,
			&movie.CreatedAt,
			&movie.Likes,
			&movie.Hates,
			&movie.UserLiked,
			&movie.UserHated,
			&movie.PostedBy,
		); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return movies, err
	}

	return movies, nil
}

// GetUserMovies returns a list of movies for a particular user.
func (sr *Repository) GetUserMovies(ctx context.Context, userID, authUsrID int, sortType string) ([]Movie, error) {
	sqlQuery := fmt.Sprintf(`SELECT 
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
    		ORDER BY %s DESC;`, ConvertSortTypeToOrderByColumn(sortType))

	rows, err := sr.read.QueryContext(ctx, sqlQuery, authUsrID, authUsrID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []Movie
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.UserID,
			&movie.CreatedAt,
			&movie.Likes,
			&movie.Hates,
			&movie.UserLiked,
			&movie.UserHated,
			&movie.PostedBy,
		); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return movies, err
	}

	return movies, nil
}

// GetUserMoviesPublic returns a list of movies for a particular user.
func (sr *Repository) GetUserMoviesPublic(ctx context.Context, userID int, sortType string) ([]Movie, error) {
	sqlQuery := fmt.Sprintf(`SELECT 
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
    		ORDER BY %s DESC;`, ConvertSortTypeToOrderByColumn(sortType))

	rows, err := sr.read.QueryContext(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []Movie
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.UserID,
			&movie.CreatedAt,
			&movie.Likes,
			&movie.Hates,
			&movie.PostedBy,
		); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return movies, err
	}

	return movies, nil
}

// CreateMovie creates a new movie.
func (sr *Repository) CreateMovie(ctx context.Context, movie *SQLMovie) error {
	sqlQuery := `INSERT INTO movies (
		title,
		user_id,
		description
	) VALUES(
		?,
		?,
		?
	);`

	_, err := sr.write.ExecContext(ctx,
		sqlQuery,
		movie.Title,
		movie.UserID,
		movie.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

// AddMovieAction adds an action for a movie.
func (sr *Repository) AddMovieAction(ctx context.Context, movieID, userID int, action string) error {
	sqlQuery := `INSERT INTO movies_users_actions (
		movie_id,
		user_id,
		action
	) VALUES(
		?,
		?,
		?
	);`

	_, err := sr.write.ExecContext(ctx,
		sqlQuery,
		movieID,
		userID,
		action,
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveMovieAction removes an action for a movie.
func (sr *Repository) RemoveMovieAction(ctx context.Context, movieID, userID int, action string) error {
	sqlQuery := `DELETE FROM movies_users_actions WHERE movie_id=? AND user_id=? AND action=?;`

	_, err := sr.write.ExecContext(ctx,
		sqlQuery,
		movieID,
		userID,
		action,
	)
	if err != nil {
		return err
	}

	return nil
}

// ConvertSortTypeToOrderByColumn converts sort type to db column name type.
func ConvertSortTypeToOrderByColumn(sortType string) string {
	var orderSQL string
	switch sortType {
	case SortCaseDate:
		orderSQL = TypeOrderByCreatedDate
		break
	case SortCaseLikes:
		orderSQL = TypeOrderByLikes
		break
	case SortCaseHates:
		orderSQL = TypeOrderByHates
		break
	default:
		orderSQL = TypeOrderByCreatedDate
		break
	}

	return orderSQL
}
