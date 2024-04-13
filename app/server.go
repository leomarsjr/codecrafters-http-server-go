// A simple implementation of own HTTP server, for learning purposes.
// For more information, visit CodeCrafters [catalog],
// section "Build your own HTTP server".
//
// The command accepts an optional flag:
//
//	--directory directory-name
//		name of the local directory for reading and writing files.
//
// [catalog]: https://app.codecrafters.io/catalog
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var directory string // name of the directory to read or write files

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
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %s", err.Error())
		}
		go NewClient(conn).run()
	}
}
