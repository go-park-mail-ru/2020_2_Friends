package user

//easyjson:json
type User struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Points   string   `json:"points"`
	Email    string   `json:"email"`
	Number   string   `json:"number"`
	Adresses []string `json:"adresses"`
}
