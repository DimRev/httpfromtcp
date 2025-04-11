package main

import (
	"fmt"
	"log"
	"net"

	"github.com/DimRev/httpfromtcp/internal/request"
)

const PORT = "42069"

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %v", err)
		}
		go handleConnection(conn) // Handle each connection in a goroutine
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		return
	}

	formattedHeaders := ""
	for key, value := range req.Headers {
		formattedHeaders += fmt.Sprintf("- %s: %s\n", key, value)
	}

	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\nHeaders:\n%sBody:\n-%s\n",
		req.RequestLine.Method,
		req.RequestLine.RequestTarget,
		req.RequestLine.HttpVersion,
		formattedHeaders,
		string(req.Body),
	)
}
