// Package httpmessage contains the basic structures and functionality
// for HTTP [Request] and [Response] handling.
package httpmessage

import (
	"fmt"
	"strings"
)

const (
	protocolVersion = "HTTP/1.1"

	emptyBody = ""

	textPlainContentType   = "text/plain"
	octetStreamContentType = "application/octet-stream"
)

// List of HTTP Status used by this implementation.
var (
	StatusOK                  = Status{200, "OK"}
	StatusCreated             = Status{201, "Created"}
	StatusNotFound            = Status{404, "Not Found"}
	StatusInternalServerError = Status{500, "Internal Server Error"}
	StatusMethodNotAllowed    = Status{503, "Method Not Allowed"}
)

// EmptyHeaders is used as a representation of empty Headers.
var EmptyHeaders = Headers{}

// Status represents HTTP status codes, indicating success or failure of the request.
// Contains the numeric code and the text description.
type Status struct {
	code int
	text string
}

// String returns the status and code text in a standard format.
func (s Status) String() string {
	return fmt.Sprintf("%d %s", s.code, s.text)
}

// Headers represent the HTTP request or response headers.
// It is a simple map of key-value pairs (strings).
type Headers map[string]string

// String returns the headers according to the HTTP standard.
// Format:
//   - each header in a single line, separated by '\r\n',
//     with key and value like 'Key: Value'
func (h Headers) String() string {
	output := make([]string, 0, len(h))
	for k, v := range h {
		output = append(output, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(output, "\r\n")
}
