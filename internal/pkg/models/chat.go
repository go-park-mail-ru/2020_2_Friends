package models

import "time"

type Message struct {
	OrderID   int       `json:"order_id,omitempty"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	SentAt    time.Time `json:"-"`
	SentAtStr string    `json:"sent_at"`
}

type Chat struct {
	OrderID          int    `json:"order_id"`
	InterlocutorID   string `json:"interlocutor_id"`
	InterlocutorName string `json:"interlocutor_name"`
	LastMsg          string `json:"last_message"`
}
