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

	hub, requesters := endpoint.NewHub()
	hub.Authenticator(authenticator(0), 0)
	hub.Hook("check_schedule.get", checkScheduleHandler(0), 0)
	hub.Hook("check_metrics.post", checkMetricsHandler(0), 0)

	go proactive(requesters)

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
