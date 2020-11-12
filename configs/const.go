package configs

import "time"

const (
	ApiUrl                 = "/api/v1"
	Port                   = ":9000"
	FileServerPort         = ":9001"
	Postgres               = "postgres"
	DataSourceNamePostgres = "host=localhost dbname=grass sslmode=disable"
	ExpireTime             = time.Duration(time.Hour * 24)
	CSRFTokenExpireTime    = time.Duration(time.Hour * 24)
	RedisAddr              = "localhost:6379"
	ReqID                  = "reqID"
	UserID                 = "userID"
	SessionID              = "session_id"
	ImgMaxSize             = 1024 * 1024
	AvatarFormFileKey      = "avatar"
	ImgFormFileKey         = "image"
	FileServerPath         = "./static"
	ProductID              = "product_id"
	UserRole               = 1
	AdminRole              = 2
)
