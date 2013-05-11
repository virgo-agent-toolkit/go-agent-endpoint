package endpoint

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Server struct {
	ep endpoint
	ln net.Listener
}

func NewServer(laddr string) (server *Server, err error) {
	server = &Server{ep: make(endpoint)}
	server.ln, err = net.Listen("tcp", laddr)
	server.bind()
	return
}

func (s *Server) Run() {
	for {
		conn, err := s.ln.Accept()
		if err == nil {
			go s.ep.ServeConn(conn)
		}
	}
}

func (s *Server) bind() {
	s.ep["heartbeat.post"] = s.serveHeartbeat
}

func (e *Server) serveHeartbeat(req *request, encoder *json.Encoder, decoder *json.Decoder) {
	rsp := respondingTo(req)
	var hb Heartbeat
	err := json.Unmarshal(req.Params, &hb)
	if err != nil {
		rsp.Err = getErr(err)
	} else {
		fmt.Printf("Got a timestamp: %v\n", hb.Timestamp)
		rsp.Result, _ = Heartbeat{Timestamp: time.Now()}.MarshalJSON()
	}
	encoder.Encode(rsp)
}
