package endpoint

import (
	"bufio"
	"errors"
	"github.com/racker/go-proxy-protocol"
	"net"
	"sync"
)

var (
	ServerAlreadyStarted = errors.New("Server is already started; can't bind more handlers")
	DuplicateMethod      = errors.New("A handler with this method name already exists")
	AuthenticationFailed = errors.New("Authentication failed")
)

type Endpoint struct {
	hub *Hub
	ln  net.Listener

	stop    chan int
	wg      *sync.WaitGroup
	once    sync.Once
	running bool
}

func NewEndpoint(laddr string, hub *Hub) (endpoint *Endpoint, err error) {
	endpoint = new(Endpoint)
	endpoint.wg = new(sync.WaitGroup)
	endpoint.stop = make(chan int, 1)
	endpoint.ln, err = net.Listen("tcp", laddr)
	endpoint.hub = hub
	return
}

func (e *Endpoint) Start() {
	go e.once.Do(func() {
		e.running = true
		for e.running {
			select {
			case <-e.stop:
				e.running = false
			default:
				conn, err := e.ln.Accept()
				if err == nil {
					e.wg.Add(1)
					go e.serveConn(conn, e.wg)
				}
			}
		}
	})
}

func (e *Endpoint) Destroy() {
	e.stop <- 1
	e.ln.Close()
	e.wg.Wait()
}

func (e *Endpoint) serveConn(conn net.Conn, wg *sync.WaitGroup) {
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
		e.hub.serveConn(newReadWriter(reader, conn))
	} else {
		logger.Printf("Got: %s; not a valid json, will pass to HTTP handler.\n", first)
		handleUpgrade(newReadWriter(reader, conn))
	}
}
