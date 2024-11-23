package models

import "time"

type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Token     string    `json:"token,omitempty"`
	ExpiresIn time.Time `json:"expiresIn,omitempty"`
}
