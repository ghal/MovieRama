package auth

// LoginRes contains the login response.
type LoginRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
