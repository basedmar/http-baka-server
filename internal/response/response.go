package response

import (
	"errors"
	"fmt"
	"http/internal/headers"
	"io"
)

type StatusCode int

const twohundred StatusCode = 1
const fourhundred StatusCode = 2
const fivehundred StatusCode = 3

type Writer struct {
	Statusline []byte
	Headers    headers.Headers
	Body       []byte
}

func (w *Writer) WriteStatusLinesW(StatusCode StatusCode) error {
	switch StatusCode {
	case twohundred:
		w.Statusline = []byte("HTTP/1.1 200 OK\r\n")
	case fourhundred:
		w.Statusline = []byte("HTTP/1.1 400 Bad Request\r\n")

	case fivehundred:
		w.Statusline = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return errors.New("not implemented status code")
	}
	return nil
}

func (w *Writer) WriteHeadersW(headers headers.Headers) error {
	if w.Headers == nil {
		w.Headers = make(map[string]string)
	}
	for k, v := range headers {
		w.Headers[k] = v
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) error {
	w.Body = p
	return nil
}

func WriteStatusLines(w io.Writer, statuscode StatusCode) error {
	switch statuscode {
	case twohundred:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case fourhundred:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case fivehundred:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	head := headers.NewHeaders()
	head["content-length"] = fmt.Sprintf("%d", contentLen)
	head["connection"] = "close"
	head["content-type"] = "text/plain"
	return head
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
