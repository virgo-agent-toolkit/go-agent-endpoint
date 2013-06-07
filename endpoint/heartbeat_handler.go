package endpoint

import (
	"time"
)

// HeartbeatHandler is a default handler for "heartbeat.post" requests
type HeartbeatHandler byte

// Handle parses a "heartbeat.post" request, and write a timestamp as response.
func (h HeartbeatHandler) Handle(req *Request, responder *Responder, connCtx ConnContext) HandleCode {
	var hb Heartbeat
	err := req.DecodeParams(&hb)
	if err != nil {
		logger.Printf("Unmarshaling heartbeat error: %v\n", err)
		responder.Respond(nil, GetErr(err))
	} else {
		logger.Printf("Got a heartbeat: %v\n", hb.Timestamp)
		responder.Respond(Heartbeat{Timestamp: time.Now()}, nil)
	}
	return OK
}
