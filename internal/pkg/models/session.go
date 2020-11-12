package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type Session struct {
	Name       string
	UserID     string
	ExpireTime time.Duration
}

func (s *Session) Sanitize() {
	p := bluemonday.UGCPolicy()
	s.Name = p.Sanitize(s.Name)
	s.UserID = p.Sanitize(s.UserID)
}
