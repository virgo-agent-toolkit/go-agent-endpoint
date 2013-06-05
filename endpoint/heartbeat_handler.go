package endpoint

import (
	"encoding/json"
	"time"
)

type HeartbeatHandler byte

func (h HeartbeatHandler) Handle(req *request, encoder *json.Encoder, connCtx ConnContext) HandleCode {
	rsp := respondingTo(req)
	var hb Heartbeat
	err := json.Unmarshal(req.Params, &hb)
	if err != nil {
		rsp.Err = getErr(err)
		logger.Printf("Unmarshaling heartbeat error: %v\n", err)
	} else {
		logger.Printf("Got a heartbeat: %v\n", hb.Timestamp)
		rsp.Result, _ = Heartbeat{Timestamp: time.Now()}.MarshalJSON()
	}
	encoder.Encode(rsp)
	return OK
}
