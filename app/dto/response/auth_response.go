package response

// RegisterResponse represents the registration response
type RegisterResponse struct {
	Customer CustomerResponse `json:"customer"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string           `json:"token"`
	Customer  CustomerResponse `json:"customer"`
	ExpiresIn int64            `json:"expires_in"` // in seconds
}

type CustomerResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}
