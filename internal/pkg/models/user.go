package models

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role"`
}

func (u *User) Sanitize() {
	p := bluemonday.UGCPolicy()
	u.ID = p.Sanitize(u.ID)
	u.Login = p.Sanitize(u.Login)
	u.Password = p.Sanitize(u.Password)
}
