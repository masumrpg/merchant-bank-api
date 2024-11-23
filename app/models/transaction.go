package models

import "time"

type Transaction struct {
	IDTransaction string    `json:"id_transaction"`
	Type          string    `json:"type"`
	FromUser      string    `json:"from_user,omitempty"`
	ToUser        string    `json:"to_user,omitempty"`
	Amount        int64     `json:"amount,omitempty"`
	Details       string    `json:"details,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}
