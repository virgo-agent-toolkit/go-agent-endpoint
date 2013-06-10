package endpoint

import (
	"encoding/json"
	"errors"
	"net"
)

var (
	// NoSuchHandler means no handler is implemented to handle the request.
	NoSuchHandler = errors.New("No such handler")

	// WrongVersion is the error returned to RPC when protocol version is not
	// compatible.
	WrongVersion = errors.New("Wrong protocol version")
)

// Request is low level RPC request object. Normally should not be used by a
// controller.  DecodeParams should be sufficient for most cases.
type Request struct {
	Version string          `json:"v"`
	ID      string          `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Params  json.RawMessage `json:"params"` // left intact for handlers to deal with
	Method  string          `json:"method"`

	// ugly hacks for detecting message type in bidirectional communication.
	// default to Request because that's way more common than Response to be
	// received on an endpoint
	PlaceHolderResult json.RawMessage `json:"result"`
	PlaceHolderErr    *Error          `json:"error"`
}

// Returns <true, nil> if r is a Request; returns <false, response>, where
// response is a Response object parsed from this message, if r is a Response.
// It's an ugly hack for detecting message type in bidirectional communication.
func (r *Request) isRequestOrGetResponse() (isRequest bool, response *Response) {
	if r.Method != "" {
		return true, nil
	}
	return false, &Response{Version: r.Version, ID: r.ID, Target: r.Target, Source: r.Source, Result: r.PlaceHolderResult, Err: r.PlaceHolderErr}
}

// DecodeParams decodes Params (json:"params") field from the request.
func (r *Request) DecodeParams(v interface{}) error {
	return json.Unmarshal(r.Params, v)
}

// Error is the errors used in the protocol. Use GetErr to convert from error
// to *Error.
type Error struct {
	Field   string `json:"field"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response is low level RPC response object. Normally should not be used by a controller.
// Responder should be sufficient for most cases.
type Response struct {
	Version string          `json:"v"`
	ID      string          `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Result  json.RawMessage `json:"result"`
	Err     *Error          `json:"error"`
}

// ConnContext contains context information about an agent connection.
type ConnContext struct {
	LocalAddr  net.Addr
	RemoteAddr net.Addr
}

func respondingTo(req *Request) *Response {
	return &Response{Version: VERSION, ID: req.ID, Target: req.Source, Source: req.Target, Result: json.RawMessage("{}"), Err: nil}
}

// GetErr converts error to *Error
func GetErr(err error) *Error {
	return &Error{Message: err.Error()}
}

// Responder is what a handler uses to respond to a request.
type Responder struct {
	encodingChan chan<- interface{}
	request      *Request
}

func newResponder(encodingChan chan<- interface{}, request *Request) *Responder {
	return &Responder{encodingChan: encodingChan, request: request}
}

// Respond wraps result with a Response object and put it into sending queue
// that goes back to the agent. This is the method that's sufficient for most
// cases for a controller.
func (r *Responder) Respond(result interface{}, e *Error) (err error) {
	rsp := respondingTo(r.request)
	rsp.Result, err = json.Marshal(result)
	if err != nil {
		return err
	}
	rsp.Err = e
	r.encodingChan <- rsp
	return nil
}

// RespondWithCustomResponse let library user construct their own Response
// object to be sent back to the agent.
func (r *Responder) RespondWithCustomResponse(rsp *Response) {
	r.encodingChan <- rsp
}
