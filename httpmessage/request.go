package httpmessage

import (
	"regexp"
	"strings"
)

var (
	requestLineRegex = regexp.MustCompile(`^([A-Z]+) ([^ ]+) (HTTP/[0-9.]+)`)
	headerRegex      = regexp.MustCompile(`^([^:]+):\s*(.*)$`)
)

// RequestLine contains the method, the URL target and the protocol version of a request.
type RequestLine struct {
	Method  string // HTTP method (e.g. GET, POST, PUT)
	Target  string // URL of the request target
	Version string // HTTP protocol version (HTTP/1.1)
}

// NewRequestLine creates a request line with method, target and version.
func NewRequestLine(method, target, version string) *RequestLine {
	return &RequestLine{Method: method, Target: target, Version: version}
}

// SplitActionAndParams interprets the request line target.
// action is the first part of the target, representing what the server needs to process.
// params are the arguments of the action the server will process.
func (l RequestLine) SplitActionAndParams() (action, params string) {
	action, params, _ = strings.Cut(l.Target[1:], "/")
	return
}

// Request is the representation of an HTTP request.
// It contains the request line, the headers and the body of a request.
type Request struct {
	RequestLine RequestLine // Request line (method, target, version)
	Headers     Headers     // HTTP Headers (key-value pairs)
	Body        string      // Body in the text format
}

// NewRequest creates a request with request line, headers and body.
func NewRequest(requestLine RequestLine, headers Headers, body string) *Request {
	return &Request{RequestLine: requestLine, Headers: headers, Body: body}
}

// ParseRequest returns a struct representing the request.
func ParseRequest(request string) *Request {
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
	return NewRequest(*requestLine, headers, body)
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
