package endpoint

import (
	"encoding/json"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

var (
	// RequesterDead is returned when Send() is called on a dead requester.
	RequesterDead = errors.New("Requester is dead already. This is probably because agent has disconnected.")
)

// Requester is an object used for proactively sending messages from controller
// (endpoint) to agents.
type Requester struct {
	request        chan<- *Request
	respondingPath map[int64]chan *Response
	mu             *sync.Mutex
	dead           int32

	remote      string
	local       string
	version     string
	id          int64
	connContext ConnContext
}

func newRequester(requestChan chan<- *Request, remote string, local string, version string, connContext ConnContext) *Requester {
	ret := new(Requester)

	ret.request = requestChan
	ret.respondingPath = make(map[int64]chan *Response)
	ret.mu = new(sync.Mutex)

	ret.local = local
	ret.remote = remote
	ret.version = version
	ret.id = 0 // initial value 0 ?
	ret.connContext = connContext

	return ret
}

// AgentName returns name of the agent that associates with this requester.
func (r *Requester) AgentName() string { return r.remote }

// Version returns protocol version of the agent that associates with this
// requester.
func (r *Requester) Version() string { return r.version }

// AgentAddr returns network address of the agent that associates with this
// requester.
func (r *Requester) AgentAddr() net.Addr { return r.connContext.RemoteAddr }

// Send a message (Request) to this agent. Returns reply, a channel that will
// be used for receiving response from the agent, if successful; or an error if
// the requester is dead.
func (r *Requester) Send(method string, params interface{}) (reply <-chan *Response, err error) {
	if 0 != atomic.LoadInt32(&r.dead) {
		return nil, RequesterDead
	}

	req := new(Request)
	req.PlaceHolderResult = json.RawMessage("{}")
	req.Method = method
	req.Source = r.local
	req.Target = r.remote
	req.Version = r.version
	req.Params, err = json.Marshal(params)
	req.ID = atomic.AddInt64(&r.id, 1)
	if err != nil {
		return
	}

	replyFull := make(chan *Response, 1)
	reply = replyFull
	r.mu.Lock()
	r.respondingPath[req.ID] = replyFull
	r.mu.Unlock()

	r.request <- req
	return
}

func (r *Requester) newResponse(response *Response) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.respondingPath[response.ID] <- response
	delete(r.respondingPath, response.ID)
}

func (r *Requester) die() {
	atomic.AddInt32(&r.dead, 1)
}
