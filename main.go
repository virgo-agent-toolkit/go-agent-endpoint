package main

import (
	//  "github.com/racker/go-agent-endpoint/monitoring"
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
	"log"
	"os"
)

func print_usage() {
	fmt.Printf("Usage: %s laddr\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		print_usage()
		os.Exit(1)
	}
	server, err := endpoint.NewEndpoint(os.Args[1], endpoint.NewHub())
	if err != nil {
		print_usage()
		log.Fatalln(err)
	}
	server.Start()
	<-make(chan int)
}
