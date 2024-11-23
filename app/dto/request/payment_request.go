package request

var PaymentRequest struct {
	Amount     int64  `json:"amount"`
	ToUsername string `json:"to_username"`
	Details    string `json:"details"`
}
