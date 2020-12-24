package configs

import "time"

const (
	APIURL             = "/api/v1"
	Port               = ":9000"
	FileServerPort     = ":9001"
	FileServerGRPCPort = ":9002"
	SessionServicePort = ":9003"
	Postgres           = "postgres"
	ExpireTime         = time.Hour * 24
	RedisAddr          = "localhost:6379"
	ReqID              = "reqID"
	UserID             = "userID"
	SessionID          = "session_id"
	CookieCSRF         = "X-CSRF-Cookie"
	ImgMaxSize         = 1024 * 1024
	AvatarFormFileKey  = "avatar"
	ImgFormFileKey     = "image"
	FileServerPath     = "./static"
	ImageDir           = "./static/img/"
	ProductID          = "product_id"
	UserRole           = 1
	AdminRole          = 2
	TimeFormat         = "02.01.2006 15:04:05"
	Longitude          = "longitude"
	Latitude           = "latitude"
)
