package main

import (
    "fmt"

    "github.com/virgo-agent-toolkit/go-agent-endpoint/endpoint"
)

const scheduleInterval = 1

type checkScheduleHandler byte

// Handle responds to check_schedule.get requests with a list of checks
func (c checkScheduleHandler) Handle(req *endpoint.Request, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
    var rsp checkScheduleGetResult
    checks := [][]interface{}{
        {"agent.memory", nil},
        {"agent.disk", map[string]string{"target": "/dev/mapper/precise64-root"}},
        {"agent.filesystem", map[string]string{"target": "/"}},
        {"agent.network", map[string]string{"target": "eth0"}},
        {"agent.network", map[string]string{"target": "eth1"}},
        {"agent.cpu", nil},
        {"agent.load_average", nil},
    }
    for i, v := range checks {
        rsp.Checks = append(rsp.Checks, check{
            ID:       fmt.Sprintf("check-%02d", i),
            Type:     v[0].(string),
            Details:  v[1],
            Period:   scheduleInterval,
            Timeout:  scheduleInterval,
            Disabled: false,
        })
    }
    responder.Respond(rsp, nil)
    return endpoint.OK
}
