package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type OrderRequest struct {
	VendorID   int            `json:"-"`
	VendorName string         `json:"-"`
	ProductIDs []int          `json:"products"`
	Products   []OrderProduct `json:"-"`
	CreatedAt  time.Time      `json:"-"`
	Address    string         `json:"address"`
	Price      int            `json:"-"`
}

//easyjson:json
type OrderResponse struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	VendorName string         `json:"vendor_name,omitempty"`
	Products   []OrderProduct `json:"products"`
	CreatedAt  time.Time      `json:"created_at"`
	Address    string         `json:"address"`
	Status     string         `json:"status"`
	Price      int            `json:"price"`
}

type VendorOrdersResponse struct {
	VendorName    string          `json:"vendor_name"`
	VendorPicture string          `json:"picture"`
	Orders        []OrderResponse `json:"orders"`
}

//easyjson:json
type OrderProduct struct {
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
