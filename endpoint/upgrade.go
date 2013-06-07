package endpoint

import (
	"io"
	"net"
)

func handleUpgrade(rw io.ReadWriter, upgradingFileServer string) {
	conn, err := net.Dial("tcp", upgradingFileServer)
	if err == nil {
		go io.Copy(conn, rw)
		io.Copy(rw, conn)
	}
}
