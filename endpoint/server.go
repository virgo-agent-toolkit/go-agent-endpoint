package endpoint

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	ep   endpoint
	ln   net.Listener
	stop chan int
	wg   *sync.WaitGroup
	once sync.Once
}

func NewServer(laddr string) (server *Server, err error) {
	server = &Server{ep: make(endpoint)}
	server.wg = new(sync.WaitGroup)
	server.stop = make(chan int, 1)
	server.ln, err = net.Listen("tcp", laddr)
	server.bind()
	return
}

func (s *Server) Start() {
	go s.once.Do(func() {
		run := true
		for run {
			select {
			case <-s.stop:
				run = false
			default:
				conn, err := s.ln.Accept()
				if err == nil {
					s.wg.Add(1)
					go s.ep.ServeConn(conn, s.wg)
				}
			}
		}
	})
}

func (s *Server) Destroy() {
	s.stop <- 1
	s.ln.Close()
	s.wg.Wait()
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
