package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type Review struct {
	UserID    string    `json:"user_id"`
	OrderID   int       `json:"order_id"`
	VendorID  int       `json:"-"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Review) Sanitize() {
	p := bluemonday.UGCPolicy()
	r.Text = p.Sanitize(r.Text)
}
