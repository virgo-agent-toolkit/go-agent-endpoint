package endpoint

import (
	"encoding/json"
	"io"
)

const (
	LowestPriority = int(^uint(0) >> 1) // maximum value for int
)

// Hub is a wrapper to all handlers. It supports multiple handlers to a same
// trigger. Internally a priority queue is used to determine the order in which
// handlers are executed. When requested, Hub tries all handlers hooked to the
// trigger, from priority number lower to higher, until a handler properly
// handles the request.
type Hub struct {
	handlers       map[string]*handlerList
	authenticators *authenticatorList
	unhandled      *handlerList
}

func NewHub() *Hub {
	ret := new(Hub)
	ret.authenticators = newAuthenticatorList()
	ret.unhandled = newHandlerList()
	ret.handlers = make(map[string]*handlerList)
	ret.Hook("heartbeat.post", HeartbeatHandler(0), LowestPriority)
	ret.Unhandled(Unhandled(0), LowestPriority)
	return ret
}

// Hook a handler (handler) to a method (trigger) with priority. There can be
// multiple handlers for the same trigger. When an rpc request on trigger is
// received, the hub is requested by the endpoint to iterate through handlers
// hooked on the trigger, starting from lowest priority number, until the
// request is (considered) properly handled.
func (h *Hub) Hook(trigger string, handler Handler, priority int) {
	if _, ok := h.handlers[trigger]; !ok {
		h.handlers[trigger] = newHandlerList()
	}
	h.handlers[trigger].Push(handlerListItem{handler: handler, priority: priority})
}

// Hook an authentication handler (authenticator) with priority. There can be
// multiple authentication handlers. When a new handshake is initiated,
// endpoint tries to authenticate the ageng. The hub is requested by the
// endpoint to iterate through all authenticators, starting from lowest
// priority number, until authentication succeeds or all authenticators are
// tried.
func (h *Hub) Authenticator(authenticator Authenticator, priority int) {
	h.authenticators.Push(constructAuthenticatorListItem(authenticator, priority))
}

// Hook an unhandled handler with priority. An unhandled handler is used when
// no handled can handle the request. There can be multiple unhandled handlers.
// The execution rule is similar to regular handler.
func (h *Hub) Unhandled(handler Handler, priority int) {
	h.unhandled.Push(handlerListItem{handler, priority})
}

func (h *Hub) serveConn(rw io.ReadWriter, connCtx ConnContext) {
	if !h.authenticate(rw, connCtx) {
		return
	}

	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(rw)
	for {
		req := new(Request)
		err := decoder.Decode(req)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Decoding error: %v\n", err)
			}
			break
		}
		if req.Version != VERSION {
			rsp := respondingTo(req)
			rsp.Err = GetErr(WrongVersion)
			encoder.Encode(rsp)
		} else {
			handlerList, ok := h.handlers[req.Method]
			if !ok {
				h.unhandled.Iterate(req, &Responder{encoder: encoder, request: req}, connCtx)
			} else {
				handlerList.Iterate(req, &Responder{encoder: encoder, request: req}, connCtx)
			}
		}
	}
}

func (h *Hub) authenticate(rw io.ReadWriter, connCtx ConnContext) (authenticated bool) {
	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(rw)
	req := new(Request)
	err := decoder.Decode(req)
	if err != nil {
		return false
	}
	responder := Responder{encoder: encoder, request: req}
	if req.Method != "handshake.hello" {
		responder.Respond(nil, GetErr(AuthenticationFailed))
		return false
	}
	if req.Version != VERSION {
		responder.Respond(nil, GetErr(WrongVersion))
		return false
	}
	var hl HelloParams
	err = json.Unmarshal(req.Params, &hl)
	if err != nil {
		responder.Respond(nil, GetErr(err))
		return false
	}
	logger.Printf("got a handshake.hello from %s\n", req.Source)
	// TODO: check process version and bundleversion
	if OK != h.authenticators.Iterate(hl.AgentName, hl.AgentId, hl.Token, connCtx) {
		responder.Respond(nil, GetErr(AuthenticationFailed))
		logger.Printf("handshake.hello from %s failed authentication\n", req.Source)
		return false
	}
	responder.Respond(HelloResult{HeartbeatInterval: "1000"}, nil)
	return true
}
