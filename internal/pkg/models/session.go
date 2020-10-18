package models

import "time"

type Session struct {
	Name       string
	UserID     string
	ExpireTime time.Duration
}
