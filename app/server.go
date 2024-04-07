package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/httpmessage"
)

const (
	bufferSize = 4096
)

func main() {
	if err := runServer(); err != nil {
		log.Fatalln(err)
	}
}

func runServer() error {
	l, err := net.Listen("tcp", "127.0.0.1:4221")
	if err != nil {
		return fmt.Errorf("%s", "failed to bind to port 4221")
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %s", err.Error())
	}
	defer conn.Close()

	buffer := make([]byte, bufferSize)
	readRequest(conn, buffer)

	resp := handleRequest(buffer)

	writeResponse(conn, resp)

	return nil
}

func readRequest(conn net.Conn, buffer []byte) {
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
}

func handleRequest(input []byte) *httpmessage.Response {
	reqStr := string(bytes.TrimRight(input, "\x00"))
	request, _ := httpmessage.ParseRequest(reqStr)
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

func writeResponse(conn net.Conn, resp *httpmessage.Response) {
	_, err := conn.Write(resp.ToByteArray())
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
}
