package httpmessage

import (
	"regexp"
	"strings"
)

var (
	requestLineRegex = regexp.MustCompile(`^([A-Z]+) ([^ ]+) (HTTP/[0-9.]+)`)
	headerRegex      = regexp.MustCompile(`^([^:]+):\s*(.*)$`)
)

type RequestLine struct {
	Method  string
	Target  string
	Version string
}

func NewRequestLine(method, target, version string) *RequestLine {
	return &RequestLine{Method: method, Target: target, Version: version}
}

func (l RequestLine) SplitActionAndParams() (action, params string) {
	action, params, _ = strings.Cut(l.Target[1:], "/")
	return
}

type Request struct {
	RequestLine RequestLine
	Headers     Headers
	Body        string
}

func NewRequest(requestLine RequestLine, headers Headers, body string) *Request {
	return &Request{RequestLine: requestLine, Headers: headers, Body: body}
}

func ParseRequest(request string) (*Request, error) {
	parts := strings.Split(request, "\r\n")
	requestLine := parseRequestLine(parts[0])
	headers := make(Headers)
	i := 1
	for parts[i] != "" {
		k, v := parseSingleHeader(parts[i])
		headers[k] = v
		i++
	}
	body := parts[i+1]
	return NewRequest(*requestLine, headers, body), nil
}

func parseSingleHeader(s string) (k, v string) {
	parts := headerRegex.FindStringSubmatch(s)
	k, v = parts[1], parts[2]
	return
}

func parseRequestLine(s string) *RequestLine {
	parts := requestLineRegex.FindStringSubmatch(s)
	return NewRequestLine(parts[1], parts[2], parts[3])
}
