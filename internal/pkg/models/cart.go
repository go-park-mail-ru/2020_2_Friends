package models

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type CartRequest struct {
	ProductID string `json:"product_id"`
	VendorID  string `json:"vendor_id"`
}

func (c *CartRequest) Sanitize() {
	p := bluemonday.UGCPolicy()
	c.ProductID = p.Sanitize(c.ProductID)
	c.VendorID = p.Sanitize(c.VendorID)
}
