package httpmessage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type Response struct {
	status  Status
	headers Headers
	body    string
}

func NewResponse(status Status, headers Headers, body string) *Response {
	return &Response{status, headers, body}
}

func StatusOnlyResponse(status Status) *Response {
	return NewResponse(status, EmptyHeaders, emptyBody)
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

func GetFileResponse(directory, fileName string) *Response {
	headers := Headers{}
	headers["Content-Type"] = octetStreamContentType
	file, err := readFile(path.Join(directory, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return NewResponse(StatusNotFound, headers, emptyBody)
	}
	headers["Content-Length"] = strconv.Itoa(len(file))
	return NewResponse(StatusOK, headers, string(file))
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func PostFileResponse(directory, fileName, body string) *Response {
	err := createFile(path.Join(directory, fileName), body)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return StatusOnlyResponse(StatusInternalServerError)
	}
	return StatusOnlyResponse(StatusCreated)
}

func createFile(file, body string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(body); err != nil {
		return err
	}
	return nil
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
