package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"
const allowed = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	header_str := string(data[:idx])
	split := strings.SplitN(header_str, ":", 2)
	if len(split) != 2 {
		return 0, false, errors.New("invalid header")
	}
	if strings.ContainsAny(split[0], " \t") {
		return 0, false, errors.New("Header cannot start withleading white space")
	}
	value := strings.ToLower(strings.TrimSpace(split[1]))
	key := strings.ToLower(strings.TrimSpace(split[0]))
	for _, char := range key {
		if !strings.Contains(allowed, string(char)) {
			return 0, false, errors.New("invalid character in key value of header")
		}
	}
	_, ok := h[key]
	if ok {
		h[key] += ", " + value
	} else {
		h[key] = value
	}
	return idx + len(crlf), false, nil
}
