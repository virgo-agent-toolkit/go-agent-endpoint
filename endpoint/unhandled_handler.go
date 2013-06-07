package endpoint

import ()

// Unhandled is a default handler for unhandled requests
type Unhandled byte

// Handle returns a NoSuchHandler error to RPC
func (h Unhandled) Handle(req *Request, responder *Responder, connCtx ConnContext) HandleCode {
	responder.Respond(nil, GetErr(NoSuchHandler))
	logger.Printf("Got a request to unimplemented handler: %s\n", req.Method)
	return OK
}
