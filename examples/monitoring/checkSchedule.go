package main

import (
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
)

type checkScheduleHandler byte

// Handle responds to check_schedule.get requests with a list of checks
func (c checkScheduleHandler) Handle(req *endpoint.Request, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
	var rsp CheckScheduleGetResult
	checks := []string{"agent.memory", "agent.disk", "agent.filesystem", "agent.network", "agent.cpu", "agent.load_average"}
	for i, v := range checks {
		rsp.Checks = append(rsp.Checks, Check{
			Id:       fmt.Sprintf("check-%2d", i),
			Type:     v,
			Details:  nil,
			Period:   16,
			Timeout:  16,
			Disabled: false,
		})
	}
	responder.Respond(rsp, nil)
	return endpoint.OK
}
