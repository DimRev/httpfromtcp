package response

import (
	"fmt"
	"io"

	"github.com/DimRev/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

const CRLF = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return &ErrorWritingStatusLine{Err: err}
		}
		return nil
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return &ErrorWritingStatusLine{Err: err}
		}
		return nil
	case StatusServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Server Error\r\n"))
		if err != nil {
			return &ErrorWritingStatusLine{Err: err}
		}
		return nil
	default:
		return &ErrorInvalidStatusCode{StatusCode: int(statusCode)}
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Replace("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Replace("Content-Type", "text/plain")
	h.Replace("Connection", "close")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return &ErrorWritingHeaders{Err: err}
		}
	}

	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return &ErrorWritingHeaders{Err: err}
	}

	return nil
}
