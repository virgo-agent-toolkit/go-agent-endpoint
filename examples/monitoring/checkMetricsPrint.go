package main

import (
	"encoding/json"
	"fmt"
	"github.com/racker/go-agent-endpoint/endpoint"
)

type checkMetricsPrintHandler byte

// Handle parses to check_metrics.post requests and print out the posted
// metrics; it does not finalize (responde to) the request
func (c checkMetricsPrintHandler) Handle(req *endpoint.Request, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
	var params checkMetricsPostParams
	err := json.Unmarshal(req.Params, &params)
	if err != nil { // parsing failed, should not go on
		fmt.Printf("parsing check_metrics.post Params failed: %v\n", err)
		return endpoint.FAIL
	}

	fmt.Printf("%+q\n", params)

	return endpoint.DECLINED // not finalized; other handler may need the data
}
