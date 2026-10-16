package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type state int

const (
	initalized state = iota
	done
)

type Request struct {
	state_val   state
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const buf_size = 8

func (r *Request) parse(data []byte) (int, error) {
	if r.state_val == initalized {
		rl, num_bytes, err := ParseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if num_bytes == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.state_val = done
		return num_bytes, nil
	}

	if r.state_val == done {
		return 0, errors.New("error: trying to read data in a done state")
	}

	return 0, errors.New("error: unknown state baka")
}

func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, buf_size)
	readToIndex := 0
	req := Request{state_val: initalized}
	for req.state_val != done {
		if readToIndex == len(buf) {
			newbuf := make([]byte, len(buf)*2)
			copy(newbuf, buf)
			buf = newbuf
		}
		read, err := r.Read(buf[readToIndex:])
		if err == io.EOF {
			req.state_val = done
			break
		}
		readToIndex += read
		read, err = req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[read:readToIndex])
		readToIndex -= read
	}
	return &req, nil
}

func ParseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	reqLineTxt := data[:idx]
	reqLine, err := RequestLineFromStr(string(reqLineTxt))
	if err != nil {
		return nil, 0, err
	}
	return reqLine, idx + len(crlf), nil
}

func RequestLineFromStr(str string) (*RequestLine, error) {
	split_parts := strings.Split(str, " ")
	if len(split_parts) != 3 {
		return nil, errors.New("not properly formatted request line")
	}

	for _, char := range split_parts[0] {
		if char < 'A' || char > 'Z' {
			return nil, errors.New("method must be upper alphabetic char")
		}
	}

	ver := strings.Split(split_parts[2], "/")
	if len(ver) != 2 {
		return nil, errors.New("http version not properly formatted")
	}

	if ver[0] != "HTTP" {
		return nil, errors.New("unproper http version")
	}

	if ver[1] != "1.1" {
		return nil, errors.New("unproper http version")
	}

	return &RequestLine{
		Method:        split_parts[0],
		RequestTarget: split_parts[1],
		HttpVersion:   ver[1],
	}, nil
}
