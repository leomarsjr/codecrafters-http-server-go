package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	bufferSize         = 4096
	httpStatusOk       = 200
	httpStatusNotFound = 404
)

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

	status := parseRequest(buffer)

	writeResponse(conn, status)
}

func parseRequest(input []byte) int {
	path := strings.Fields(string(input))[1]
	if path != "/" {
		return httpStatusNotFound
	}
	return httpStatusOk
}

func readRequest(conn net.Conn, buffer []byte) {
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
}

func writeResponse(conn net.Conn, status int) {
	responseHeader := createResponseHeader(status)
	_, err := conn.Write([]byte(responseHeader))
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
}

func createResponseHeader(status int) string {
	switch status {
	case httpStatusOk:
		return "HTTP/1.1 200 OK\r\n\r\n"
	case httpStatusNotFound:
		return "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	}
	return ""
}
