package endpoint

import (
	"encoding/json"
	"io"
)

const (
	// LowestPriority is maximum value for int
	LowestPriority = int(^uint(0) >> 1)
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
	newRequesters  chan *Requester
}

// NewSimpleHub creates a Hub that supports only single-directional
// communication from agents to endpoint
func NewSimpleHub() (hub *Hub) {
	ret := new(Hub)
	ret.authenticators = newAuthenticatorList()
	ret.unhandled = newHandlerList()
	ret.handlers = make(map[string]*handlerList)
	ret.Hook("heartbeat.post", HeartbeatHandler(0), LowestPriority)
	ret.Unhandled(Unhandled(0), LowestPriority)
	return ret
}

// NewHub creates a Hub that supports bi-directional communication with agents
func NewHub() (hub *Hub, newRequesters <-chan *Requester) {
	ret := NewSimpleHub()
	ret.newRequesters = make(chan *Requester)
	return ret, ret.newRequesters
}

// Hook hooks a handler (handler) to a method (trigger) with priority. There
// can be multiple handlers for the same trigger. When an rpc request on
// trigger is received, the hub is requested by the endpoint to iterate through
// handlers hooked on the trigger, starting from lowest priority number, until
// the request is (considered) properly handled.
func (h *Hub) Hook(trigger string, handler Handler, priority int) {
	if _, ok := h.handlers[trigger]; !ok {
		h.handlers[trigger] = newHandlerList()
	}
	h.handlers[trigger].Push(handlerListItem{handler: handler, priority: priority})
}

// Authenticator hooks an authentication handler (authenticator) with priority.
// There can be multiple authentication handlers. When a new handshake is
// initiated, endpoint tries to authenticate the ageng. The hub is requested by
// the endpoint to iterate through all authenticators, starting from lowest
// priority number, until authentication succeeds or all authenticators are
// tried.
func (h *Hub) Authenticator(authenticator Authenticator, priority int) {
	h.authenticators.Push(constructAuthenticatorListItem(authenticator, priority))
}

// Unhandled hooks an "unhandled" handler with priority. An "unhandled" handler
// is used when no handled can handle the request. There can be multiple
// unhandled handlers.  The execution rule is similar to regular handler.
func (h *Hub) Unhandled(handler Handler, priority int) {
	h.unhandled.Push(handlerListItem{handler, priority})
}

func (h *Hub) serveConn(rw io.ReadWriter, connCtx ConnContext) {
	encodingChan := make(chan interface{})
	go func() { // encoder worker (aggregating different messages)
		encoder := json.NewEncoder(rw)
		for {
			encoder.Encode(<-encodingChan)
		}
	}()

	authenticated, authReq := h.authenticate(rw, connCtx, encodingChan)
	if !authenticated {
		return
	}

	var requester *Requester
	if h.newRequesters != nil {
		requestChan := make(chan *Request)
		requester = newRequester(requestChan, authReq.Source, authReq.Target, authReq.Version, connCtx)
		go func() { // requestChan --> encodingChan
			h.newRequesters <- requester
			for {
				encodingChan <- <-requestChan
			}
		}()
	}

	decoder := json.NewDecoder(rw)
	for { // decoder worker
		// Decode a message, could be response as well
		req := new(Request)
		err := decoder.Decode(req)
		if err != nil {
			if err != io.EOF {
				logger.Printf("Decoding error: %v\n", err)
			}
			break
		}

		// determine whether it's request or response
		isReq, rsp := req.isRequestOrGetResponse()
		if !isReq && requester != nil {
			// it's a response to a request from requester, should be given back to
			// the requester
			requester.newResponse(rsp)
		} else {
			// it's a request; should use handlers
			responder := newResponder(encodingChan, req)
			if responder.request.Version != VERSION {
				responder.Respond(nil, GetErr(WrongVersion))
			} else {
				handlerList, ok := h.handlers[responder.request.Method]
				if !ok {
					h.unhandled.Iterate(responder.request, responder, connCtx)
				} else {
					handlerList.Iterate(responder.request, responder, connCtx)
				}
			}
		}
	}
}

func (h *Hub) authenticate(rw io.Reader, connCtx ConnContext, encodingChan chan interface{}) (authenticated bool, req *Request) {
	decoder := json.NewDecoder(rw)
	req = new(Request)
	err := decoder.Decode(req)
	if err != nil {
		return false, req
	}
	responder := newResponder(encodingChan, req)
	if req.Method != "handshake.hello" {
		responder.Respond(nil, GetErr(AuthenticationFailed))
		return false, req
	}
	if req.Version != VERSION {
		responder.Respond(nil, GetErr(WrongVersion))
		return false, req
	}
	var hl HelloParams
	err = json.Unmarshal(req.Params, &hl)
	if err != nil {
		responder.Respond(nil, GetErr(err))
		return false, req
	}
	logger.Printf("got a handshake.hello from %s\n", req.Source)
	// TODO: check process version and bundleversion
	if OK != h.authenticators.Iterate(hl.AgentName, hl.AgentID, hl.Token, responder, connCtx) {
		responder.Respond(nil, GetErr(AuthenticationFailed))
		logger.Printf("handshake.hello from %s failed authentication\n", req.Source)
		return false, req
	}
	return true, req
}
