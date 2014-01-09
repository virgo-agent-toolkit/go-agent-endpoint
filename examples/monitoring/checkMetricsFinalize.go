package main

import (
    "github.com/virgo-agent-toolkit/go-agent-endpoint/endpoint"
)

type checkMetricsFinalizeHandler byte

// Handle finalizes check_metrics.post request, writing a non-error response
func (c checkMetricsFinalizeHandler) Handle(req *endpoint.Request, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
    responder.Respond(0, nil)
    return endpoint.OK
}
