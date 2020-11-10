package models

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type Vendor struct {
	ID       int       `json:"id"`
	Name     string    `json:"store_name"`
	Products []Product `json:"products,omitempty"`
}

//easyjson:json
type Product struct {
	ID       int    `json:"id"`
	Picture  string `json:"picture"`
	Name     string `json:"food_name"`
	Price    string `json:"food_price"`
	VendorID int    `json:"vendor_id"`
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
	p.Price = pol.Sanitize(p.Price)
}
