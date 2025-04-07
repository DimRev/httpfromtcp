package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const BUFFER_SIZE = 8

type ParsingStage int

const (
	RequestLineParsingStage ParsingStage = iota
	HeadersParsingStage
	BodyParsingStage
)

func NewRequest() *Request {
	newHeaders := make(map[string]string)
	newRequestLine := NewRequestLine()

	return &Request{
		RequestLine: newRequestLine,
		Headers:     newHeaders,
		Body:        []byte{},
	}
}

func NewRequestLine() RequestLine {
	return RequestLine{
		HttpVersion:   "",
		RequestTarget: "",
		Method:        "",
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, BUFFER_SIZE)
	buffer := ""
	req := NewRequest()

	var phase ParsingStage = RequestLineParsingStage
outer:
	for {
		n, err := reader.Read(b)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading request: %v", err)
		}

		buffer += string(b[:n])
		lines := strings.Split(buffer, "\r\n")
		buffer = lines[len(lines)-1]

		for _, line := range lines[:len(lines)-1] {
			switch phase {
			case RequestLineParsingStage:
				if err := req.parseRequestLine(line); err != nil {
					return nil, err
				}
				phase = HeadersParsingStage
			case HeadersParsingStage:
				if line == "" {
					break outer
				}
				parts := strings.SplitN(line, ": ", 2)
				if len(parts) != int(BodyParsingStage) {
					return nil, fmt.Errorf("invalid header line: %s", line)
				}
				req.Headers[parts[0]] = parts[1]
			}
		}

		if err == io.EOF {
			break
		}
	}

	return req, nil
}

func (r *Request) parseRequestLine(line string) error {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return fmt.Errorf("invalid request line: %s", line)
	}
	Methods := []string{"GET", "POST", "PUT", "DELETE"}
	if !slices.Contains(Methods, parts[0]) {
		return fmt.Errorf("invalid method: %s", parts[0])
	} else {
		r.RequestLine.Method = parts[0]
	}

	if strings.HasPrefix(parts[1], "/") {
		r.RequestLine.RequestTarget = parts[1]
	} else {
		return fmt.Errorf("invalid request target: %s", parts[1])
	}

	if parts[2] != "HTTP/1.1" {
		return fmt.Errorf("invalid http version: %s", parts[2])
	} else {
		r.RequestLine.HttpVersion = strings.Split(parts[2], "/")[1]
	}

	return nil
}
