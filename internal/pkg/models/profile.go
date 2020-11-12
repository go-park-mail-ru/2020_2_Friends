package models

import "github.com/microcosm-cc/bluemonday"

//easyjson:json
type Profile struct {
	UserID    string   `json:"userId"`
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	Addresses []string `json:"addresses"`
	Points    int      `json:"points"`
	Avatar    string   `json:"avatar"`
}

//easyjson:json
type ImgResponse struct {
	Avatar string `json:"avatar"`
}

func (p *Profile) Sanitize() {
	pol := bluemonday.UGCPolicy()
	p.UserID = pol.Sanitize(p.UserID)
	p.Name = pol.Sanitize(p.Name)
	p.Phone = pol.Sanitize(p.Phone)
	p.Avatar = pol.Sanitize(p.Avatar)
	for i, address := range p.Addresses {
		p.Addresses[i] = pol.Sanitize(address)
	}
}
