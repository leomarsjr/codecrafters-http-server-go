package main

import (
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/httpmessage"
)

const bufferSize = 4096

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) run() {
	defer c.conn.Close()

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
	request, _ := httpmessage.ParseRequest(string(req))
	action, params := request.RequestLine.SplitActionAndParams()
	switch action {
	case "":
		return httpmessage.StatusOnlyResponse(httpmessage.StatusOK)
	case "echo":
		return httpmessage.EchoResponse(params)
	case "user-agent":
		return httpmessage.UserAgentResponse(request.Headers["User-Agent"])
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
