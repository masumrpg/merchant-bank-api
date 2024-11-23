package models

type Customer struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
	Balance  int64  `json:"balance"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}
