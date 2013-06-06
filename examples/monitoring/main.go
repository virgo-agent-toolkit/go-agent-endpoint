package main

import (
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
	"log"
	"os"
)

func print_usage() {
	fmt.Printf("Usage: %s laddr upgradingServerAddr\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		print_usage()
		os.Exit(1)
	}

	hub := endpoint.NewSimpleHub()
	hub.Authenticator(endpoint.DumbAuthenticatorDontUseMe(0), 0)

	config := endpoint.EndpointConfig{}
	config.ListenAddr = os.Args[1]
	config.UpgradingFileServerAddr = os.Args[2]
	config.Hub = hub

	server, err := endpoint.NewEndpoint(config)
	if err != nil {
		print_usage()
		log.Fatalln(err)
	}
	server.Start()
	<-make(chan int)
}
