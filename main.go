package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const PATH_TO_FILE = "./messages.txt"
const BUFFER_SIZE = 8

func main() {
	file, err := os.Open(PATH_TO_FILE)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	b := make([]byte, BUFFER_SIZE)
	currentLine := ""

	for {
		n, err := file.Read(b)
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading file: %v", err)
		}
		if n > 0 {
			currentLine += string(b[:n])
		}

		splitLines := strings.Split(currentLine, "\n")
		for i := 0; i < len(splitLines)-1; i++ {
			fmt.Printf("read: %s\n", splitLines[i])
		}
		currentLine = splitLines[len(splitLines)-1]

		if err == io.EOF {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
			}
			break
		}
	}

}
