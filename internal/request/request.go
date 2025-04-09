package request

import (
	"io"
	"slices"
	"strings"
	// Import errors package if you use custom error types defined elsewhere
	// "errors"
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

const BUFFER_SIZE = 4096
const CRLF = "\r\n"

type ParsingStage int

const (
	RequestLineParsingStage ParsingStage = iota
	HeadersParsingStage
	BodyParsingStage
)

func NewRequest() *Request {
	return &Request{
		RequestLine: NewRequestLine(),
		Headers:     make(map[string]string),
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
	buf := make([]byte, BUFFER_SIZE)
	buffer := ""
	req := NewRequest()
	stage := RequestLineParsingStage

	for {
		remainingBuffer, parseErr := req.parse(buffer, &stage)
		if parseErr != nil {
			return nil, parseErr
		}
		buffer = remainingBuffer

		if stage == BodyParsingStage {
			req.Body = []byte(buffer)
			return req, nil
		}

		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, &ErrorReadingRequest{Err: err}
		}

		if n == 0 {
			continue
		}

		buffer += string(buf[:n])
	}
}

func (r *Request) parse(currentBuffer string, stage *ParsingStage) (string, error) {
	for *stage != BodyParsingStage {
		newLineIndex := strings.Index(currentBuffer, CRLF)
		if newLineIndex == -1 {
			return currentBuffer, nil
		}

		line := currentBuffer[:newLineIndex]
		remainingBuffer := currentBuffer[newLineIndex+len(CRLF):]

		switch *stage {
		case RequestLineParsingStage:
			if line == "" {
				currentBuffer = remainingBuffer
				continue
			}
			if err := r.parseRequestLine(line); err != nil {
				return "", err
			}
			*stage = HeadersParsingStage
			currentBuffer = remainingBuffer

		case HeadersParsingStage:
			if line == "" {
				*stage = BodyParsingStage
				return remainingBuffer, nil
			}
			if err := r.parseHeaders(line); err != nil {
				return "", err
			}
			currentBuffer = remainingBuffer
		}
	}
	return currentBuffer, nil
}

func (r *Request) parseRequestLine(line string) error {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ErrorInvalidRequestLine{
			Line: line,
		}
	}

	if err := r.assignMethod(parts[0]); err != nil {
		return err
	}

	if err := r.assignRequestTarget(parts[1]); err != nil {
		return err
	}

	if err := r.assignHttpVersion(parts[2]); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseHeaders(line string) error {
	line = strings.TrimSpace(line)
	headerParts := strings.SplitN(line, ": ", 2)
	if len(headerParts) != 2 {
		return ErrorInvalidHeaderLine{
			HeaderLine: line,
		}
	}
	key := strings.TrimSpace(headerParts[0])
	value := strings.TrimSpace(headerParts[1])
	r.Headers[key] = value
	return nil
}

func (r *Request) assignMethod(m string) error {
	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	if !slices.Contains(validMethods, m) {
		return ErrorInvalidMethod{
			Method: m,
		}
	}
	r.RequestLine.Method = m
	return nil
}

func (r *Request) assignRequestTarget(target string) error {
	if !strings.HasPrefix(target, "/") {
		return ErrorInvalidRequestTarget{
			Target: target,
		}
	}
	r.RequestLine.RequestTarget = target
	return nil
}

func (r *Request) assignHttpVersion(version string) error {
	if version != "HTTP/1.1" {
		return ErrorInvalidHTTPVersion{
			Version: version,
		}
	}
	parts := strings.Split(version, "/")
	if len(parts) != 2 {
		return ErrorInvalidHTTPVersion{
			Version: version,
		}
	}
	r.RequestLine.HttpVersion = parts[1]
	return nil
}
