package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var directory string

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--directory" {
		directory = os.Args[2]
	}
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

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %s", err.Error())
		}
		go NewClient(conn).run()
	}
}
