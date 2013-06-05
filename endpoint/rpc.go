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

type request struct {
	Version string          `json:"v"`
	Id      int             `json:"id"`
	Target  string          `json:"target"`
	Source  string          `json:"source"`
	Params  json.RawMessage `json:"params"` // left intact for handles to deal with
	Method  string          `json:"method"`
}

type Error struct {
	Field   string `json:"field"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type response struct {
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

func respondingTo(req *request) *response {
	return &response{Version: VERSION, Id: req.Id, Target: req.Source, Source: req.Target, Result: json.RawMessage("{}"), Err: nil}
}

func getErr(err error) *Error {
	return &Error{Message: err.Error()}
}
