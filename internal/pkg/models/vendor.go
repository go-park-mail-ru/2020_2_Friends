package models

//easyjson:json
type Vendor struct {
	ID       int       `json:"id"`
	Name     string    `json:"storeName"`
	Products []Product `json:"products,omitempty"`
}

//easyjson:json
type Product struct {
	ID       int    `json:"id"`
	Picture  string `json:"picture"`
	Name     string `json:"foodName"`
	Price    string `json:"foodPrice"`
	VendorID int    `json:"vendorID"`
}
