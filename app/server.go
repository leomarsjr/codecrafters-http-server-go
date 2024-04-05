package main

import (
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/codecrafters-io/http-server-starter-go/httpmessage"
)

const (
	bufferSize = 4096
)

var pathRegex = regexp.MustCompile(`/(\w*)/?(\S*)`)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, bufferSize)
	readRequest(conn, buffer)

	resp := handleRequest(buffer)

	writeResponse(conn, resp)
}

func readRequest(conn net.Conn, buffer []byte) {
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
}

func handleRequest(input []byte) *httpmessage.Response {
	matches := pathRegex.FindStringSubmatch(string(input))
	switch matches[1] {
	case "":
		return httpmessage.StatusOnlyResponse(httpmessage.StatusOK)
	case "echo":
		return httpmessage.EchoResponse(matches[2])
	default:
		return httpmessage.StatusOnlyResponse(httpmessage.StatusNotFound)
	}
}

func writeResponse(conn net.Conn, resp *httpmessage.Response) {
	_, err := conn.Write(resp.ToByteArray())
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
}
