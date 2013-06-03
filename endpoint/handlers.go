package endpoint

import (
	"encoding/json"
	"time"
)

func (e *endpoint) handleHeartbeat(req *request, encoder *json.Encoder, decoder *json.Decoder) {
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
}
