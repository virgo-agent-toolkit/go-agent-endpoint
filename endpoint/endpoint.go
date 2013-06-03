package endpoint

import (
	"encoding/json"
	"errors"
	"io"
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

func respondingTo(req *request) *response {
	return &response{Version: VERSION, Id: req.Id, Target: req.Source, Source: req.Target, Result: json.RawMessage("{}"), Err: nil}
}

func getErr(err error) *Error {
	return &Error{Message: err.Error()}
}

type Handler func(*request, *json.Encoder, *json.Decoder)

type endpoint struct {
	Handlers map[string]Handler
	ctrl     *controller
}

func newEndpoint(ctrl *controller) *endpoint {
	ret := new(endpoint)
	ret.ctrl = ctrl
	ret.Handlers = make(map[string]Handler)
	ret.Handlers["heartbeat.post"] = ret.handleHeartbeat
	return ret
}

func (e endpoint) ServeConn(rw io.ReadWriter) {
	if !e.authenticate(rw) {
		return
	}

	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(rw)
	for {
		req := new(request)
		err := decoder.Decode(req)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Decoding error: %v\n", err)
			}
			break
		}
		if req.Version != VERSION {
			rsp := respondingTo(req)
			rsp.Err = getErr(WrongVersion)
			encoder.Encode(rsp)
		} else {
			handler, ok := e.Handlers[req.Method]
			if !ok {
				rsp := respondingTo(req)
				rsp.Err = getErr(NoSuchHandler)
				logger.Printf("Got a request to unimplemented handler: %s\n", req.Method)
				encoder.Encode(rsp)
			} else {
				handler(req, encoder, decoder)
			}
		}
	}
}

func (e *endpoint) authenticate(rw io.ReadWriter) (authenticated bool) {
	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(rw)
	req := new(request)
	err := decoder.Decode(req)
	if err != nil {
		return false
	}
	rsp := respondingTo(req)
	defer encoder.Encode(rsp)
	if req.Method != "handshake.hello" {
		rsp.Err = getErr(AuthenticationFailed)
		return false
	}
	if req.Version != VERSION {
		rsp.Err = getErr(WrongVersion)
		return false
	}
	var hl HelloParams
	err = json.Unmarshal(req.Params, &hl)
	if err != nil {
		rsp.Err = getErr(err)
		return false
	}
	logger.Printf("got a handshake.hello from %s\n", req.Source)
	// TODO: check process version and bundleversion
	if !e.ctrl.Authenticate(hl.AgentName, hl.AgentId, hl.Token) {
		rsp.Err = getErr(AuthenticationFailed)
		logger.Printf("handshake.hello from %s failed authentication\n", req.Source)
		return false
	}
	rsp.Result, _ = json.Marshal(HelloResult{HeartbeatInterval: "1000"})
	return true
}
