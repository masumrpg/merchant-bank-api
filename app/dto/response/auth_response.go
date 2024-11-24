package response

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User UserResponse `json:"user"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresIn int64        `json:"expires_in"` // in seconds
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}
