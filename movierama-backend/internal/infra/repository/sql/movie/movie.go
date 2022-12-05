package movie

import "time"

// Movie struct.
type Movie struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	PostedBy    string    `db:"posted_by"`
	Likes       int       `db:"likes"`
	Hates       int       `db:"hates"`
	UserLiked   bool      `db:"usr_liked"`
	UserHated   bool      `db:"usr_hated"`
	CreatedAt   time.Time `db:"created_at"`
}

// SQLMovie struct.
type SQLMovie struct {
	ID          *int      `db:"id"`
	UserID      int       `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
