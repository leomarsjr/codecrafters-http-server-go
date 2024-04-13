package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/httpmessage"
)

const bufferSize = 4096

// Client of the HTTP server that handles the connection.
type Client struct {
	conn net.Conn
}

// NewClient returns a client to handle the connection conn.
func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) run() {
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing client connection: %v", err)
		}
	}(c.conn)

	req, err := c.readRequest()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := c.handleRequest(req)

	if err := c.writeResponse(resp); err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Client) readRequest() ([]byte, error) {
	buffer := make([]byte, bufferSize)
	var data []byte
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return nil, err
		}
		data = append(data, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	return data, nil
}

func (c *Client) handleRequest(req []byte) *httpmessage.Response {
	request := httpmessage.ParseRequest(string(req))
	action, params := request.RequestLine.SplitActionAndParams()
	switch action {
	case "":
		return httpmessage.StatusOnlyResponse(httpmessage.StatusOK)
	case "echo":
		return httpmessage.EchoResponse(params)
	case "user-agent":
		return httpmessage.UserAgentResponse(request.Headers["User-Agent"])
	case "files":
		switch request.RequestLine.Method {
		case "GET":
			return httpmessage.GetFileResponse(directory, params)
		case "POST":
			return httpmessage.PostFileResponse(directory, params, request.Body)
		default:
			return httpmessage.StatusOnlyResponse(httpmessage.StatusMethodNotAllowed)
		}
	default:
		return httpmessage.StatusOnlyResponse(httpmessage.StatusNotFound)
	}
}

func (c *Client) writeResponse(resp *httpmessage.Response) error {
	_, err := c.conn.Write(resp.ToByteArray())
	if err != nil {
		return fmt.Errorf("failed to write response: %s", err.Error())
	}
	return nil
}
