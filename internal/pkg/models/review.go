package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type Review struct {
	UserID    string    `json:"-"`
	Username  string    `json:"username"`
	OrderID   int       `json:"order_id"`
	VendorID  int       `json:"-"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

//easyjson:json
type VendorReviewsResponse struct {
	VendorName    string   `json:"vendor_name"`
	VendorPicture string   `json:"vendor_picture"`
	Reviews       []Review `json:"reviews"`
}

func (r *Review) Sanitize() {
	p := bluemonday.UGCPolicy()
	r.Text = p.Sanitize(r.Text)
}
