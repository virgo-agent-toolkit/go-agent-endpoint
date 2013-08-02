package main

import (
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
)

func proactive(requesters <-chan *endpoint.Requester) {
	for {
		requester := <-requesters
		go ping(requester)
	}
}

func ping(requester *endpoint.Requester) {
  req := &WriteFileTextParam{Path: "/tmp/hello", Content: "Hello, world!"}
	reply, err := requester.Send("write_file.text", req)
	if err != nil {
		fmt.Println(err)
	}
	rsp := <-reply
  fmt.Printf("WriteTextFile response: %#v\n", rsp)
}
