package response

import (
	"errors"
	"fmt"
	"http/internal/headers"
	"io"
	"net"
)

type StatusCode int

const twohundred StatusCode = 1
const fourhundred StatusCode = 2
const fivehundred StatusCode = 3

type Writer struct {
	write_stream io.Writer
}

func NewWriter(con net.Conn) *Writer {
	return &Writer{write_stream: con}
}

func (w *Writer) WriteStatusLinesW(StatusCode StatusCode) error {
	var Statusline []byte
	switch StatusCode {
	case twohundred:
		Statusline = []byte("HTTP/1.1 200 OK\r\n")
	case fourhundred:
		Statusline = []byte("HTTP/1.1 400 Bad Request\r\n")

	case fivehundred:
		Statusline = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return errors.New("not implemented status code")
	}
	_, err := w.write_stream.Write(Statusline)
	return err
}

func (w *Writer) WriteHeadersW(headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.write_stream.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.write_stream.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) error {
	_, err := w.write_stream.Write(p)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	if contentLen == -1 {
		head := headers.NewHeaders()
		head["connection"] = "close"
		head["content-type"] = "text/plain"
		return head
	}
	head := headers.NewHeaders()
	head["content-length"] = fmt.Sprintf("%d", contentLen)
	head["connection"] = "close"
	head["content-type"] = "text/plain"
	return head
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	w.write_stream.Write([]byte(fmt.Sprintf("%x\r\n", len(p))))
	w.write_stream.Write(p)
	return w.write_stream.Write([]byte("\r\n"))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.write_stream.Write([]byte("0\r\n\r\n"))
}
