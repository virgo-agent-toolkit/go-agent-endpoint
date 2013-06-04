package endpoint

import (
	"encoding/json"
	"io"
)

type Hub struct {
	handlers       map[string]*handlerList
	authenticators *authenticatorList
}

func NewHub() *Hub {
	ret := new(Hub)
	ret.authenticators = newAuthenticatorList()
	ret.handlers = make(map[string]*handlerList)
	ret.Hook("heartbeat.post", HeartbeatHandler(0), 0)
	return ret
}

func (h *Hub) Hook(trigger string, handler Handler, priority int) {
	if _, ok := h.handlers[trigger]; !ok {
		h.handlers[trigger] = newHandlerList()
	}
	h.handlers[trigger].Push(handlerListItem{handler: handler, priority: priority})
}

func (h *Hub) Authenticator(authenticator Authenticator, priority int) {
	h.authenticators.Push(constructAuthenticatorListItem(authenticator, priority))
}

func (h *Hub) serveConn(rw io.ReadWriter) {
	if !h.authenticate(rw) {
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
			handlerList, ok := h.handlers[req.Method]
			if !ok {
				rsp := respondingTo(req)
				rsp.Err = getErr(NoSuchHandler)
				logger.Printf("Got a request to unimplemented handler: %s\n", req.Method)
				encoder.Encode(rsp)
			} else {
				handlerList.Iterate(req, encoder, decoder)
			}
		}
	}
}

func (h *Hub) authenticate(rw io.ReadWriter) (authenticated bool) {
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
	if OK != h.authenticators.Iterate(hl.AgentName, hl.AgentId, hl.Token) {
		rsp.Err = getErr(AuthenticationFailed)
		logger.Printf("handshake.hello from %s failed authentication\n", req.Source)
		return false
	}
	rsp.Result, _ = json.Marshal(HelloResult{HeartbeatInterval: "1000"})
	return true
}
