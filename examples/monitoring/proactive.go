package main

import (
	"encoding/json"
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
)

func proactive(requesters <-chan *endpoint.Requester) {
	for {
		requester := <-requesters
		go askForSystemInfo(requester)
	}
}

func askForSystemInfo(requester *endpoint.Requester) {
	reply, err := requester.Send("system.info", nil)
	if err != nil {
		fmt.Println(err)
	}
	rsp := <-reply
	var sysInfo systemInfoResponse
	json.Unmarshal(rsp.Result, &sysInfo)
	fmt.Printf("%#v\n", sysInfo)
}
