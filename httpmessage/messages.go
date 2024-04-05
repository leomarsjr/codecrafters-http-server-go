package httpmessage

import (
	"fmt"
	"strings"
)

const (
	protocolVersion = "HTTP/1.1"
	EmptyBody       = ""
)

var (
	StatusOK       = Status{200, "OK"}
	StatusNotFound = Status{404, "Not Found"}

	EmptyHeaders = make(Headers)
)

type Status struct {
	code int
	text string
}

func (s Status) String() string {
	return fmt.Sprintf("%d %s", s.code, s.text)
}

type Headers map[string]string

func (h Headers) String() string {
	output := make([]string, 0, len(h))
	for k, v := range h {
		output = append(output, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(output, "\r\n")
}

type Response struct {
	status  Status
	headers Headers
	body    string
}

func NewResponse(status Status, headers Headers, body string) *Response {
	return &Response{status, headers, body}
}

func StatusOnlyResponse(status Status) *Response {
	return NewResponse(status, EmptyHeaders, EmptyBody)
}

func (r Response) String() string {
	var resp strings.Builder
	fmt.Fprintf(&resp, "%s %s\r\n", protocolVersion, r.status)
	if len(r.headers) > 0 {
		fmt.Fprintf(&resp, "%s\r\n", r.headers)
	}
	if r.body != "" {
		fmt.Fprintf(&resp, "\r\n%s\r\n", r.body)
	}
	fmt.Fprintf(&resp, "\r\n")
	return resp.String()
}

func (r Response) ToByteArray() []byte {
	return []byte(r.String())
}
