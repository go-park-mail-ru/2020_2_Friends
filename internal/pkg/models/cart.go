package models

//easyjson:json
type CartRequest struct {
	ProductID string `json:"product_id"`
	VendorID  string `json:"vendor_id"`
}
