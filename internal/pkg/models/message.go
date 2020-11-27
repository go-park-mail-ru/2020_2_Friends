package models

import "time"

type Message struct {
	OrderID int       `json:"order_id"`
	UserID  string    `json:"user_id"`
	Text    string    `json:"text"`
	SentAt  time.Time `json:"sent_at"`
}
