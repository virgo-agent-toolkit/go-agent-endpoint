package endpoint

import (
	"encoding/json"
	"errors"
	"net"
)

var (
	NoSuchHandler = errors.New("No such handler")
	WrongVersion  = errors.New("Wrong protocol version")
)

// Low level RPC request object. Normally should not be used by a controller.
// DecodeParams should be sufficient for most cases.
type Request struct {
	Version string          `json:"v"`
	Id      int64           `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Params  json.RawMessage `json:"params"` // left intact for handlers to deal with
	Method  string          `json:"method"`

	// ugly hacks for detecting message type in bidirectional communication.
	// default to Request because that's way more common than Response to be
	// received on an endpoint
	Place_holder_result json.RawMessage `json:"result"`
	Place_holder_err    *Error          `json:"error"`
}

// Returns <true, nil> if r is a Request; returns <false, response>, where
// response is a Response object parsed from this message, if r is a Response.
// It's an ugly hack for detecting message type in bidirectional communication.
func (r *Request) isRequestOrGetResponse() (isRequest bool, response *Response) {
	if r.Method != "" {
		return true, nil
	} else {
		return false, &Response{Version: r.Version, Id: r.Id, Target: r.Target, Source: r.Source, Result: r.Place_holder_result, Err: r.Place_holder_err}
	}
}

func (r *Request) DecodeParams(v interface{}) error {
	return json.Unmarshal(r.Params, v)
}

type Error struct {
	Field   string `json:"field"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Low level RPC response object. Normally should not be used by a controller.
// Responder should be sufficient for most cases.
type Response struct {
	Version string          `json:"v"`
	Id      int64           `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Result  json.RawMessage `json:"result"`
	Err     *Error          `json:"error"`
}

type ConnContext struct {
	LocalAddr  net.Addr
	RemoteAddr net.Addr
}

func respondingTo(req *Request) *Response {
	return &Response{Version: VERSION, Id: req.Id, Target: req.Source, Source: req.Target, Result: json.RawMessage("{}"), Err: nil}
}

func GetErr(err error) *Error {
	return &Error{Message: err.Error()}
}

type Responder struct {
	encodingChan chan<- interface{}
	request      *Request
}

func newResponder(encodingChan chan<- interface{}, request *Request) *Responder {
	return &Responder{encodingChan: encodingChan, request: request}
}

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

func (r *Responder) RespondWithCustomResponse(rsp *Response) {
	r.encodingChan <- rsp
}
