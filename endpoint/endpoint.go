package endpoint

import (
	"encoding/json"
	"errors"
	"net"
	"sync"
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
	Params  json.RawMessage `json:"params"`
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
	Err     Error           `json:"error"`
}

func respondingTo(req *request) *response {
	return &response{Version: VERSION, Id: req.Id, Target: req.Source, Source: req.Target, Result: json.RawMessage("{}")}
}

func getErr(err error) Error {
	return Error{Message: err.Error()}
}

type handler func(*request, *json.Encoder, *json.Decoder)

type endpoint map[string]handler

func (e endpoint) ServeConn(conn net.Conn, wg *sync.WaitGroup) {
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	req := new(request)
	err := decoder.Decode(req)
	if err != nil {
		return
	}
	if req.Version != VERSION {
		rsp := respondingTo(req)
		rsp.Err = getErr(WrongVersion)
		encoder.Encode(rsp)
	} else {
		handler, ok := e[req.Method]
		if !ok {
			rsp := respondingTo(req)
			rsp.Err = getErr(NoSuchHandler)
			encoder.Encode(rsp)
		} else {
			handler(req, encoder, decoder)
		}
	}
	conn.Close()
	wg.Done()
}
