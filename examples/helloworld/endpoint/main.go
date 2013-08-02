package main

import (
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
	"log"
	"os"
)

func printUsage() {
	fmt.Printf("Usage: %s laddr upgradingServerAddr\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	hub, requesters := endpoint.NewHub()

	go proactive(requesters)

	config := endpoint.EndpointConfig{}
	config.ListenAddr = os.Args[1]
	config.UpgradingFileServerAddr = os.Args[2]
	config.Hub = hub

	server, err := endpoint.NewEndpoint(config)
	if err != nil {
		printUsage()
		log.Fatalln(err)
	}
	server.Start()
	select {}
}
