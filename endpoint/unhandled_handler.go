package endpoint

import (
	"encoding/json"
)

type Unhandled byte

func (h Unhandled) Handle(req *request, encoder *json.Encoder, connCtx ConnContext) HandleCode {
	rsp := respondingTo(req)
	rsp.Err = getErr(NoSuchHandler)
	logger.Printf("Got a request to unimplemented handler: %s\n", req.Method)
	encoder.Encode(rsp)
	return OK
}
