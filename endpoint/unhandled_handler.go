package endpoint

import ()

type Unhandled byte

func (h Unhandled) Handle(req *Request, responder *Responder, connCtx ConnContext) HandleCode {
	responder.Respond(nil, GetErr(NoSuchHandler))
	logger.Printf("Got a request to unimplemented handler: %s\n", req.Method)
	return OK
}
