package movie

// NewMovie contains the newly created movie payload struct.
type NewMovie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
