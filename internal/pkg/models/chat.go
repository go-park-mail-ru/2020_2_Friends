package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type Message struct {
	Type      string    `json:"type"`
	OrderID   int       `json:"order_id,omitempty"`
	UserID    string    `json:"-"`
	VendorID  int       `json:"vendor_id,omitempty"`
	IsYourMsg bool      `json:"is_your_msg"`
	Text      string    `json:"text"`
	SentAt    time.Time `json:"-"`
	SentAtStr string    `json:"sent_at"`
}

//easyjson:json
type Chat struct {
	OrderID          int    `json:"order_id"`
	InterlocutorID   string `json:"interlocutor_id"`
	InterlocutorName string `json:"interlocutor_name"`
	LastMsg          string `json:"last_message"`
}

//easyjson:json
type VendorChatsWithInfo struct {
	Chats         []Chat `json:"chats"`
	VendorName    string `json:"vendor_name"`
	VendorPicture string `json:"picture"`
}

func (m *Message) Sanitaze() {
	p := bluemonday.UGCPolicy()
	m.Text = p.Sanitize(m.Text)
}
