package httpmessage

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

// Response is the representation of an HTTP response.
// It contains the status, the headers and the body of a response.
type Response struct {
	status  Status  // HTTP status (numeric code and text)
	headers Headers // HTTP headers (key-value pairs)
	body    string  // Body in the text format
}

// NewResponse creates a response with given status, headers and body.
func NewResponse(status Status, headers Headers, body string) *Response {
	return &Response{status, headers, body}
}

// StatusOnlyResponse returns a response with status, but empty headers and body.
func StatusOnlyResponse(status Status) *Response {
	return NewResponse(status, EmptyHeaders, emptyBody)
}

// EchoResponse returns a response with the given body,
// status "200 OK" and content type and length headers.
func EchoResponse(body string) *Response {
	headers := Headers{}
	headers["Content-Type"] = textPlainContentType
	headers["Content-Length"] = strconv.Itoa(len(body))
	return NewResponse(StatusOK, headers, body)
}

// UserAgentResponse returns a response with the given agent as body content,
// status "200 OK" and content type and length headers.
func UserAgentResponse(agent string) *Response {
	// Despite this implementation being the same as the EchoResponse,
	// we keep them separated with some duplication,
	// in case they evolve differently in the future.
	headers := Headers{}
	headers["Content-Type"] = textPlainContentType
	headers["Content-Length"] = strconv.Itoa(len(agent))
	return NewResponse(StatusOK, headers, agent)
}

// GetFileResponse returns a response with the content of a file
// identified by fileName and contained in directory.
// In case of success, the response body will be the content of the file,
// with status "200 OK" and content type and length headers.
// In case the file doesn't exist in directory, status "404 Not Found" is returned.
func GetFileResponse(directory, fileName string) *Response {
	file, err := readFile(path.Join(directory, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return StatusOnlyResponse(StatusNotFound)
	}
	headers := Headers{}
	headers["Content-Type"] = octetStreamContentType
	headers["Content-Length"] = strconv.Itoa(len(file))
	return NewResponse(StatusOK, headers, string(file))
}

func readFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(f)
	return io.ReadAll(f)
}

// PostFileResponse creates a file named fileName in directory, with body as content.
// If the file was created successfully, returns a response with status "201 Created".
// If an error occurred, returns a response with status "500 Internal Server Error".
func PostFileResponse(directory, fileName, body string) *Response {
	err := createFile(path.Join(directory, fileName), body)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return StatusOnlyResponse(StatusInternalServerError)
	}
	return StatusOnlyResponse(StatusCreated)
}

func createFile(file, body string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(f)
	if _, err := f.WriteString(body); err != nil {
		return err
	}
	return nil
}

// String returns the response as text, following the HTTP standard.
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

// ToByteArray returns the response as an array of bytes.
func (r Response) ToByteArray() []byte {
	return []byte(r.String())
}
