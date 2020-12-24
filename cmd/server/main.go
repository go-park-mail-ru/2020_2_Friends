package main

import (
	"os"

	"github.com/friends/internal/app/server"
)

func main() {
	dsn := os.Getenv("dsn")
	server.StartAPIServer(dsn)
}
