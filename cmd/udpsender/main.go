package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const DOMAIN = "127.0.0.1"
const PORT = "42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", DOMAIN, PORT))
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading: %v", err)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Fatalf("Error writing: %v", err)
		}
	}
}
