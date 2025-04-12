package response

import (
	"fmt"
	"io"
	"net"

	"github.com/DimRev/httpfromtcp/internal/headers"
)

type Writer struct {
	StatusLine []byte
	Headers    headers.Headers
	Body       []byte

	writerState writerState
	writer      io.Writer
}

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateDone
)

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer:      conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateStatusLine}
	}
	w.StatusLine = getStatusLine(statusCode)
	_, err := w.writer.Write(w.StatusLine)
	if err != nil {
		return &ErrorWritingStatusLine{Err: err}
	}
	w.writerState = writerStateHeaders
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateHeaders}
	}
	w.Headers = headers
	for key, value := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return &ErrorWritingHeaders{Err: err}
		}
	}
	_, err := w.writer.Write([]byte(CRLF))
	if err != nil {
		return &ErrorWritingHeaders{Err: err}
	}
	w.writerState = writerStateBody
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateBody}
	}
	w.Body = p
	_, err := w.writer.Write(p)
	if err != nil {
		return 0, &ErrorWritingBody{Err: err}
	}
	w.writerState = writerStateDone
	return len(p), nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, &ErrorInvalidWriterState{CurrentState: w.writerState, ExpectedState: writerStateBody}
	}

	chunkSize := len(p)
	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, &ErrorWritingChunkedBody{Err: err}
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, &ErrorWritingChunkedBody{Err: err}
	}
	nTotal += n

	n, err = fmt.Fprintf(w.writer, "\r\n")
	if err != nil {
		return nTotal, &ErrorWritingChunkedBody{Err: err}
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
}
