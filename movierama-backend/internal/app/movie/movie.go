package movie

// Movie contains the movie data.
type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	PostedBy    string `json:"posted_by"`
	Likes       int    `json:"likes"`
	Hates       int    `json:"hates"`
	UserLiked   bool   `json:"user_liked"`
	UserHated   bool   `json:"user_hated"`
	IsSameUser  bool   `json:"is_same_user"`
	TimeAgo     string `json:"time_ago"`
}

// NewMovie contains the new movie data.
type NewMovie struct {
	Title       string
	Description string
}

// GetmoviesRes contains the response of get movies.
type GetmoviesRes struct {
	Movies []Movie `json:"movies"`
}
