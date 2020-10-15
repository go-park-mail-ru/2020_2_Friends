package configs

import "time"

const (
	ApiUrl                 = "/api/v1"
	Port                   = ":9000"
	Postgres               = "postgres"
	DataSourceNamePostgres = "host=localhost dbname=grass sslmode=disable"
	ExpireTime             = time.Duration(time.Hour * 24)
	RedisAddr              = "localhost:6379"
)