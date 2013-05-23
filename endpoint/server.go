package endpoint

import (
	"errors"
	"net"
	"sync"
)

var (
	ServerAlreadyStarted = errors.New("Server is already started; can't bind more handlers")
	DuplicateMethod      = errors.New("A handler with this method name already exists")
)

type Server struct {
	ep endpoint
	ln net.Listener

	stop    chan int
	wg      *sync.WaitGroup
	once    sync.Once
	running bool

	agents map[string]*agent
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
		s.running = true
		for s.running {
			select {
			case <-s.stop:
				s.running = false
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
	s.ep["handshake.hello"] = s.handleHandshakeHello
	s.ep["heartbeat.post"] = s.handleHeartbeat
}

func (s *Server) Bind(method string, handler Handler) error {
	if s.running {
		return ServerAlreadyStarted
	}
	if _, ok := s.ep[method]; ok {
		return DuplicateMethod
	}
	s.ep[method] = handler
	return nil
}
