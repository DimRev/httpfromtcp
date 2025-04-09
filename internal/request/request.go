package request

import (
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

const BUFFER_SIZE = 1
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
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, &ErrorReadingRequest{
				Err: err,
			}
		}

		buffer += string(buf[:n])

		var parseErr error
		buffer, parseErr = req.parse(buffer, &stage)
		if parseErr != nil {
			return nil, parseErr
		}

		if err == io.EOF {
			if stage == BodyParsingStage {
				req.Body = []byte(buffer)
			}
			break
		}
	}

	return req, nil
}

func (r *Request) parse(buffer string, stage *ParsingStage) (string, error) {
	for {
		newLineIndex := strings.Index(buffer, CRLF)
		if newLineIndex == -1 {
			break
		}

		line := buffer[:newLineIndex]
		buffer = buffer[newLineIndex+len(CRLF):]

		switch *stage {
		case RequestLineParsingStage:
			if line == "" {
				continue
			}
			if err := r.parseRequestLine(line); err != nil {
				return "", err
			}
			*stage = HeadersParsingStage

		case HeadersParsingStage:
			if line == "" {
				*stage = BodyParsingStage
				continue
			}
			if err := r.parseHeaders(line); err != nil {
				return "", err
			}
		}
	}
	return buffer, nil
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
