package models

//easyjson:json
type User struct {
	ID       string   `json:"id"`
	Login    string   `json:"login"`
	Password string   `json:"password,omitempty"`
	Name     string   `json:"name"`
	Points   string   `json:"points"`
	Email    string   `json:"email"`
	Adresses []string `json:"adresses"`
}
