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

// Low level RPC request object. Normally should not be used. Use DecodeParams
// to retrieve parameters
type Request struct {
	Version string          `json:"v"`
	Id      int             `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Params  json.RawMessage `json:"params"` // left intact for handlers to deal with
	Method  string          `json:"method"`
}

func (r *Request) DecodeParams(v interface{}) error {
	return json.Unmarshal(r.Params, v)
}

type Error struct {
	Field   string `json:"field"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Low level RPC response object. Normally should not be used. Use Responder to
// respond to RPC calls
type Response struct {
	Version string          `json:"v"`
	Id      int             `json:"id"`
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
	encoder *json.Encoder
	request *Request
}

func (r *Responder) Respond(result interface{}, e *Error) (err error) {
	rsp := respondingTo(r.request)
	rsp.Result, err = json.Marshal(result)
	if err != nil {
		return err
	}
	rsp.Err = e
	return r.encoder.Encode(rsp)
}

func (r *Responder) RespondWithCustomResponse(rsp *Response) (err error) {
	return r.encoder.Encode(rsp)
}
