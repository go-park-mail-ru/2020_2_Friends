package models

//easyjson:json
type Vendor struct {
	ID          int       `json:"id"`
	Name        string    `json:"store_name"`
	Products    []Product `json:"products,omitempty"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
}

//easyjson:json
type Product struct {
	ID       int    `json:"id"`
	Picture  string `json:"picture"`
	Name     string `json:"food_name"`
	Price    string `json:"food_price"`
	VendorID int    `json:"vendor_id"`
}

type AddResponse struct {
	ID int `json:"id"`
}
