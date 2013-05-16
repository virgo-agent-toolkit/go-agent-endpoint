package endpoint

import (
	"encoding/json"
	"time"
)

func (s *Server) handleHeartbeat(req *request, encoder *json.Encoder, decoder *json.Decoder) {
	rsp := respondingTo(req)
	var hb Heartbeat
	err := json.Unmarshal(req.Params, &hb)
	if err != nil {
		rsp.Err = getErr(err)
	} else {
		//logger.Printf("Got a timestamp: %v\n", hb.Timestamp)
		rsp.Result, _ = Heartbeat{Timestamp: time.Now()}.MarshalJSON()
	}
	encoder.Encode(rsp)
}

func (s *Server) handleHandshakeHello(req *request, encoder *json.Encoder, decoder *json.Decoder) {
}
