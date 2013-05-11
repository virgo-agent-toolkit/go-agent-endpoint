package main

import (
	"github.com/racker/go-agent-endpoint/endpoint"
)

func main() {
	server, _ := endpoint.NewServer(":9876")
	server.Run()
}
