package request

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strings"
	// Import errors package if you use custom error types defined elsewhere
	// "errors"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const BUFFER_SIZE = 8
const CRLF = "\r\n"

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

func NewRequestLine() RequestLine {
	return RequestLine{
		HttpVersion:   "",
		RequestTarget: "",
		Method:        "",
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE)
	readToIndex := 0
	req := &Request{
		state: requestStateInitialized,
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, &ErrorIncompleteRequest{}
				}
				break
			}
			return nil, &ErrorUnexpectedReadError{Err: err}
		}

		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func (r *Request) parse(currentBuffer []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(currentBuffer)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateDone
		return n, nil
	case requestStateDone:
		return 0, &ErrorParsingTryingToReadAfterDone{}
	default:
		return 0, &ErrorParsingUnknownState{State: r.state}
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return nil, 0, nil
	}
	line := string(data[:idx])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, &ErrorParsingRequestLineMalformed{Line: line}
	}

	method, err := assignMethod(parts[0])
	if err != nil {
		return nil, 0, err
	}
	target, err := assignRequestTarget(parts[1])
	if err != nil {
		return nil, 0, err
	}
	version, err := assignHttpVersion(parts[2])
	if err != nil {
		return nil, 0, err
	}

	requestLine := &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}

	return requestLine, idx + len(CRLF), nil
}

func assignMethod(m string) (string, error) {
	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	if !slices.Contains(validMethods, m) {
		return "", &ErrorParsingRequestInvalidMethod{
			Method: m,
		}
	}
	return m, nil
}

func assignRequestTarget(target string) (string, error) {
	if !strings.HasPrefix(target, "/") {
		return "", &ErrorParsingRequestInvalidTarget{
			Target: target,
		}
	}
	return target, nil
}

func assignHttpVersion(version string) (string, error) {
	if version != "HTTP/1.1" {
		return "", &ErrorParsingRequestInvalidVersion{
			Version: version,
		}
	}
	parts := strings.Split(version, "/")
	if len(parts) != 2 {
		return "", &ErrorParsingRequestInvalidVersion{
			Version: version,
		}
	}
	return parts[1], nil
}
