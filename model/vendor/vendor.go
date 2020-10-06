package vendor

//easyjson:json
type Vendor struct {
	Name     string    `json:"storeName"`
	Products []Product `json:"products"`
}

//easyjson:json
type Product struct {
	PicturePath string `json:"picturePath"`
	Name        string `json:"foodName"`
	Price       string `json:"foodPrice"`
}
