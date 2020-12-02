package models

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type Vendor struct {
	ID          int       `json:"id"`
	Name        string    `json:"store_name"`
	HintContent string    `json:"hintContent"`
	Products    []Product `json:"products"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
	Longitude   float32   `json:"longitude"`
	Latitude    float32   `json:"latitude"`
	Radius      int       `json:"distance"`
}

//easyjson:json
type Product struct {
	ID          int    `json:"id"`
	Picture     string `json:"picture"`
	Name        string `json:"food_name"`
	Description string `json:"description"`
	Price       int    `json:"food_price"`
	VendorID    int    `json:"vendor_id"`
}

//easyjson:json
type AddResponse struct {
	ID int `json:"id"`
}

func NewEmptyVendor() Vendor {
	return Vendor{
		Products: make([]Product, 0),
	}
}

func (v *Vendor) Sanitize() {
	p := bluemonday.UGCPolicy()
	v.Name = p.Sanitize(v.Name)
	for idx := range v.Products {
		v.Products[idx].Sanitize()
	}
}

func (p *Product) Sanitize() {
	pol := bluemonday.UGCPolicy()
	p.Picture = pol.Sanitize(p.Picture)
	p.Name = pol.Sanitize(p.Name)
	p.Description = pol.Sanitize(p.Description)
}
