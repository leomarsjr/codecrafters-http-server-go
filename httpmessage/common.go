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
