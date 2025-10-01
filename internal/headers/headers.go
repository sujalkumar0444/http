package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= 'a' && ch <= 'z' ||
			ch >= 'A' && ch <= 'Z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

var CRLF = []byte("\r\n")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) (string, bool) {
	str, ok := h.headers[strings.ToLower(name)]
	return str, ok
}

func (h *Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	if _, ok := h.headers[name]; ok {
		h.headers[name] = value
	}
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], CRLF)
		if idx == -1 {
			break
		}

		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(CRLF)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed field name")
		}

		read += idx + len(CRLF)
		h.Set(name, value)
	}

	return read, done, nil
}
