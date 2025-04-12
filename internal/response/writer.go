package response

import (
	"fmt"
	"io"

	"github.com/DimRev/httpfromtcp/internal/headers"
)

type Writer struct {
	StatusLine []byte
	Headers    headers.Headers
	Body       []byte

	writerState writerState
}

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateDone
)

func NewWriter() *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateStatusLine}
	}
	w.StatusLine = getStatusLine(statusCode)
	w.writerState = writerStateHeaders
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateHeaders}
	}
	w.Headers = headers
	w.writerState = writerStateBody
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateBody}
	}
	w.Body = p
	w.writerState = writerStateDone
	return len(p), nil
}

func (w *Writer) Write(conn io.Writer) error {
	if w.writerState != writerStateDone {
		return &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateDone}
	}
	_, err := conn.Write(w.StatusLine)
	if err != nil {
		return &ErrorWritingStatusLine{Err: err}
	}
	headers := w.Headers
	for key, value := range headers {
		_, err := conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return &ErrorWritingHeaders{Err: err}
		}
	}
	_, err = conn.Write([]byte(CRLF))
	if err != nil {
		return &ErrorWritingHeaders{Err: err}
	}
	_, err = conn.Write(w.Body)
	if err != nil {
		return &ErrorWritingBody{Err: err}
	}
	return nil
}
