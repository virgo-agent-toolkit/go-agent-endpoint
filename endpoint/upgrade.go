package endpoint

import (
	"io"
	"net"
)

var (
	upgradingFileServer = "localhost:8989" // should be in config in the future
)

func handleUpgrade(rw io.ReadWriter) {
	conn, err := net.Dial("tcp", upgradingFileServer)
	if err == nil {
		go io.Copy(conn, rw)
		io.Copy(rw, conn)
	}
}
