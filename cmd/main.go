package main

import (
	"github.com/friends/internal/app/fileserver"
	"github.com/friends/internal/app/server"
)

func main() {
	go fileserver.StartFileServer()
	server.StartApiServer()
}
