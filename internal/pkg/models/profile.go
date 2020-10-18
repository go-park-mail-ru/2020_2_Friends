package models

//easyjson:json
type Profile struct {
	UserID    string   `json:"userId"`
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	Addresses []string `json:"adresses"`
	Points    int      `json:"points"`
}
