package endpoint

import (
	"bufio"
	"errors"
	"github.com/racker/go-proxy-protocol"
	"net"
	"net/http"
	"sync"
)

var (
	ServerAlreadyStarted = errors.New("Server is already started; can't bind more handlers")
	DuplicateMethod      = errors.New("A handler with this method name already exists")
)

type Server struct {
	ep          endpoint
	ln          net.Listener
	httpHandler http.Handler

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
	server.httpHandler = http.NotFoundHandler()
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
					go s.serveConn(conn, s.wg)
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

func (s *Server) serveConn(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	var err error
	reader := bufio.NewReader(conn)
	_, err = proxyProtocol.ConsumeProxyLine(reader)
	if err != nil {
		return
	}
	first, err := reader.Peek(1)
	for err == nil && (first[0] == ' ' || first[0] == '\t' || first[0] == '\n' || first[0] == '\r') {
		reader.ReadByte()
		first, err = reader.Peek(1)
	}
	if err != nil {
		return
	}
	if first[0] == '{' {
		// writing shouldn't be buffered
		s.ep.ServeConn(newReadWriter(reader, conn))
	} else {
		logger.Printf("Got: %s; not a valid json, will pass to HTTP handler.\n", first)
		handleUpgrade(newReadWriter(reader, conn))
	}
}
