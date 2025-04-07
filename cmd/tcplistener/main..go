package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const PATH_TO_FILE = "./messages.txt"
const BUFFER_SIZE = 8
const PORT = "42069"

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting: %v", err)
	}
	defer conn.Close()

	linesCh := getLinesChannel(conn)

	for line := range linesCh {
		fmt.Printf("%s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan string)
	go func() {
		b := make([]byte, BUFFER_SIZE)
		currentLine := ""
		for {
			n, err := f.Read(b)
			if err != nil && err != io.EOF {
				close(lineCh)
				log.Fatalf("Error reading file: %v", err)
			}
			if n > 0 {
				currentLine += string(b[:n])
			}

			splitLines := strings.Split(currentLine, "\n")
			for i := 0; i < len(splitLines)-1; i++ {
				lineCh <- splitLines[i]
			}
			currentLine = splitLines[len(splitLines)-1]

			if err == io.EOF {
				if currentLine != "" {
					lineCh <- currentLine
				}
				close(lineCh)
				break
			}
		}
	}()

	return lineCh
}
