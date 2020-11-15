package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type OrderRequest struct {
	VendorID   int       `json:"-"`
	VendorName string    `json:"-"`
	Products   []int     `json:"products"`
	CreatedAt  time.Time `json:"created_at"`
	Address    string    `json:"address"`
}

//easyjson:json
type OrderResponse struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	VendorName string         `json:"vendor_name"`
	Products   []OrderProduct `json:"products"`
	CreatedAt  time.Time      `json:"created_at"`
	Address    string         `json:"address"`
	Status     string         `json:"status"`
	ProductIDs []string       `json:"-"`
}

//easyjson:json
type OrderProduct struct {
	ID      int    `json:"id"`
	Picture string `json:"picture"`
	Name    string `json:"food_name"`
	Price   int    `json:"food_price"`
}

//easyjson:json
type IDRequest struct {
	ID int `json:"id"`
}

//easyjson:json
type IDResponse struct {
	ID int `json:"id"`
}

//easyjson:json
type OrderStatusRequest struct {
	Status string `json:"status"`
}

func (o *OrderRequest) Sanitize() {
	p := bluemonday.UGCPolicy()
	o.VendorName = p.Sanitize(o.VendorName)
	o.Address = p.Sanitize(o.Address)
}

func (s *OrderStatusRequest) Sanitize() {
	p := bluemonday.UGCPolicy()
	s.Status = p.Sanitize(s.Status)
}
