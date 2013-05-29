package endpoint

import (
	"io"
)

type readWriter struct {
	reader io.Reader
	writer io.Writer
}

func newReadWriter(reader io.Reader, writer io.Writer) readWriter {
	return readWriter{reader: reader, writer: writer}
}

func (rw readWriter) Read(p []byte) (n int, err error) {
	n, err = rw.reader.Read(p)
	return
}

func (rw readWriter) Write(p []byte) (n int, err error) {
	n, err = rw.writer.Write(p)
	return
}
