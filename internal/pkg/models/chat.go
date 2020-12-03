package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type Message struct {
	OrderID   int       `json:"order_id,omitempty"`
	UserID    string    `json:"-"`
	IsYourMsg bool      `json:"is_your_msg"`
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

func (m *Message) Sanitaze() {
	p := bluemonday.UGCPolicy()
	m.Text = p.Sanitize(m.Text)
}
