package user

// AuthUserDetails struct.
type AuthUserDetails struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

// SQLUser struct.
type SQLUser struct {
	ID        *int   `db:"id"`
	Username  string `db:"username"`
	Password  string `db:"password"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
