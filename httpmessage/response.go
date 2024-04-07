package httpmessage

import (
	"fmt"
	"strconv"
	"strings"
)

const textPlainContentType = "text/plain"

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

func EchoResponse(body string) *Response {
	headers := Headers{}
	headers["Content-Type"] = textPlainContentType
	headers["Content-Length"] = strconv.Itoa(len(body))
	return NewResponse(StatusOK, headers, body)
}

func UserAgentResponse(agent string) *Response {
	headers := Headers{}
	headers["Content-Type"] = textPlainContentType
	headers["Content-Length"] = strconv.Itoa(len(agent))
	return NewResponse(StatusOK, headers, agent)
}

func (r Response) String() string {
	var resp strings.Builder
	fmt.Fprintf(&resp, "%s %s\r\n", protocolVersion, r.status)
	if len(r.headers) > 0 {
		fmt.Fprintf(&resp, "%s\r\n", r.headers)
	}
	fmt.Fprintf(&resp, "\r\n")
	if r.body != "" {
		fmt.Fprintf(&resp, "%s", r.body)
	}
	return resp.String()
}

func (r Response) ToByteArray() []byte {
	return []byte(r.String())
}
