package endpoint

import (
	"bufio"
	"github.com/racker/go-proxy-protocol"
	"net"
	"sync"
)

type Endpoint struct {
	config  EndpointConfig
	ln      net.Listener
	stop    chan int
	wg      *sync.WaitGroup
	once    sync.Once
	running bool
}

type EndpointConfig struct {
	ListenAddr              string
	UpgradingFileServerAddr string

	Hub *Hub
}

func NewEndpoint(config EndpointConfig) (endpoint *Endpoint, err error) {
	endpoint = new(Endpoint)
	endpoint.config = config
	endpoint.wg = new(sync.WaitGroup)
	endpoint.stop = make(chan int, 1)
	endpoint.ln, err = net.Listen("tcp", config.ListenAddr)
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
		e.config.Hub.serveConn(newReadWriter(reader, conn), ConnContext{LocalAddr: conn.LocalAddr(), RemoteAddr: conn.RemoteAddr()})
	} else {
		logger.Printf("Got: %s; not a valid json, will pass to HTTP handler.\n", first)
		handleUpgrade(newReadWriter(reader, conn), e.config.UpgradingFileServerAddr)
	}
}
